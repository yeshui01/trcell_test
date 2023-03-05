/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-06-15 14:14:17
 * @LastEditTime: 2022-10-14 13:47:36
 * @Brief:
 */
package trframe

import (
	"errors"
	"fmt"
	"time"

	"trcell/pkg/evhub"
	"trcell/pkg/loghlp"
	"trcell/pkg/protocol"
	"trcell/pkg/timeutil"
	"trcell/pkg/trframe/iframe"
	"trcell/pkg/trframe/tframeconfig"
	"trcell/pkg/trframe/tframedispatcher"
	"trcell/pkg/trframe/trnode"
	"trcell/pkg/trframe/trnode/tcellcenter"
	"trcell/pkg/trframe/trnode/tcelldata"
	"trcell/pkg/trframe/trnode/tcellgame"
	"trcell/pkg/trframe/trnode/tcellgate"
	"trcell/pkg/trframe/trnode/tcelllog"
	"trcell/pkg/trframe/trnode/tcelllogic"
	"trcell/pkg/trframe/trnode/tcellroot"
	"trcell/pkg/trframe/trnode/tcellsocial"
	"trcell/pkg/trframe/trnode/tcelltrans"
	"trcell/pkg/trframe/trnode/tcellview"

	"github.com/sirupsen/logrus"
)

type EFrameStep int32

const (
	ETRFrameStepCheck  EFrameStep = 0 // 运行检测
	ETRFrameStepInit   EFrameStep = 1 // 初始化
	ETRFrameStepPreRun EFrameStep = 2 // 准备运行
	ETRFrameStepRun    EFrameStep = 3 // 正常运行
	ETRFrameStepStop   EFrameStep = 4 // 停止
	ETRFrameStepEnd    EFrameStep = 5 // 结束
	ETRFrameStepExit   EFrameStep = 6 // 退出
	ETRFrameStepFinal  EFrameStep = 7 // 边界值
)

type FrameRunFunc func(curTimeMs int64)

// 包装一下hubcommand
type TRFrameCommand struct {
	UserCmd *evhub.HubCommand
	Hub     *evhub.EventHub
}
type TRFrame struct {
	frameNodeMgr    *FrameNodeMgr
	evHub           *evhub.EventHub
	runStep         EFrameStep
	loopFuncList    []FrameRunFunc
	curWorkNode     ITRFrameWorkNode // 当前工作节点
	frameConfig     *tframeconfig.FrameConfig
	nodeType        int32
	nodeIndex       int32
	netSessionList  map[int32]*FrameSession
	userStepRun     []func(curTimeMs int64) bool
	userCmdHandle   func(userCmd *TRFrameCommand)
	remoteMsgMgr    *RemoteMsgCallMgr
	multiMsgCallMgr *MultiMsgCallMgr
	msgDispatcher   *tframedispatcher.FrameMsgDispatcher
	keepNodeTime    int64
	nowFrameTimeMs  int64
	msgDoneList     []func() // 消息处理后的执行列表
}

func newTRFrame(configPath string, nType int32, nIndex int32) *TRFrame {
	tf := &TRFrame{
		frameNodeMgr:   NewFrameNodeMgr(),
		evHub:          evhub.NewHub(),
		runStep:        ETRFrameStepCheck,
		frameConfig:    tframeconfig.NewFrameConfig(),
		nodeType:       nType,
		nodeIndex:      nIndex,
		netSessionList: make(map[int32]*FrameSession),
		userStepRun:    make([]func(curTimeMs int64) bool, int(ETRFrameStepFinal)),
		nowFrameTimeMs: timeutil.NowTimeMs(),
	}
	// 加载配置
	err := tframeconfig.ReadFrameConfigFromFile(configPath, tf.frameConfig)
	if err != nil {
		panic(fmt.Sprintf("load frame config error:%s", err.Error()))
	}

	// 分发器
	tf.msgDispatcher = tframedispatcher.NewFrameMsgDispatcher(tf)
	// 初始化当前节点
	tf.curWorkNode = tf.makeInitWorkNode(nType, nIndex)
	if tf.curWorkNode == nil {
		panic(fmt.Sprintf("worknode(%d,%d) is nil", nType, nIndex))
	}
	tf.initNodeSetting()
	tf.remoteMsgMgr = newRemoteMsgMgr(tf)
	tf.multiMsgCallMgr = newMultiMsgCallMgr(tf)
	tf.regCallback()
	tf.evHub.AddFrameLoopFunc(func(curTimeMs int64) {
		tf.frameRun(curTimeMs)
	})
	return tf
}

