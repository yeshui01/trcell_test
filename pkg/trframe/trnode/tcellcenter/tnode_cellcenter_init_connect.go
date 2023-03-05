/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 11:36:31
 * @LastEditTime: 2022-09-27 18:26:01
 * @FilePath: \trcell\pkg\trframe\trnode\tcellcenter\tnode_cellcenter_init_connect.go
 */
package tcellcenter

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

func (frameNode *FrameNodeCellCenter) InitConnectServer() bool {
	// 连接rootserver
	frameConfig := frameNode.tframeObj.GetFrameConfig()
	evHub := frameNode.tframeObj.GetEvHub()
	for idx, cfg := range frameConfig.CellRootCfgs {
		connDo := func() {
			// 发送注册消息
			reqMsg := &pbframe.FrameMsgRegisterServerInfoReq{
				ZoneID:    frameNode.tframeObj.GetFrameConfig().ZoneID,
				NodeType:  trnode.ETRNodeTypeCellCenter,
				NodeIndex: frameNode.nodeIndex,
				NodeDes:   fmt.Sprintf("%d_ETRNodeTypeCellCenter%d", frameNode.tframeObj.GetFrameConfig().ZoneID, frameNode.nodeIndex),
			}
			cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
				loghlp.Infof("center register to cellroot callback succ:%d", okCode)
			}
			frameNode.tframeObj.ForwardMessage(
				protocol.EMsgClassFrame,
				protocol.EFrameMsgRegisterServerInfo,
				reqMsg,
				trnode.ETRNodeTypeCellRoot,
				int32(idx),
				cb,
				nil,
			)
		}
		userData := &iframe.SessionUserData{
			DataType:       iframe.ESessionDataTypeNetInfo,
			NodeType:       trnode.ETRNodeTypeCellRoot,
			NodeIndex:      int32(idx),
			DesInfo:        fmt.Sprintf("%dETRNodeTypeCellRoot%d", frameNode.tframeObj.GetFrameConfig().ZoneID, idx),
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
				loghlp.Warnf("connect to rootserver fail")
			} else {
				break
			}
			if failCount >= 10 {
				loghlp.Errorf("connect rootserver fail,exit")
				return false
			}
			time.Sleep(time.Second * 1)
		}
	}
	time.Sleep(time.Second)
	// 连接dataserver
	for idx, cfg := range frameConfig.CellDataCfgs {
		connDo := func() {
			// 发送注册消息
			reqMsg := &pbframe.FrameMsgRegisterServerInfoReq{
				ZoneID:    frameNode.tframeObj.GetFrameConfig().ZoneID,
				NodeType:  trnode.ETRNodeTypeCellCenter, // 当前自身节点
				NodeIndex: frameNode.nodeIndex,
				NodeDes:   fmt.Sprintf("%d_ETRNodeTypeCellCenter%d", frameNode.tframeObj.GetFrameConfig().ZoneID, frameNode.nodeIndex),
			}
			cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
				loghlp.Infof("center register to celldata callback succ:%d", okCode)
			}
			frameNode.tframeObj.ForwardMessage(
				protocol.EMsgClassFrame,
				protocol.EFrameMsgRegisterServerInfo,
				reqMsg,
				trnode.ETRNodeTypeCellData, // 目标节点
				int32(idx),
				cb,
				nil,
			)
		}
		var listenMode int32 = evhub.ListenModeTcp
		if cfg.ListenMode == "unix" {
			listenMode = evhub.ListenModeUnix
		}
		userData := &iframe.SessionUserData{
			DataType:       iframe.ESessionDataTypeNetInfo,
			NodeType:       trnode.ETRNodeTypeCellData, // 当前连接的节点
			NodeIndex:      int32(idx),
			DesInfo:        fmt.Sprintf("%dETRNodeTypeCellData%d", frameNode.tframeObj.GetFrameConfig().ZoneID, idx),
			IsServerClient: true,
		}
		failCount := 0
		for {
			if !evHub.Connect2(listenMode, cfg.ListenAddr, true, userData, connDo) {
				failCount++
				loghlp.Warnf("connect to dataserver fail")
			} else {
				break
			}
			if failCount >= 10 {
				loghlp.Errorf("connect dataserver fail,exit")
				return false
			}
			time.Sleep(time.Second * 1)
		}
	}
	return true
}
