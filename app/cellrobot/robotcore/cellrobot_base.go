package robotcore

import (
	"fmt"
	"io"
	"net"
	"strings"
	"time"
	"trcell/pkg/crossdef"
	"trcell/pkg/loghlp"
	"trcell/pkg/pb/pbclient"
	"trcell/pkg/pb/pbtools"
	"trcell/pkg/protocol"
	"trcell/pkg/sconst"
	"trcell/pkg/timeutil"
	"trcell/pkg/webreq"

	"github.com/goccy/go-json"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

type ICellRobot interface {
	GetRobotName() string
	LogRecvMsgInfo(clientMsg *ClientMessage, pbMsg proto.Message)
}
type RobotMsgHandler func(robotIns ICellRobot, clientMsg *ClientMessage)

type RobotCallEnv struct {
	CallbackFunc func(clientMsg *ClientMessage)
	BeginTime    int64
	MsgClass     int32
	MsgType      int32
}
type IRobotAI interface {
	Update(curTime int64)
}
type CellRobot struct {
	tcpPeerConn net.Conn
	recvMsgCh   chan *ClientMessage // 收到的消息队列
	//sendMsgCh   chan *evhub.NetMessage // 发送消息的队列
	RobotName    string
	UserID       int64
	Icon         int32
	SeqID        int32
	ClosedWrite  bool
	HeartTime    int64
	StopRun      bool
	MsgHandlers  map[int32]RobotMsgHandler
	AsyncCall    map[int32]*RobotCallEnv
	RoomID       int64 // 当前房间id
	TargetRoomID int64 // 指定进入的房间id
	UpdateFunc   func(curTime int64)
	// status
	RobotStatus  int32
	StatusTime   int64
	StatusStep   int32 // 状态阶段	0-进入 1-状态完成
	Token        string
	RealRobotIns ICellRobot
	//
	AiInstance IRobotAI
	// // data
	GateAddr string
	// RoomDetail *pbclient.RoomData
	// IsReady    bool // 是否准备游戏了
	// robot action
	LastActionTime int64  // 最近执行的行为时间
	ActionName     string // 行为名字标示
	ActionTime     int64  // 行为开始时间 0-没有执行任何行为
	ActionStep     int32  // 行为阶段 0--进行中 1-完成结束了
	ActionSeq      int32
}

func NewCellRobot(robotName string) *CellRobot {
	r := &CellRobot{
		recvMsgCh:    make(chan *ClientMessage),
		RobotName:    robotName,
		UserID:       0,
		SeqID:        0,
		ClosedWrite:  false,
		HeartTime:    0,
		StopRun:      false,
		MsgHandlers:  make(map[int32]RobotMsgHandler),
		AsyncCall:    make(map[int32]*RobotCallEnv),
		RoomID:       0,
		RealRobotIns: nil,
		ActionTime:   0,
		ActionStep:   0,
		ActionSeq:    0,
	}
	//r.initRegisterRobotHandle()
	return r
}

func (robotObj *CellRobot) ConnectServer(servAddr string) error {
	// u := url.URL{Scheme: "ws", Host: servAddr, Path: "/ws"}
	conn, err := net.Dial("tcp", servAddr)
	if err != nil {
		logrus.Error("connect server error:" + err.Error())
		return err
	}
	robotObj.tcpPeerConn = conn
	go func() {
		robotObj.startRead()
	}()

	return nil
}

// 日志操作
func (robotObj *CellRobot) Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}
func (robotObj *CellRobot) Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}
func (robotObj *CellRobot) Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}
func (robotObj *CellRobot) Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