func (tf *TRFrame) frameRun(curTimeMs int64) {
	tf.nowFrameTimeMs = curTimeMs
	tf.remoteMsgMgr.update(curTimeMs)
	tf.multiMsgCallMgr.update(curTimeMs)
	switch tf.runStep {
	case ETRFrameStepCheck:
		{
			if tf.RunStepCheck(curTimeMs) {
				tf.changeNextStep()
			}
			break
		}
	case ETRFrameStepInit:
		{
			if tf.RunStepInit(curTimeMs) {
				errListen := tf.listen()
				if errListen != nil {
					loghlp.Errorf("listenError:%s", errListen.Error())
					time.Sleep(time.Second * 2)
				} else {
					loghlp.Infof("listen succ")
					tf.changeNextStep()
				}
			}
			break
		}
	case ETRFrameStepPreRun:
		{
			if tf.RunStepPreRun(curTimeMs) {
				tf.changeNextStep()
			}
			break
		}
	case ETRFrameStepRun:
		{
			tf.RunStepRun(curTimeMs)
			break
		}
	case ETRFrameStepStop:
		{
			if tf.RunStepStop(curTimeMs) {
				tf.changeNextStep()
			}
			break
		}
	case ETRFrameStepEnd:
		{
			if tf.RunStepEnd(curTimeMs) {
				tf.changeNextStep()
			}
			break
		}
	}
}

func (tf *TRFrame) Run() {
	tf.evHub.Run()
}

func (tf *TRFrame) Stop() {
	tf.evHub.Stop()
}

func (tf *TRFrame) changeNextStep() {
	tf.runStep++
	loghlp.Infof("change trframe_step(%d) before step:%d", tf.runStep, tf.runStep-1)
}

