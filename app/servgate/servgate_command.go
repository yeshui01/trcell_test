/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 15:54:59
 * @LastEditTime: 2022-09-28 09:52:58
 * @FilePath: \trcell\app\servgate\servgate_command.go
 */
package servgate

import (
	"net"
	"trcell/pkg/loghlp"
	"trcell/pkg/pb/pbcmd"
	"trcell/pkg/pb/pbserver"
	"trcell/pkg/protocol"
	"trcell/pkg/sconst"
	"trcell/pkg/trframe"
	"trcell/pkg/trframe/iframe"
	"trcell/pkg/trframe/trnode"
)

const (
	MsgFactor = 1000
)

func (serv *CellServGate) HandleCommand(frameCmd *trframe.TRFrameCommand) {
	hKey := GetCmdHandlerKey(frameCmd.UserCmd.GetCmdClass(), frameCmd.UserCmd.GetCmdType())
	if h, ok := serv.CmdHandlerMap[hKey]; ok {
		h(frameCmd)
	} else {
		loghlp.Errorf("not find cmd handler(%d_%d)", frameCmd.UserCmd.GetCmdClass(), frameCmd.UserCmd.GetCmdType())
	}
}
func (serv *CellServGate) HandleCmdTcpConnect(frameCmd *trframe.TRFrameCommand) {
	tcpConn, ok := frameCmd.UserCmd.GetCmdData().(net.Conn)
	if !ok {
		loghlp.Errorf("HandleCmdTcpConnect data error")
		return
	}
	loghlp.Infof("HandleCmdTcpConnect newTcpConn,remoteAddr:%s", tcpConn.RemoteAddr().String())
	serv.ConnMgr.AddConnection(tcpConn)
}

func (serv *CellServGate) HandleCmdTcpClose(frameCmd *trframe.TRFrameCommand) {
	tcpConn, ok := frameCmd.UserCmd.GetCmdData().(net.Conn)
	if !ok {
		loghlp.Errorf("HandleCmdTcpClose data error")
	}
	servClient := serv.ConnMgr.GetConnection(tcpConn)
	if servClient != nil {
		servClient.Stop()
	}
	loghlp.Info("HandleCmdTcpClose")

	// 通知centerserver玩家断线
	if servClient != nil && servClient.UserID > 0 {
		userCsIndex := int32(0)
		gateUser := serv.UserMgr.GetGateUser(servClient.UserID)
		// if gateUser != nil {
		// 	hallNode := gateUser.GetCenterNode()
		// 	if hallNode != nil {
		// 		userHallIndex = hallNode.NodeIndex
		// 	}
		// }
		offlineReq := &pbserver.ESMsgPlayerDisconnectReq{
			RoleID: servClient.UserID,
			Reason: sconst.EPlayerOfflineReasonNormal,
		}
		cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
			loghlp.Infof("offline notify success,okCode:%d", okCode)
		}
		trframe.ForwardZoneMessage(
			protocol.ESMsgClassPlayer,
			protocol.ESMsgPlayerDisconnect,
			offlineReq,
			trnode.ETRNodeTypeCellCenter,
			userCsIndex,
			cb,
			nil,
		)
		if gateUser != nil {
			logicIdx := gateUser.GetNetPeerIndex(trnode.ETRNodeTypeCellLogic)
			trframe.PushZoneClientPBMessage(protocol.ESMsgClassPlayer,
				protocol.ESMsgPlayerDisconnect,
				offlineReq,
				trnode.ETRNodeTypeCellLogic,
				logicIdx,
				servClient.UserID,
			)
		}

		// 删除用户
		serv.UserMgr.DelGateUser(servClient.UserID)
	}

	serv.ConnMgr.RemoveConnection(tcpConn)
}