func (robotObj *CellRobot) startRead() {
	headBytes := make([]byte, sconst.ClientMsgHeadSize)
	for {
		robotObj.tcpPeerConn.SetReadDeadline(time.Now().Add(time.Second * 40))
		_, err := io.ReadFull(robotObj.tcpPeerConn, headBytes)
		if err != nil {
			var errInfo string = err.Error()
			if strings.Contains(errInfo, "closed by the remote host") {
				loghlp.Warnf("remote(%s)closed!!!", robotObj.tcpPeerConn.RemoteAddr(), robotObj.UserID)
			} else {
				loghlp.Errorf("tcpRecv error:%s", err.Error())
			}
			break
		}
		var clientHeader = NewClientHeader(0, 0, 0)
		clientHeader.Decode(headBytes)
		// 消息体
		var clientMsg = NewClientMessage(clientHeader)
		bodyLen := clientHeader.Len - ClientHeaderSize
		clientMsg.Data = make([]byte, bodyLen)
		if bodyLen > 0 {
			_, err = io.ReadFull(robotObj.tcpPeerConn, clientMsg.Data)
			if err != nil {
				var errInfo string = err.Error()
				if strings.Contains(errInfo, "closed by the remote host") {
					loghlp.Warnf("remote client(%s)[%d] closed!!!", robotObj.tcpPeerConn.RemoteAddr(), robotObj.UserID)
				} else {
					loghlp.Errorf("tcpRecv body error:%s", err.Error())
				}
				break
			}
		}
		// 特殊消息hook转换一下
		if clientMsg.Head.MainType == protocol.ECMsgClassBase && clientMsg.Head.SubType == protocol.ECMsgBaseWrapOpt {
			wrapMsg := &pbclient.ECMsgBaseWrapOptRsp{}
			errWrap := proto.Unmarshal(clientMsg.Data, wrapMsg)
			if errWrap == nil {
				clientMsg.Head.MainType = uint8(wrapMsg.MsgClass)
				clientMsg.Head.SubType = uint8(wrapMsg.MsgType)
				clientMsg.Head.Result = uint16(wrapMsg.ErrCode)
				clientMsg.Head.ID = uint32(wrapMsg.AckID)
				clientMsg.Data = wrapMsg.Data
				clientMsg.Head.Len = ClientHeaderSize + uint32(len(wrapMsg.Data))
			}
		}

		robotObj.recvMsgCh <- clientMsg
	}
	robotObj.ClosedWrite = true
	robotObj.Infof("robot(%s) exit read", robotObj.RobotName)
}

func (robotObj *CellRobot) GenActionSeq() int32 {
	robotObj.ActionSeq++
	if robotObj.ActionSeq >= 99999999 {
		robotObj.ActionSeq = 1
	}
	return robotObj.ActionSeq
}

// 执行一个行为(一连串的消息),设定,一个机器人，同时只能做一件事情
func (robotObj *CellRobot) DoAction(actionName string, doAnything func()) bool {
	if robotObj.ActionTime != 0 {
		robotObj.Errorf("robot(%s) do action fail, action(%s) is doing", robotObj.RobotName, robotObj.ActionName)
		return false
	}
	// 可以执行
	robotObj.ActionName = fmt.Sprintf("%s_%d", actionName, robotObj.GenActionSeq())
	robotObj.ActionTime = timeutil.NowTime()
	robotObj.ActionStep = 0
	robotObj.Infof("robot(%s) start action(%s)", robotObj.RobotName, robotObj.ActionName)
	doAnything()
	return true
}

func (robotObj *CellRobot) EndAction() {
	robotObj.Infof("robot(%s) end action(%s)", robotObj.RobotName, robotObj.ActionName)
	robotObj.LastActionTime = robotObj.ActionTime
	robotObj.ActionTime = 0
	robotObj.ActionStep = 0
	robotObj.ActionName = ""
}

