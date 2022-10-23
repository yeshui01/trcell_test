/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-16 18:18:35
 * @LastEditTime: 2022-09-27 18:27:20
 * @FilePath: \trcell\pkg\trframe\trnode\tcelldata\trnode_celldata_init_connnect.go
 */
package tcelldata

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

func (frameNode *FrameNodeCellData) InitConnectServer() bool {
	// 连接rootserver
	frameConfig := frameNode.tframeObj.GetFrameConfig()
	evHub := frameNode.tframeObj.GetEvHub()

	for idx, cfg := range frameConfig.CellRootCfgs {
		connDo := func() {
			// 发送注册消息
			reqMsg := &pbframe.FrameMsgRegisterServerInfoReq{
				ZoneID:    frameNode.tframeObj.GetFrameConfig().ZoneID,
				NodeType:  trnode.ETRNodeTypeCellData,
				NodeIndex: frameNode.nodeIndex,
				NodeDes:   fmt.Sprintf("%d_ETRNodeTypeCellData%d", frameNode.tframeObj.GetFrameConfig().ZoneID, frameNode.nodeIndex),
			}
			cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
				loghlp.Infof("data register to cellroot callback suss:%d", okCode)
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
	return true
}