func (tf *TRFrame) listen() error {
	var listenMode int32
	var listenAddr string

	switch tf.nodeType {
	case trnode.ETRNodeTypeCellRoot:
		{
			rootCfg := tf.frameConfig.CellRootCfgs[tf.nodeIndex]
			if rootCfg.ListenMode == "unix" {
				listenMode = evhub.ListenModeUnix
			} else {
				listenMode = evhub.ListenModeTcp
			}
			listenAddr = rootCfg.ListenAddr
			break
		}
	case trnode.ETRNodeTypeCellGate:
		{
			listenCfg := tf.frameConfig.CellGateCfgs[tf.nodeIndex]
			if listenCfg.ListenMode == "unix" {
				listenMode = evhub.ListenModeUnix
			} else {
				listenMode = evhub.ListenModeTcp
			}
			listenAddr = listenCfg.ListenAddr
			break
		}
	case trnode.ETRNodeTypeCellData:
		{
			listenCfg := tf.frameConfig.CellDataCfgs[tf.nodeIndex]
			if listenCfg.ListenMode == "unix" {
				listenMode = evhub.ListenModeUnix
			} else {
				listenMode = evhub.ListenModeTcp
			}
			listenAddr = listenCfg.ListenAddr
			break
		}
	case trnode.ETRNodeTypeCellCenter:
		{
			listenCfg := tf.frameConfig.CellCenterCfgs[tf.nodeIndex]
			if listenCfg.ListenMode == "unix" {
				listenMode = evhub.ListenModeUnix
			} else {
				listenMode = evhub.ListenModeTcp
			}
			listenAddr = listenCfg.ListenAddr
			break
		}
	case trnode.ETRNodeTypeCellLogic:
		{
			listenCfg := tf.frameConfig.CellLogicCfgs[tf.nodeIndex]
			if listenCfg.ListenMode == "unix" {
				listenMode = evhub.ListenModeUnix
			} else {
				listenMode = evhub.ListenModeTcp
			}
			listenAddr = listenCfg.ListenAddr
			break
		}
	case trnode.ETRNodeTypeCellGame:
		{
			listenCfg := tf.frameConfig.CellGameCfgs[tf.nodeIndex]
			if listenCfg.ListenMode == "unix" {
				listenMode = evhub.ListenModeUnix
			} else {
				listenMode = evhub.ListenModeTcp
			}
			listenAddr = listenCfg.ListenAddr
			break
		}
	case trnode.ETRNodeTypeCellView:
		{
			listenCfg := tf.frameConfig.CellViewCfgs[tf.nodeIndex]
			if listenCfg.ListenMode == "unix" {
				listenMode = evhub.ListenModeUnix
			} else {
				listenMode = evhub.ListenModeTcp
			}
			listenAddr = listenCfg.ListenAddr
			break
		}
	case trnode.ETRNodeTypeCellLog:
		{
			listenCfg := tf.frameConfig.CellLogCfgs[tf.nodeIndex]
			if listenCfg.ListenMode == "unix" {
				listenMode = evhub.ListenModeUnix
			} else {
				listenMode = evhub.ListenModeTcp
			}
			listenAddr = listenCfg.ListenAddr
			break
		}
	case trnode.ETRNodeTypeCellSocial:
		{
			listenCfg := tf.frameConfig.CellSocialCfgs[tf.nodeIndex]
			if listenCfg.ListenMode == "unix" {
				listenMode = evhub.ListenModeUnix
			} else {
				listenMode = evhub.ListenModeTcp
			}
			listenAddr = listenCfg.ListenAddr
			break
		}
	case trnode.ETRNodeTypeCellTrans:
		{
			listenCfg := tf.frameConfig.CellTransCfgs[tf.nodeIndex]
			if listenCfg.ListenMode == "unix" {
				listenMode = evhub.ListenModeUnix
			} else {
				listenMode = evhub.ListenModeTcp
			}
			listenAddr = listenCfg.ListenAddr
			break
		}
	default:
		{
			loghlp.Warnf("not find listen config,nodeType:%d", tf.nodeType)
			break
		}
	}
	if len(listenAddr) == 0 {
		return errors.New("not find listenaddr")
	}
	loghlp.Infof("listen addr:%s", listenAddr)
	return tf.evHub.Listen(listenMode, listenAddr)
}