func (robotObj *CellRobot) GenSeqID() int32 {
	robotObj.SeqID++
	if robotObj.SeqID >= 999999999 {
		robotObj.SeqID = 1
	}
	return robotObj.SeqID
}
func (robotObj *CellRobot) SendMsgToServer(msgClass int32, msgType int32, pbMsg proto.Message) {
	if robotObj.ClosedWrite {
		robotObj.Warnf("writeclosed, robot(%s) ingnore send msg(%d_%d)", robotObj.RobotName, msgClass, msgType)
		return
	}

	clientMsg := MakeClientTcpMessage(uint8(msgClass),
		uint8(msgType),
		0,
		0)

	clientMsg.Data, _ = proto.Marshal(pbMsg)
	sendData := clientMsg.Encode()
	robotObj.tcpPeerConn.SetWriteDeadline(time.Now().Add(time.Second * 40))
	sendN, errSend := robotObj.tcpPeerConn.Write(sendData)
	if sendN != len(sendData) {
		loghlp.Errorf("sendToClient sendN(%d) != len(sendData)(%d)", sendN, len(sendData))
	}
	if errSend != nil {
		loghlp.Errorf("sendToServer error:%s", errSend.Error())
	}
	robotObj.Debugf("robot(%s) send msg(%s_%s)[%s]:%+v", robotObj.RobotName, msgClass, msgType, pbtools.GetFullNameByMessage(pbMsg), pbMsg)
}

func (robotObj *CellRobot) LogCbMsgInfo(clientMsg *ClientMessage, pbMsg proto.Message) {
	robotObj.Debugf("log robot(%s) cbmsg(%d_%d)ack_seqid[%d][%s]isok[%d]",
		robotObj.RobotName,
		clientMsg.Head.MainType,
		clientMsg.Head.SubType,
		clientMsg.Head.ID,
		pbtools.GetFullNameByMessage(pbMsg),
		clientMsg.Head.Result,
	)
}
func (robotObj *CellRobot) LogRecvMsgInfo(clientMsg *ClientMessage, pbMsg proto.Message) {
	robotObj.Debugf("log robot(%s) recv msg(%d_%d)[%d][%s]isok[%d]",
		robotObj.RobotName,
		clientMsg.Head.MainType,
		clientMsg.Head.SubType,
		clientMsg.Head.ID,
		pbtools.GetFullNameByMessage(pbMsg),
		clientMsg.Head.Result,
	)
}
func (robotObj *CellRobot) SendKeepHeart() {
	reqMsg := &pbclient.ECMsgPlayerKeepHeartReq{}
	robotObj.RemoteCall(protocol.ECMsgClassPlayer, protocol.ECMsgPlayerKeepHeart, reqMsg, func(clientMsg *ClientMessage) {
		rsp := &pbclient.ECMsgPlayerKeepHeartRsp{}
		errParse := proto.Unmarshal(clientMsg.Data, rsp)
		if errParse != nil {
			return
		}
		robotObj.LogCbMsgInfo(clientMsg, rsp)
	})
}
func (robotObj *CellRobot) GetRobotName() string {
	return robotObj.RobotName
}
func (robotObj *CellRobot) Update(curTime int64) {
	if robotObj.AiInstance != nil {
		robotObj.AiInstance.Update(curTime)
	}
}
func (robotObj *CellRobot) Run() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for !robotObj.StopRun {
		select {
		case <-ticker.C:
			{
				curTime := time.Now().Unix()
				if curTime-robotObj.HeartTime >= 3 {
					// 发送心跳
					robotObj.SendKeepHeart()
					robotObj.HeartTime = curTime
				}
				// 异步调用超时处理
				for k, v := range robotObj.AsyncCall {
					if curTime-v.BeginTime >= 30 {
						robotObj.Errorf("asynccall(%d_%d),seqid(%d) timeout", v.MsgClass, v.MsgType, k)
						delete(robotObj.AsyncCall, k)
					}
				}
				// 机器人行为时间检测
				if robotObj.ActionTime > 0 {
					if robotObj.ActionStep == 0 {
						if curTime-robotObj.ActionTime > 30 {
							// 行为超时了,异常,解除锁定
							robotObj.LastActionTime = robotObj.ActionTime
							robotObj.ActionTime = 0
							robotObj.ActionStep = 0
						}
					}
				}
				robotObj.Update(curTime)
				break
			}
		case clientMsg, ok := <-robotObj.recvMsgCh:
			{
				if ok {
					robotObj.Debugf("recv Message(%d_%d) seqid(%d)", clientMsg.Head.MainType, clientMsg.Head.SubType, clientMsg.Head.ID)
					robotObj.DispatchRobotMsg(clientMsg)
				} else {
					robotObj.Errorf("maybe recvMsgCh closed !!!")
				}
				break
			}
		}
	}

	robotObj.Infof("robot(%s) stop", robotObj.RobotName)

	close(robotObj.recvMsgCh)
	time.Sleep(time.Second * 2)
	robotObj.Infof("robot(%s) exit run", robotObj.RobotName)
}

