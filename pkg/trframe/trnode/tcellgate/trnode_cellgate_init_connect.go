/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-16 18:11:36
 * @LastEditTime: 2022-09-27 18:30:21
 * @FilePath: \trcell\pkg\trframe\trnode\tcellgate\trnode_cellgate_init_connect.go
 */
package tcellgate

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

func (frameNode *FrameNodeCellGate) InitConnectServer() bool {
	// 连接rootserver
	frameConfig := frameNode.tframeObj.GetFrameConfig()
	evHub := frameNode.tframeObj.GetEvHub()
	for idx, cfg := range frameConfig.CellRootCfgs {
		connDo := func() {
			// 发送注册消息
			reqMsg := &pbframe.FrameMsgRegisterServerInfoReq{
				ZoneID:    frameNode.tframeObj.GetFrameConfig().ZoneID,
				NodeType:  trnode.ETRNodeTypeCellGate,
				NodeIndex: frameNode.nodeIndex,
				NodeDes:   fmt.Sprintf("%d_ETRNodeTypeCellGate%d", frameNode.tframeObj.GetFrameConfig().ZoneID, frameNode.nodeIndex),
			}
			cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
				loghlp.Infof("gate register to rootserver callback suss:%d", okCode)
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
		var listenMode int32 = evhub.ListenModeTcp
		if cfg.ListenMode == "unix" {
			listenMode = evhub.ListenModeUnix
		}

		userData := &iframe.SessionUserData{
			DataType:       iframe.ESessionDataTypeNetInfo,
			NodeType:       trnode.ETRNodeTypeCellRoot,
			NodeIndex:      frameNode.nodeIndex,
			DesInfo:        fmt.Sprintf("%dETRNodeTypeCellRoot%d", frameNode.tframeObj.GetFrameConfig().ZoneID, idx),
			IsServerClient: true,
		}
		failCount := 0
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
	// 连接centerserver
	for idx, cfg := range frameConfig.CellCenterCfgs {
		connDo := func() {
			// 发送注册消息
			reqMsg := &pbframe.FrameMsgRegisterServerInfoReq{
				ZoneID:    frameNode.tframeObj.GetFrameConfig().ZoneID,
				NodeType:  trnode.ETRNodeTypeCellGate,
				NodeIndex: frameNode.nodeIndex,
				NodeDes:   fmt.Sprintf("%d_ETRNodeTypeCellGate%d", frameNode.tframeObj.GetFrameConfig().ZoneID, frameNode.nodeIndex),
			}
			cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
				loghlp.Infof("gate register to centerserver callback suss:%d", okCode)
			}
			frameNode.tframeObj.ForwardMessage(
				protocol.EMsgClassFrame,
				protocol.EFrameMsgRegisterServerInfo,
				reqMsg,
				trnode.ETRNodeTypeCellCenter,
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
			NodeType:       trnode.ETRNodeTypeCellCenter,
			NodeIndex:      frameNode.nodeIndex,
			DesInfo:        fmt.Sprintf("%dETRNodeTypeCellCenter%d", frameNode.tframeObj.GetFrameConfig().ZoneID, idx),
			IsServerClient: true,
		}
		failCount := 0
		for {
			if !evHub.Connect2(listenMode, cfg.ListenAddr, true, userData, connDo) {
				failCount++
				loghlp.Warnf("connect to centerserver fail")
			} else {
				break
			}
			if failCount >= 10 {
				loghlp.Errorf("connect centerserver fail,exit")
				return false
			}
			time.Sleep(time.Second * 1)
		}
	}
	// 连接logicserver
	for idx, cfg := range frameConfig.CellLogicCfgs {
		connDo := func() {
			// 发送注册消息
			reqMsg := &pbframe.FrameMsgRegisterServerInfoReq{
				ZoneID:    frameNode.tframeObj.GetFrameConfig().ZoneID,
				NodeType:  trnode.ETRNodeTypeCellGate,
				NodeIndex: frameNode.nodeIndex,
				NodeDes:   fmt.Sprintf("%d_ETRNodeTypeCellGate%d", frameNode.tframeObj.GetFrameConfig().ZoneID, frameNode.nodeIndex),
			}
			cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
				loghlp.Infof("gate register to centerserver callback suss:%d", okCode)
			}
			frameNode.tframeObj.ForwardMessage(
				protocol.EMsgClassFrame,
				protocol.EFrameMsgRegisterServerInfo,
				reqMsg,
				trnode.ETRNodeTypeCellLogic,
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
			NodeType:       trnode.ETRNodeTypeCellLogic,
			NodeIndex:      frameNode.nodeIndex,
			DesInfo:        fmt.Sprintf("%dETRNodeTypeCellLogic%d", frameNode.tframeObj.GetFrameConfig().ZoneID, idx),
			IsServerClient: true,
		}
		failCount := 0
		for {
			if !evHub.Connect2(listenMode, cfg.ListenAddr, true, userData, connDo) {
				failCount++
				loghlp.Warnf("connect to logicserver fail")
			} else {
				break
			}
			if failCount >= 10 {
				loghlp.Errorf("connect to logicserver fail,exit")
				return false
			}
			time.Sleep(time.Second * 1)
		}
	}
	// 连接viewserver
	for idx, cfg := range frameConfig.CellViewCfgs {
		connDo := func() {
			// 发送注册消息
			reqMsg := &pbframe.FrameMsgRegisterServerInfoReq{
				ZoneID:    frameNode.tframeObj.GetFrameConfig().ZoneID,
				NodeType:  trnode.ETRNodeTypeCellGate, // 当前自身节点
				NodeIndex: frameNode.nodeIndex,
				NodeDes:   fmt.Sprintf("%d_ETRNodeTypeCellGate%d", frameNode.tframeObj.GetFrameConfig().ZoneID, frameNode.nodeIndex),
			}
			cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
				loghlp.Infof("gate register to cellview callback suss:%d", okCode)
			}
			frameNode.tframeObj.ForwardMessage(
				protocol.EMsgClassFrame,
				protocol.EFrameMsgRegisterServerInfo,
				reqMsg,
				trnode.ETRNodeTypeCellView, // 目标节点
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
			NodeType:       trnode.ETRNodeTypeCellView, // 当前连接的节点
			NodeIndex:      frameNode.nodeIndex,
			DesInfo:        fmt.Sprintf("%dETRNodeTypeCellView%d", frameNode.tframeObj.GetFrameConfig().ZoneID, idx),
			IsServerClient: true,
		}
		failCount := 0
		for {
			if !evHub.Connect2(listenMode, cfg.ListenAddr, true, userData, connDo) {
				failCount++
				loghlp.Warnf("connect to viewserver fail")
			} else {
				break
			}
			if failCount >= 10 {
				loghlp.Errorf("connect viewserver fail,exit")
				return false
			}
			time.Sleep(time.Second * 1)
		}
	}
	return true
}