func (serv *CellServGate) HandleCmdTcpMessage(frameCmd *trframe.TRFrameCommand) {
	tcpMessage, ok := frameCmd.UserCmd.GetCmdData().(*pbcmd.CmdTypeTcpsocketMessageData)
	if !ok {
		loghlp.Errorf("HandleCmdTcpMessage data error")
	}
	loghlp.Infof("recv tcpClientMessage,totalLen(%d),dataLen(%d)", tcpMessage.HubMsg.Head.Len, len(tcpMessage.HubMsg.Data))
	// 处理用户消息
	hgc := serv.ConnMgr.GetConnection(tcpMessage.TcpConn)
	if hgc != nil {
		if serv.IsForwardToView(int32(tcpMessage.HubMsg.Head.MsgClass), int32(tcpMessage.HubMsg.Head.MsgType)) {
			cbEnv := trframe.MakeMsgEnv(0,
				tcpMessage.HubMsg)
			cbEnv.UserData = hgc
			cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
				loghlp.Infof("trans view client msg callback succ,okCode:%d", okCode)
				serv.SendTcpClientReplyMessage2(okCode, msgData, env)
			}
			gateUser := serv.UserMgr.GetGateUser(hgc.UserID)
			if gateUser != nil {
				if gateUser.GetNetPeer(trnode.ETRNodeTypeCellView) != nil {
					trframe.ForwardZoneClientMessage(
						int32(tcpMessage.HubMsg.Head.MsgClass),
						int32(tcpMessage.HubMsg.Head.MsgType),
						tcpMessage.HubMsg.Data,
						trnode.ETRNodeTypeCellView,
						gateUser.GetNetPeerIndex(trnode.ETRNodeTypeCellView),
						cb,
						cbEnv,
						hgc.UserID,
					)
				} else {
					loghlp.Errorf("not find player(%d) view node", hgc.UserID)
					trframe.ForwardZoneClientMessage(
						int32(tcpMessage.HubMsg.Head.MsgClass),
						int32(tcpMessage.HubMsg.Head.MsgType),
						tcpMessage.HubMsg.Data,
						trnode.ETRNodeTypeCellView,
						0,
						cb,
						cbEnv,
						hgc.UserID,
					)
				}
			} else {
				loghlp.Errorf("not find gate user:%d", hgc.UserID)
			}
		} else if serv.IsForwardToCenter(int32(tcpMessage.HubMsg.Head.MsgClass), int32(tcpMessage.HubMsg.Head.MsgType)) {
			cbEnv := trframe.MakeMsgEnv(0,
				tcpMessage.HubMsg)
			cbEnv.UserData = hgc
			cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
				loghlp.Infof("trans center client msg callback succ,okCode:%d", okCode)
				serv.SendTcpClientReplyMessage2(okCode, msgData, env)
			}
			gateUser := serv.UserMgr.GetGateUser(hgc.UserID)
			if gateUser != nil {
				trframe.ForwardZoneClientMessage(
					int32(tcpMessage.HubMsg.Head.MsgClass),
					int32(tcpMessage.HubMsg.Head.MsgType),
					tcpMessage.HubMsg.Data,
					trnode.ETRNodeTypeCellCenter,
					0,
					cb,
					cbEnv,
					hgc.UserID,
				)
			} else {
				loghlp.Errorf("not find gate user:%d", hgc.UserID)
			}
		} else if serv.IsForwardToGame(int32(tcpMessage.HubMsg.Head.MsgClass), int32(tcpMessage.HubMsg.Head.MsgType)) {
			cbEnv := trframe.MakeMsgEnv(0,
				tcpMessage.HubMsg)
			cbEnv.UserData = hgc
			cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
				loghlp.Infof("trans game client msg callback succ,okCode:%d", okCode)
				serv.SendTcpClientReplyMessage2(okCode, msgData, env)
			}
			gateUser := serv.UserMgr.GetGateUser(hgc.UserID)
			if gateUser != nil {
				trframe.ForwardZoneClientMessage(
					int32(tcpMessage.HubMsg.Head.MsgClass),
					int32(tcpMessage.HubMsg.Head.MsgType),
					tcpMessage.HubMsg.Data,
					trnode.ETRNodeTypeCellGame,
					0,
					cb,
					cbEnv,
					hgc.UserID,
				)
			} else {
				loghlp.Errorf("not find gate user:%d", hgc.UserID)
			}
		} else if !trframe.DispatchMsg(hgc, tcpMessage.HubMsg, nil) {
			// 默认处理,网关不处理的,转发到logic处理
			cbEnv := trframe.MakeMsgEnv(0,
				tcpMessage.HubMsg)
			cbEnv.UserData = hgc
			cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
				loghlp.Infof("trans logic client msg callback succ,okCode:%d", okCode)
				serv.SendTcpClientReplyMessage2(okCode, msgData, env)
			}
			gateUser := serv.UserMgr.GetGateUser(hgc.UserID)
			if gateUser != nil {
				if gateUser.GetNetPeer(trnode.ETRNodeTypeCellLogic) != nil {
					trframe.ForwardZoneClientMessage(
						int32(tcpMessage.HubMsg.Head.MsgClass),
						int32(tcpMessage.HubMsg.Head.MsgType),
						tcpMessage.HubMsg.Data,
						trnode.ETRNodeTypeCellLogic,
						gateUser.GetNetPeerIndex(trnode.ETRNodeTypeCellLogic),
						cb,
						cbEnv,
						hgc.UserID,
					)
				} else {
					trframe.ForwardZoneClientMessage(
						int32(tcpMessage.HubMsg.Head.MsgClass),
						int32(tcpMessage.HubMsg.Head.MsgType),
						tcpMessage.HubMsg.Data,
						trnode.ETRNodeTypeCellLogic,
						0,
						cb,
						cbEnv,
						hgc.UserID,
					)
				}
			} else {
				trframe.ForwardZoneClientMessage(
					int32(tcpMessage.HubMsg.Head.MsgClass),
					int32(tcpMessage.HubMsg.Head.MsgType),
					tcpMessage.HubMsg.Data,
					trnode.ETRNodeTypeCellLogic,
					0,
					cb,
					cbEnv,
					hgc.UserID,
				)
			}
		}
	} else {
		loghlp.Error("not find hgc object!!!!")
	}
}