func (robotObj *CellRobot) ChangeRobotStatus(robotStatus int32, step int32) {
	robotObj.RobotStatus = robotStatus
	robotObj.StatusTime = timeutil.NowTime()
	robotObj.StatusStep = step
}

func (robotObj *CellRobot) RemoteCall(msgClass int32, msgType int32, pbMsg proto.Message, cbFunc func(clientMsg *ClientMessage)) {
	if robotObj.ClosedWrite {
		robotObj.Warnf("writeclosed, robot(%s) remotecall ingnore send msg(%d_%d)", robotObj.RobotName, msgClass, msgType)
		return
	}
	clientMsg := MakeClientTcpMessage(uint8(msgClass),
		uint8(msgType),
		uint32(robotObj.GenSeqID()),
		0)
	robotObj.AsyncCall[int32(clientMsg.Head.ID)] = &RobotCallEnv{
		BeginTime:    timeutil.NowTime(),
		CallbackFunc: cbFunc,
		MsgClass:     msgClass,
		MsgType:      msgType,
	}

	clientMsg.Data, _ = proto.Marshal(pbMsg)

	// 消息适配一下
	wrapMsg := &pbclient.ECMsgBaseWrapOptReq{}
	wrapMsg.MsgClass = msgClass
	wrapMsg.MsgType = msgType
	wrapMsg.ReqID = int32(clientMsg.Head.ID)
	wrapMsg.Data = clientMsg.Data
	wrapData, wrapErr := proto.Marshal(wrapMsg)
	if wrapErr == nil {
		clientMsg.Head.MainType = protocol.ECMsgClassBase
		clientMsg.Head.SubType = protocol.ECMsgBaseWrapOpt
		clientMsg.Data = wrapData
	}

	sendData := clientMsg.Encode()
	robotObj.tcpPeerConn.SetWriteDeadline(time.Now().Add(time.Second * 5))
	sendN, errSend := robotObj.tcpPeerConn.Write(sendData)
	if sendN != len(sendData) {
		loghlp.Errorf("sendToClient sendN(%d) != len(sendData)(%d)", sendN, len(sendData))
	}
	if errSend != nil {
		loghlp.Errorf("sendToServer error:%s", errSend.Error())
	}

	robotObj.Debugf("robot(%s) send remotecall msg(%d_%d)(%s)seqid(%d),userid(%d):%+v",
		robotObj.RobotName,
		msgClass,
		msgType,
		pbtools.GetFullNameByMessage(pbMsg),
		clientMsg.Head.ID,
		robotObj.UserID,
		pbMsg)
}

// 注册账号
func (robotObj *CellRobot) RegisterAccount(hostAddr string) bool {
	regReq := &AccountRegisterReq{
		UserName: robotObj.RobotName,
		Pswd:     "123456",
	}
	repData, err := webreq.PostJson(fmt.Sprintf("http://%s/account/register", hostAddr), regReq)
	if err != nil {
		robotObj.Errorf("register account(%s) error:%s", robotObj.RobotName, err.Error())
		return false
	}
	type RegisterMsgRsp struct {
		Data *AccountRegisterRsp
		Code int32
		Msg  string
	}
	regRep := &RegisterMsgRsp{}
	errJv := json.Unmarshal(repData, regRep)
	if errJv != nil {
		robotObj.Errorf("unmarshal RegisterMsgRsp error:%s", errJv.Error())
		return false
	}
	robotObj.UserID = regRep.Data.UserID
	robotObj.Infof("robot(%s) register success,user_id(%d)", robotObj.RobotName, robotObj.UserID)
	return true
}

