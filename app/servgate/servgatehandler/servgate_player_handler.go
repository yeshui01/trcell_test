/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-20 14:24:39
 * @LastEditTime: 2022-09-28 09:55:04
 * @FilePath: \trcell\app\servgate\servgatehandler\servgate_player_handler.go
 */
package servgatehandler

import (
	"trcell/app/servgate/servgatemain"
	"trcell/pkg/crossdef"
	"trcell/pkg/evhub"
	"trcell/pkg/loghlp"
	"trcell/pkg/pb/pbclient"
	"trcell/pkg/pb/pbserver"
	"trcell/pkg/pb/pbtools"
	"trcell/pkg/protocol"
	"trcell/pkg/trframe"
	"trcell/pkg/trframe/iframe"
	"trcell/pkg/trframe/trnode"

	"github.com/sirupsen/logrus"
)

// 心跳
func HandlePlayerHeart(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbclient.ECMsgPlayerKeepHeartReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, pbtools.MakeErrorParams("pberror"), iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	rep := &pbclient.ECMsgPlayerKeepHeartRsp{}

	// 发送到logic更新heartTime
	// 这里的session是hateconnection
	hgsession := tmsgCtx.Session.(*servgatemain.HGateConnction)
	if hgsession.UserID > 0 {
		gateUser := servGate.GetUserManager().GetGateUser(hgsession.UserID)
		if gateUser != nil {
			loghlp.Debugf("recv gateuser(%d) heart", hgsession.UserID)
			logicServIndex := gateUser.GetNetPeerIndex(trnode.ETRNodeTypeCellLogic)
			trframe.PushZoneClientPBMessage(protocol.ECMsgClassPlayer,
				protocol.ECMsgPlayerKeepHeart,
				req,
				trnode.ETRNodeTypeCellLogic,
				logicServIndex,
				hgsession.UserID,
			)
			if gateUser.GetNetPeer(trnode.ETRNodeTypeCellView) != nil {
				trframe.PushZoneClientPBMessage(protocol.ECMsgClassPlayer,
					protocol.ECMsgPlayerKeepHeart,
					req,
					trnode.ETRNodeTypeCellView,
					gateUser.GetNetPeerIndex(trnode.ETRNodeTypeCellView),
					hgsession.UserID,
				)
			}
		}
	}

	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}

// 踢人
func HandlePlayerKickout(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbserver.ESMsgPlayerKickOutReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	loghlp.Warnf("HandlePlayerKickout:%+v", req)
	rep := &pbserver.ESMsgPlayerKickOutRep{}
	gateUser := servGate.GetUserManager().GetGateUser(req.RoleID)
	if gateUser != nil {
		loghlp.Warnf("kickout player(%d), reason:%d succ!!!", gateUser.GetGateConnect().UserID, req.Reason)
		gateConnect := gateUser.GetGateConnect()
		emptyMsg := evhub.MakeEmptyMessage()
		emptyMsg.Head.HasSecond = 1
		emptyMsg.SecondHead = &evhub.NetMsgSecondHead{
			ID: req.RoleID,
		}
		gateConnect.SendMsg(emptyMsg)

		gateConnect.UserID = 0 // 解除关联
		// 删除用户
		servGate.GetUserManager().DelGateUser(req.RoleID)
		servGate.GetGateConnMgr().RemoveConnection(gateConnect.TcpPeerConn)
	} else {
		loghlp.Warnf("kickout player(%d), reason:%d, but not find gate user", req.RoleID, req.Reason)
	}

	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}