func (tf *TRFrame) makeInitWorkNode(nType int32, index int32) ITRFrameWorkNode {
	switch nType {
	// 分布式cell
	case trnode.ETRNodeTypeCellGate:
		{
			return tcellgate.New(tf, index)
		}
	case trnode.ETRNodeTypeCellRoot:
		{
			return tcellroot.New(tf, index)
		}
	case trnode.ETRNodeTypeCellData:
		{
			return tcelldata.New(tf, index)
		}
	case trnode.ETRNodeTypeCellCenter:
		{
			return tcellcenter.New(tf, index)
		}
	case trnode.ETRNodeTypeCellLogic:
		{
			return tcelllogic.New(tf, index)
		}
	case trnode.ETRNodeTypeCellGame:
		{
			return tcellgame.New(tf, index)
		}
	case trnode.ETRNodeTypeCellView:
		{
			return tcellview.New(tf, index)
		}
	case trnode.ETRNodeTypeCellLog:
		{
			return tcelllog.New(tf, index)
		}
	case trnode.ETRNodeTypeCellSocial:
		{
			return tcellsocial.New(tf, index)
		}
	case trnode.ETRNodeTypeCellTrans:
		{
			return tcelltrans.New(tf, index)
		}
	default:
		{
			loghlp.Info("unknonw init worknode type:%d", nType)
			break
		}
	}
	return nil
}
func (tf *TRFrame) initNodeSetting() {
	switch tf.nodeType {
	case trnode.ETRNodeTypeCellRoot:
		{
			nodeCfg := tf.frameConfig.CellRootCfgs[tf.nodeIndex]
			loghlp.SetConsoleLogLevel(logrus.Level(nodeCfg.LogLevel))
			if len(nodeCfg.LogPath) > 0 {
				loghlp.ActiveFileLog(nodeCfg.LogPath, fmt.Sprintf("cellroot%d", tf.nodeIndex))
				loghlp.SetFileLogLevel(logrus.Level(nodeCfg.LogLevel))
			}
		}
	case trnode.ETRNodeTypeCellGate:
		{
			nodeCfg := tf.frameConfig.CellGateCfgs[tf.nodeIndex]
			loghlp.SetConsoleLogLevel(logrus.Level(nodeCfg.LogLevel))
			if len(nodeCfg.LogPath) > 0 {
				loghlp.ActiveFileLog(nodeCfg.LogPath, fmt.Sprintf("cellgate%d", tf.nodeIndex))
				loghlp.SetFileLogLevel(logrus.Level(nodeCfg.LogLevel))
			}
		}
	case trnode.ETRNodeTypeCellData:
		{
			nodeCfg := tf.frameConfig.CellDataCfgs[tf.nodeIndex]
			loghlp.SetConsoleLogLevel(logrus.Level(nodeCfg.LogLevel))
			if len(nodeCfg.LogPath) > 0 {
				loghlp.ActiveFileLog(nodeCfg.LogPath, fmt.Sprintf("celldata%d", tf.nodeIndex))
				loghlp.SetFileLogLevel(logrus.Level(nodeCfg.LogLevel))
			}
		}
	case trnode.ETRNodeTypeCellCenter:
		{
			nodeCfg := tf.frameConfig.CellCenterCfgs[tf.nodeIndex]
			loghlp.SetConsoleLogLevel(logrus.Level(nodeCfg.LogLevel))
			if len(nodeCfg.LogPath) > 0 {
				loghlp.ActiveFileLog(nodeCfg.LogPath, fmt.Sprintf("cellcenter%d", tf.nodeIndex))
				loghlp.SetFileLogLevel(logrus.Level(nodeCfg.LogLevel))
			}
		}
	case trnode.ETRNodeTypeCellLogic:
		{
			nodeCfg := tf.frameConfig.CellLogicCfgs[tf.nodeIndex]
			loghlp.SetConsoleLogLevel(logrus.Level(nodeCfg.LogLevel))
			if len(nodeCfg.LogPath) > 0 {
				loghlp.ActiveFileLog(nodeCfg.LogPath, fmt.Sprintf("celllogic%d", tf.nodeIndex))
				loghlp.SetFileLogLevel(logrus.Level(nodeCfg.LogLevel))
			}
		}
	case trnode.ETRNodeTypeCellGame:
		{
			nodeCfg := tf.frameConfig.CellGameCfgs[tf.nodeIndex]
			loghlp.SetConsoleLogLevel(logrus.Level(nodeCfg.LogLevel))
			if len(nodeCfg.LogPath) > 0 {
				loghlp.ActiveFileLog(nodeCfg.LogPath, fmt.Sprintf("cellgame%d", tf.nodeIndex))
				loghlp.SetFileLogLevel(logrus.Level(nodeCfg.LogLevel))
			}
		}
	case trnode.ETRNodeTypeCellView:
		{
			nodeCfg := tf.frameConfig.CellViewCfgs[tf.nodeIndex]
			loghlp.SetConsoleLogLevel(logrus.Level(nodeCfg.LogLevel))
			if len(nodeCfg.LogPath) > 0 {
				loghlp.ActiveFileLog(nodeCfg.LogPath, fmt.Sprintf("cellview%d", tf.nodeIndex))
				loghlp.SetFileLogLevel(logrus.Level(nodeCfg.LogLevel))
			}
		}
	case trnode.ETRNodeTypeCellLog:
		{
			nodeCfg := tf.frameConfig.CellLogCfgs[tf.nodeIndex]
			loghlp.SetConsoleLogLevel(logrus.Level(nodeCfg.LogLevel))
			if len(nodeCfg.LogPath) > 0 {
				loghlp.ActiveFileLog(nodeCfg.LogPath, fmt.Sprintf("celllog%d", tf.nodeIndex))
				loghlp.SetFileLogLevel(logrus.Level(nodeCfg.LogLevel))
			}
		}
	case trnode.ETRNodeTypeCellSocial:
		{
			nodeCfg := tf.frameConfig.CellSocialCfgs[tf.nodeIndex]
			loghlp.SetConsoleLogLevel(logrus.Level(nodeCfg.LogLevel))
			if len(nodeCfg.LogPath) > 0 {
				loghlp.ActiveFileLog(nodeCfg.LogPath, fmt.Sprintf("cellsocial%d", tf.nodeIndex))
				loghlp.SetFileLogLevel(logrus.Level(nodeCfg.LogLevel))
			}
		}
	case trnode.ETRNodeTypeCellTrans:
		{
			nodeCfg := tf.frameConfig.CellTransCfgs[tf.nodeIndex]
			loghlp.SetConsoleLogLevel(logrus.Level(nodeCfg.LogLevel))
			if len(nodeCfg.LogPath) > 0 {
				loghlp.ActiveFileLog(nodeCfg.LogPath, fmt.Sprintf("celltrans%d", tf.nodeIndex))
				loghlp.SetFileLogLevel(logrus.Level(nodeCfg.LogLevel))
			}
		}
	default:
		{
			loghlp.Errorf("unknonw init worknode type:%d", tf.nodeType)
			break
		}
	}
}
func (tf *TRFrame) stopFrame() {
	if tf.runStep < ETRFrameStepStop {
		tf.runStep = ETRFrameStepStop
	}
}
func (tf *TRFrame) GetEvHub() *evhub.EventHub {
	return tf.evHub
}
func (tf *TRFrame) GetFrameConfig() *tframeconfig.FrameConfig {
	return tf.frameConfig
}
func (tf *TRFrame) AfterMsgJob(doJob func()) {
	tf.msgDoneList = append(tf.msgDoneList, doJob)
}
func (tf *TRFrame) updateKeepNodeAlive(curTimeMs int64) {
	if curTimeMs-tf.keepNodeTime < 3000 {
		return
	}
	tf.keepNodeTime = curTimeMs
	// 心跳
	for _, tn := range tf.frameNodeMgr.nodeList {
		if !tn.IsServerClient() {
			continue
		}
		if curTimeMs/1000-tn.LastHeartTime() >= 15 {
			sendMsg := MakeInnerMsg(protocol.EMsgClassFrame, protocol.EFrameMsgKeepNodeHeart, make([]byte, 0))
			cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
				loghlp.Infof("keep node heart callback succ:%d", okCode)
			}
			callInfo := tf.remoteMsgMgr.makeCallInfo(protocol.EMsgClassFrame,
				protocol.EFrameMsgKeepNodeHeart,
				cb,
				nil)
			setSecondHead(sendMsg, 0, callInfo.reqID, 0)
			tn.SendMsg(sendMsg)
			tn.SetHeartTime(curTimeMs / 1000)
		}
	}
}

func (tf *TRFrame) GetCurNodeType() int32 {
	return tf.curWorkNode.NodeType()
}

func (tf *TRFrame) registerFrameHandler() {
	tf.msgDispatcher.RegisterMsgHandler(
		protocol.EMsgClassFrame,
		protocol.EFrameMsgRegisterServerInfo,
		handleRegisterNodeInfo)

	tf.msgDispatcher.RegisterMsgHandler(
		protocol.EMsgClassFrame,
		protocol.EFrameMsgTransMsg,
		handleTransNodeMsg)
}