// 服务器列表
func (robotObj *CellRobot) GetServerList(hostAddr string) bool {
	regReq := &QueryServerListReq{}
	repData, err := webreq.PostJson(fmt.Sprintf("http://%s/server/list", hostAddr), regReq)
	if err != nil {
		robotObj.Errorf("server list(%s) error:%s", robotObj.RobotName, err.Error())
		return false
	}
	type ServerListMsgRsp struct {
		Data *QueryServerListRsp
		Code int32
		Msg  string
	}
	serverlistRep := &ServerListMsgRsp{}
	errJv := json.Unmarshal(repData, serverlistRep)
	if errJv != nil {
		robotObj.Errorf("unmarshal RegisterMsgRsp error:%s", errJv.Error())
		return false
	}
	if serverlistRep.Code != protocol.ECodeSuccess {
		robotObj.Errorf("get serverList error(%d)", serverlistRep.Code)
		return false
	}
	// 这里默认用第一个
	if len(serverlistRep.Data.ServerList) > 0 {
		robotObj.GateAddr = fmt.Sprintf("%s:%d", serverlistRep.Data.ServerList[0].Addr,
			serverlistRep.Data.ServerList[0].Port)
	}

	return true
}

// 登录账号
func (robotObj *CellRobot) LoginAccount(hostAddr string) bool {
	regReq := &AccountLoginReq{
		UserName: robotObj.RobotName,
		Pswd:     "123456",
	}
	repData, err := webreq.PostJson(fmt.Sprintf("http://%s/account/login", hostAddr), regReq)
	if err != nil {
		robotObj.Errorf("login account(%s) error:%s", robotObj.RobotName, err.Error())
		return false
	}
	type AccountLoginRsp struct {
		ServerAddr string `json:"server_addr"` // 大厅地址
		Token      string `json:"token"`       // 返回token信息
		Statu      int32  `json:"status"`      // 当前账号状态 0-未验证 1-认证通过
		RestTime   int32  `json:"rest_time"`   // 剩余的认证时间,为0则不可用
	}
	type LoginMsgRsp struct {
		Data *AccountLoginRsp
		Code int32
		Msg  string
	}
	loginRep := &LoginMsgRsp{}
	errJv := json.Unmarshal(repData, loginRep)
	if errJv != nil {
		robotObj.Errorf("unmarshal RegisterMsgRsp error:%s", errJv.Error())
		return false
	}
	if loginRep.Code != protocol.ECodeSuccess {
		robotObj.Errorf("login account error(%d)", loginRep.Code)
		return false
	}
	// 验证touken
	ok, tokenRes := crossdef.TokenAuthClaims(loginRep.Data.Token, crossdef.SignKey)
	if !ok {
		robotObj.Errorf("token parse fail")
		return false
	} else {
		loghlp.Infof("parse player token success:%+v", *tokenRes)
	}
	//robotObj.GateAddr = loginRep.Data.ServerAddr
	robotObj.UserID = tokenRes.UserID
	robotObj.Token = loginRep.Data.Token
	robotObj.Infof("robot(%s) login account success,user_id(%d):%+v", robotObj.RobotName, robotObj.UserID, loginRep.Data)
	return true
}