// 创角
func HandleECMsgPlayerCreateRole(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbclient.ECMsgPlayerCreateRoleReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, pbtools.MakeErrorParams("pb error"), iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	//rep := &pbclient.ECMsgPlayerCreateRoleRsp{}
	// 验证touken
	ok, tokenRes := crossdef.TokenAuthClaims(req.Token, crossdef.SignKey)
	if !ok {
		logrus.Error("token parse fail")
		return protocol.ECodeTokenExpire, pbtools.MakeErrorParams("tokenRes error"), iframe.EHandleContent
	} else {
		loghlp.Infof("parse player token success:%+v", *tokenRes)
	}

	// ->center
	csReq := &pbserver.ESMsgPlayerCreateRoleReq{
		UserID:   tokenRes.UserID,
		UserName: tokenRes.Account,
		Nickname: req.Nickname,
		GateInfo: &pbserver.ServerNodeInfo{
			ZoneID:    trframe.GetFrameConfig().ZoneID,
			NodeType:  trnode.ETRNodeTypeCellGate,
			NodeIndex: trframe.GetCurNodeIndex(),
		},
	}
	hgsession := tmsgCtx.Session.(*servgatemain.HGateConnction)
	// 发送消息
	cbDo := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
		loghlp.Infof("create role callback success,okCode:%d", okCode)
		if okCode == protocol.ECodeSuccess {
			cbRep := &pbserver.ESMsgPlayerCreateRoleRep{}
			if !trframe.DecodePBMessage2(msgData, cbRep) {
				loghlp.Error("decode cbRep error")
				// trframe.SendReplyMessage(protocol.ECodePBDecodeError, nil, env)
				servGate.SendTcpClientReplyMessage(okCode, nil, env)
				return
			}
			cltRep := &pbclient.ECMsgPlayerCreateRoleRsp{
				RoleData: cbRep.RoleData,
			}
			// 关联玩家数据
			hgsession.UserID = cbRep.RoleData.RoleID
			loghlp.Infof("user_create_role_succ, roleid:%d, ipAddr:%s",
				hgsession.UserID,
				hgsession.TcpPeerConn.RemoteAddr(),
			)
			// gate user数据初始化
			gateUser := servgatemain.NewGateUser(cbRep.RoleData.RoleID)
			gateUser.SetNetPeer(trnode.ETRNodeTypeCellLogic, &trnode.TRNodeInfo{
				ZoneID:    trframe.GetFrameConfig().ZoneID,
				NodeType:  trnode.ETRNodeTypeCellLogic,
				NodeIndex: cbRep.LogicInfo.NodeIndex,
			})
			gateUser.SetGateConnect(hgsession)
			servGate.GetUserManager().AddGateUser(cbRep.RoleData.RoleID, gateUser)
			servGate.SendTcpClientReplyMessage(okCode, cltRep, env)
			loghlp.Info("gate player(%d) create succ,logicIndex(%d),roleInfo:%+v",
				cbRep.RoleData.RoleID,
				cbRep.LogicInfo.NodeIndex,
				cbRep.RoleData,
			)
		} else {
			loghlp.Errorf("create role fail!!!")
			servGate.SendTcpClientReplyMessage(okCode, nil, env)
		}
	}

	cbEnv := trframe.MakeMsgEnv(0,
		tmsgCtx.NetMessage)
	cbEnv.UserData = hgsession

	trframe.ForwardZoneMessage(
		protocol.ESMsgClassPlayer,
		protocol.ESMsgPlayerCreateRole,
		csReq,
		trnode.ETRNodeTypeCellCenter,
		0,
		cbDo,
		cbEnv,
	)
	return protocol.ECodeSuccess, nil, iframe.EHandlePending
}

