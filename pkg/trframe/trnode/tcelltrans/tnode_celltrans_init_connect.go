/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-10-14 13:43:58
 * @LastEditTime: 2022-10-14 13:45:35
 * @FilePath: \trcell\pkg\trframe\trnode\tcelltrans\tnode_celltrans_init_connect.go
 */
package tcelltrans

func (frameNode *FrameNodeCellTrans) InitConnectServer() bool {
	// // 连接rootserver
	// frameConfig := frameNode.tframeObj.GetFrameConfig()
	// evHub := frameNode.tframeObj.GetEvHub()
	// for idx, cfg := range frameConfig.CellRootCfgs {
	// 	connDo := func() {
	// 		// 发送注册消息
	// 		reqMsg := &pbframe.FrameMsgRegisterServerInfoReq{
	// 			ZoneID:    frameNode.tframeObj.GetFrameConfig().ZoneID,
	// 			NodeType:  trnode.ETRNodeTypeCellCenter,
	// 			NodeIndex: frameNode.nodeIndex,
	// 			NodeDes:   fmt.Sprintf("%d_ETRNodeTypeCellCenter%d", frameNode.tframeObj.GetFrameConfig().ZoneID, frameNode.nodeIndex),
	// 		}
	// 		cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
	// 			loghlp.Infof("center register to cellroot callback suss:%d", okCode)
	// 		}
	// 		frameNode.tframeObj.ForwardMessage(
	// 			protocol.EMsgClassFrame,
	// 			protocol.EFrameMsgRegisterServerInfo,
	// 			reqMsg,
	// 			trnode.ETRNodeTypeCellRoot,
	// 			int32(idx),
	// 			cb,
	// 			nil,
	// 		)
	// 	}
	// 	userData := &iframe.SessionUserData{
	// 		DataType:       iframe.ESessionDataTypeNetInfo,
	// 		NodeType:       trnode.ETRNodeTypeCellRoot,
	// 		NodeIndex:      frameNode.nodeIndex,
	// 		DesInfo:        fmt.Sprintf("%dETRNodeTypeCellRoot%d", frameNode.tframeObj.GetFrameConfig().ZoneID, idx),
	// 		IsServerClient: true,
	// 	}
	// 	failCount := 0
	// 	var listenMode int32 = evhub.ListenModeTcp
	// 	if cfg.ListenMode == "unix" {
	// 		listenMode = evhub.ListenModeUnix
	// 	}
	// 	for {
	// 		if !evHub.Connect2(listenMode, cfg.ListenAddr, true, userData, connDo) {
	// 			failCount++
	// 			loghlp.Warnf("connect to rootserver fail")
	// 		} else {
	// 			break
	// 		}
	// 		if failCount >= 10 {
	// 			loghlp.Errorf("connect rootserver fail,exit")
	// 			return false
	// 		}
	// 		time.Sleep(time.Second * 1)
	// 	}
	// }
	// time.Sleep(time.Second)
	return true
}
