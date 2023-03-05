/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-16 17:57:09
 * @LastEditTime: 2022-09-16 17:59:04
 * @FilePath: \trcell\pkg\trframe\trnode\tcellroot\trnode_rootserver_init_connnect.go
 */
package tcellroot

import (
	"fmt"
	"time"
	"trcell/pkg/evhub"
	"trcell/pkg/loghlp"
	"trcell/pkg/pb/pbframe"
	"trcell/pkg/protocol"
	"trcell/pkg/trframe/iframe"
	"trcell/pkg/trframe/trnode"
)

func (frameNode *FrameNodeRootServer) InitConnectServer() bool {
	// // 连接transserver
	frameConfig := frameNode.tframeObj.GetFrameConfig()
	if !frameConfig.ConnectTrans {
		loghlp.Info("config set not connect to trans")
		return true
	}
	evHub := frameNode.tframeObj.GetEvHub()
	for _, cfg := range frameConfig.CellTransCfgs {
		connDo := func() {
			// 发送注册消息
			reqMsg := &pbframe.FrameMsgRegisterServerInfoReq{
				ZoneID:    frameNode.tframeObj.GetFrameConfig().ZoneID,
				NodeType:  trnode.ETRNodeTypeCellRoot,
				NodeIndex: frameNode.nodeIndex,
				NodeDes:   fmt.Sprintf("%d_ETRNodeTypeCellRoot%d", frameNode.tframeObj.GetFrameConfig().ZoneID, frameNode.nodeIndex),
			}
			cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
				loghlp.Infof("root register to celltrans callback succ:%d", okCode)
			}
			nodeUid := trnode.GenNodeUid(cfg.ZoneID, trnode.ETRNodeTypeCellTrans, cfg.NodeIndex)
			frameNode.tframeObj.ForwardNodePBMessageByNodeUid(
				protocol.EMsgClassFrame,
				protocol.EFrameMsgRegisterServerInfo,
				reqMsg,
				nodeUid,
				cb,
				nil,
			)
		}
		userData := &iframe.SessionUserData{
			ZoneID:         cfg.ZoneID, // 指定zoneID,和普通节点不同
			DataType:       iframe.ESessionDataTypeNetInfo,
			NodeType:       trnode.ETRNodeTypeCellTrans,
			NodeIndex:      cfg.NodeIndex,
			DesInfo:        fmt.Sprintf("%dETRNodeTypeCellTrans%d", frameNode.tframeObj.GetFrameConfig().ZoneID, cfg.NodeIndex),
			IsServerClient: true,
		}
		failCount := 0
		var listenMode int32 = evhub.ListenModeTcp
		if cfg.ListenMode == "unix" {
			listenMode = evhub.ListenModeUnix
		}
		for {
			if !evHub.Connect2(listenMode, cfg.ListenAddr, true, userData, connDo) {
				failCount++
				loghlp.Warnf("connect to transserver fail")
			} else {
				break
			}
			if failCount >= 10 {
				loghlp.Errorf("connect transserver fail,exit")
				return false
			}
			time.Sleep(time.Second * 1)
		}
	}
	time.Sleep(time.Second)
	return true
}