// 登录游戏
func HandleECMsgPlayerLoginGame(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbclient.ECMsgPlayerLoginGameReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, pbtools.MakeErrorParams("pb error"), iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	// 验证touken
	ok, tokenRes := crossdef.TokenAuthClaims(req.Token, crossdef.SignKey)
	if !ok {
		logrus.Error("token parse fail")
		return protocol.ECodeTokenExpire, pbtools.MakeErrorParams("tokenRes error"), iframe.EHandleContent
	} else {
		loghlp.Infof("parse player token success:%+v", *tokenRes)
	}

	// ->center
	csReq := &pbserver.ESMsgPlayerLoginGameReq{
		UserID:  tokenRes.UserID,
		Account: tokenRes.Account,
		RoleID:  req.RoleID,
		GateInfo: &pbserver.ServerNodeInfo{
			ZoneID:    trframe.GetFrameConfig().ZoneID,
			NodeType:  trnode.ETRNodeTypeCellGate,
			NodeIndex: trframe.GetCurNodeIndex(),
		},
	}
	hgsession := tmsgCtx.Session.(*servgatemain.HGateConnction)
	// 发送消息
	cbDo := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
		loghlp.Infof("user(%d)[%s] logingame callback success,okCode:%d", okCode, tokenRes.UserID, tokenRes.Account)
		if okCode == protocol.ECodeSuccess {
			cbRep := &pbserver.ESMsgPlayerLoginGameRep{}
			if !trframe.DecodePBMessage2(msgData, cbRep) {
				loghlp.Error("decode cbRep error")
				// trframe.SendReplyMessage(protocol.ECodePBDecodeError, nil, env)
				servGate.SendTcpClientReplyMessage(protocol.ECodePBDecodeError,
					pbtools.MakeErrorParams("logingame error"), env)
				return
			}

			preGateUser := servGate.GetUserManager().GetGateUser(cbRep.RoleData.RoleID)
			if preGateUser != nil {
				// 踢掉这个连接
				preGateUser.GetGateConnect().UserID = 0
				loghlp.Warnf("tmp kick pregate user(%d)", cbRep.RoleData.RoleID)
				// preGateUser.GetGateConnect().TcpPeerConn.Close()
				// 发送一个主动关闭的消息
				closeMsg := evhub.MakeEmptyMessage()
				closeMsg.Head.HasSecond = 1
				closeMsg.SecondHead = &evhub.NetMsgSecondHead{
					ID: req.RoleID,
				}
				preGateUser.GetGateConnect().SendMsg(closeMsg)
				servGate.GetGateConnMgr().RemoveConnection(preGateUser.GetGateConnect().TcpPeerConn)
				servGate.GetUserManager().DelGateUser(cbRep.RoleData.RoleID)
				//preGateUser.GetGateConnect().Stop()
			}
			cltRep := &pbclient.ECMsgPlayerLoginGameRsp{
				RoleData: cbRep.RoleData,
			}
			// 关联玩家数据
			hgsession.UserID = cbRep.RoleData.RoleID
			loghlp.Infof("user_login_game_succ, roleid:%d, ipAddr:%s",
				hgsession.UserID,
				hgsession.TcpPeerConn.RemoteAddr(),
			)
			// gate user数据初始化
			gateUser := servgatemain.NewGateUser(cbRep.RoleData.RoleID)
			gateUser.SetNetPeer(trnode.ETRNodeTypeCellLogic, &trnode.TRNodeInfo{
				ZoneID:    trframe.GetFrameConfig().ZoneID,
				NodeType:  trnode.ETRNodeTypeCellLogic,
				NodeIndex: cbRep.LogicInfo.NodeIndex,
			})
			gateUser.SetGateConnect(hgsession)

			servGate.GetUserManager().AddGateUser(cbRep.RoleData.RoleID, gateUser)
			servGate.SendTcpClientReplyMessage(okCode, cltRep, env)
			loghlp.Info("gate player(%d) create succ,logicIndex(%d),roleInfo:%+v",
				cbRep.RoleData.RoleID,
				cbRep.LogicInfo.NodeIndex,
				cbRep.RoleData,
			)
		} else {
			loghlp.Errorf("user(%d)[%s]login game fail(%d)!!!",
				tokenRes.UserID,
				tokenRes.Account,
				okCode)
			servGate.SendTcpClientReplyMessage2(okCode, msgData, env)
		}
	}

	cbEnv := trframe.MakeMsgEnv(0,
		tmsgCtx.NetMessage)
	cbEnv.UserData = hgsession

	trframe.ForwardZoneMessage(
		protocol.ESMsgClassPlayer,
		protocol.ESMsgPlayerLoginGame,
		csReq,
		trnode.ETRNodeTypeCellCenter,
		0,
		cbDo,
		cbEnv,
	)
	return protocol.ECodeSuccess, nil, iframe.EHandlePending
}