// 登录游戏
func (robotObj *CellRobot) LoginHall() {
	// robotObj.ChangeRobotStatus(ERobotStatusLoginHall, ERobotStatusStepIng)
	// reqMsg := &pbclient.ECMsgPlayerLoginHallReq{
	// 	Token: robotObj.Token,
	// }
	// robotObj.RemoteCall(protocol.ECMsgClassPlayer,
	// 	protocol.ECMsgPlayerLoginHall,
	// 	reqMsg,
	// 	func(sMsg *evhub.NetMessage) {
	// 		rsp := &pbclient.ECMsgPlayerLoginHallRsp{}
	// 		if !trframe.DecodePBMessage(sMsg, rsp) {
	// 			return
	// 		}
	// 		robotObj.LogCbMsgInfo(sMsg, rsp)
	// 		if sMsg.Head.Result == protocol.ECodeSuccess {
	// 			robotObj.Debugf("robot(%s) login hall succ", robotObj.RobotName)
	// 			robotObj.Icon = rsp.RoleData.Icon
	// 			robotObj.ChangeRobotStatus(ERobotStatusLoginHall, ERobotStatusStepFinish)
	// 		}
	// 	},
	// )
}

// 进入游戏
func (robotObj *CellRobot) SendEnterRoom(roomID int64) {
	// reqMsg := &pbclient.ECMsgRoomEnterReq{
	// 	RoomID: roomID,
	// }
	// robotObj.ChangeRobotStatus(ERobotStatusEnterRoom, ERobotStatusStepIng)
	// robotObj.RemoteCall(protocol.ECMsgClassRoom,
	// 	protocol.ECMsgRoomEnter,
	// 	reqMsg,
	// 	func(sMsg *evhub.NetMessage) {
	// 		rsp := &pbclient.ECMsgRoomEnterRsp{}
	// 		if !trframe.DecodePBMessage(sMsg, rsp) {
	// 			return
	// 		}
	// 		robotObj.LogCbMsgInfo(sMsg, rsp)
	// 		if sMsg.Head.Result == protocol.ECodeSuccess {
	// 			robotObj.Debugf("robot(%s) enter room(%d) succ:%+v", robotObj.RobotName, roomID, rsp.RoomDetail)
	// 			robotObj.RoomDetail = rsp.RoomDetail
	// 			robotObj.RoomID = rsp.RoomDetail.RoomID
	// 			robotObj.ChangeRobotStatus(ERobotStatusEnterRoom, ERobotStatusStepFinish)
	// 		}
	// 	},
	// )
}
func (robotObj *CellRobot) CreateRole() {
	reqMsg := &pbclient.ECMsgPlayerCreateRoleReq{
		Token:    robotObj.Token,
		Nickname: robotObj.RobotName,
	}
	robotObj.RemoteCall(protocol.ECMsgClassPlayer, protocol.ECMsgPlayerCreateRole, reqMsg, func(clientMsg *ClientMessage) {
		rsp := &pbclient.ECMsgPlayerCreateRoleRsp{}
		errParse := proto.Unmarshal(clientMsg.Data, rsp)
		if errParse != nil {
			return
		}
		robotObj.LogCbMsgInfo(clientMsg, rsp)
		if clientMsg.Head.Result == protocol.ECodeSuccess {
			robotObj.Infof("%s create role succ", robotObj.RobotName)
		}
	})
}
func (robotObj *CellRobot) LoginGame() {
	reqMsg := &pbclient.ECMsgPlayerLoginGameReq{
		Token: robotObj.Token,
	}
	robotObj.RemoteCall(protocol.ECMsgClassPlayer, protocol.ECMsgPlayerLoginGame, reqMsg, func(clientMsg *ClientMessage) {
		rsp := &pbclient.ECMsgPlayerLoginGameRsp{}
		errParse := proto.Unmarshal(clientMsg.Data, rsp)
		if errParse != nil {
			return
		}
		robotObj.LogCbMsgInfo(clientMsg, rsp)
		if clientMsg.Head.Result == protocol.ECodeRoleNotExisted {
			robotObj.CreateRole()
		} else if clientMsg.Head.Result == protocol.ECodeSuccess {
			robotObj.Infof("[%s] login game succ,roleName:%s", robotObj.RobotName, rsp.RoleData.Nickname)
		}
	})
}
