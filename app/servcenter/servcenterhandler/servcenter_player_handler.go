/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-22 16:37:27
 * @LastEditTime: 2022-09-28 09:59:44
 * @FilePath: \trcell\app\servcenter\servcenterhandler\servcenter_player_handler.go
 */
package servcenterhandler

import (
	"trcell/pkg/loghlp"
	"trcell/pkg/pb/pbserver"
	"trcell/pkg/pb/pbtools"
	"trcell/pkg/protocol"
	"trcell/pkg/trframe"
	"trcell/pkg/trframe/iframe"
	"trcell/pkg/trframe/trnode"
)

// 创角-center
func HandleESMsgPlayerCreateRole(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbserver.ESMsgPlayerCreateRoleReq{}
	rep := &pbserver.ESMsgPlayerCreateRoleRep{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, pbtools.MakeErrorParams("pb error"), iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	csGlobal := servCenter.GetCenterGlobal()
	if csGlobal.IsUserCreateLock(req.UserID) {
		return protocol.ECodeRoleCreatingLock,
			pbtools.MakeErrorParams("user create locking"),
			iframe.EHandleContent
	}
	if csGlobal.FindCsPlayerByName(req.Nickname) != nil {
		return protocol.ECodeRoleNickNameExisted,
			pbtools.MakeErrorParams("ECodeRoleNickNameExisted"),
			iframe.EHandleContent
	}

	if csGlobal.IsNicknameLock(req.Nickname) {
		return protocol.ECodeRoleNickNameExisted,
			pbtools.MakeErrorParams("nickname locked"),
			iframe.EHandleContent
	}
	logicList := trframe.GetNodeListByType(trnode.ETRNodeTypeCellLogic)
	if len(logicList) == 0 {
		return protocol.ECodeLogicServException,
			pbtools.MakeErrorParams("logicserv error"),
			iframe.EHandleContent
	}
	// 锁定
	csGlobal.LockUserCreateRole(req.UserID)
	csGlobal.LockNickname(req.Nickname)

	// logicIndex := rand.Intn(len(logicList))
	logicIndex := int(req.UserID) % len(logicList)
	cbDo := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
		loghlp.Infof("center create role callback okCode:%d", okCode)
		if okCode == protocol.ECodeSuccess {
			cbRep := &pbserver.ESMsgPlayerCreateRoleRep{}
			if !trframe.DecodePBMessage2(msgData, cbRep) {
				loghlp.Error("decode cbRep error")
				trframe.SendReplyMessage(protocol.ECodePBDecodeError,
					pbtools.MakeErrorParams("ESMsgPlayerCreateRoleRep pberror"),
					env)
				//ser.SendTcpClientReplyMessage(okCode, nil, env)
				return
			}
			trframe.LogCbMsgInfo(env.SrcMessage, cbRep)
			// 创建csplayer
			centerPlayer := csGlobal.CreateCsPlayer(cbRep.RoleData.RoleID,
				cbRep.RoleData.Nickname,
				cbRep.RoleData.Level,
				0)
			if centerPlayer == nil {
				trframe.SendReplyMessage(protocol.ECodePBDecodeError,
					pbtools.MakeErrorParams("centerPlayer create nil"),
					env)
				return
			}
			// 更新 player node info
			centerPlayer.SetNetPeer(trnode.ETRNodeTypeCellGate, trnode.NewNodeInfo(req.GateInfo.ZoneID, req.GateInfo.NodeType, req.GateInfo.NodeIndex))
			centerPlayer.SetNetPeer(trnode.ETRNodeTypeCellLogic, trnode.NewNodeInfo(cbRep.LogicInfo.ZoneID, cbRep.LogicInfo.NodeType, cbRep.LogicInfo.NodeIndex))
			csGlobal.UnlockNickname(cbRep.RoleData.Nickname)
			csGlobal.UnlockUserCreate(req.UserID)
			trframe.SendReplyMessage(okCode,
				cbRep,
				env)
		} else {
			loghlp.Errorf("create role fail!!!")
			trframe.SendReplyErrorMessage(
				okCode,
				msgData,
				env,
			)
		}
	}

	trframe.ForwardZoneMessage(protocol.ESMsgClassPlayer,
		protocol.ESMsgPlayerCreateRole,
		req,
		trnode.ETRNodeTypeCellLogic,
		int32(logicIndex),
		cbDo,
		trframe.MakeMsgEnv2(tmsgCtx.Session, tmsgCtx.NetMessage),
	)
	return protocol.ECodeSuccess, rep, iframe.EHandlePending
}

// 登录游戏-center
func HandleESMsgPlayerLoginGame(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbserver.ESMsgPlayerLoginGameReq{}
	rep := &pbserver.ESMsgPlayerLoginGameRep{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, pbtools.MakeErrorParams("pb error"), iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	csGlobal := servCenter.GetCenterGlobal()
	logicList := trframe.GetNodeListByType(trnode.ETRNodeTypeCellLogic)
	if len(logicList) == 0 {
		return protocol.ECodeLogicServException,
			pbtools.MakeErrorParams("logicserv error"),
			iframe.EHandleContent
	}
	// 锁定登录状态TODO

	// logicIndex := rand.Intn(len(logicList))
	logicIndex := int(req.UserID) % len(logicList)
	cbDo := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
		loghlp.Infof("logingame callback okCode:%d", okCode)
		if okCode == protocol.ECodeSuccess {
			cbRep := &pbserver.ESMsgPlayerLoginGameRep{}
			if !trframe.DecodePBMessage2(msgData, cbRep) {
				loghlp.Error("decode cbRep error")
				trframe.SendReplyMessage(protocol.ECodePBDecodeError,
					pbtools.MakeErrorParams("ESMsgPlayerLoginGameRep pberror"),
					env)
				//ser.SendTcpClientReplyMessage(okCode, nil, env)
				return
			}
			trframe.LogCbMsgInfo(env.SrcMessage, cbRep)
			centerPlayer := csGlobal.FindCsPlayer(cbRep.RoleData.RoleID)
			if centerPlayer == nil {
				centerPlayer = csGlobal.CreateCsPlayer(cbRep.RoleData.RoleID,
					cbRep.RoleData.Nickname,
					cbRep.RoleData.Level,
					0)
			}
			if centerPlayer == nil {
				trframe.SendReplyMessage(protocol.ECodePBDecodeError,
					pbtools.MakeErrorParams("centerPlayer init nil"),
					env)
				return
			}
			// 更新 player node info
			centerPlayer.SetNetPeer(trnode.ETRNodeTypeCellGate,
				trnode.NewNodeInfo(req.GateInfo.ZoneID,
					req.GateInfo.NodeType,
					req.GateInfo.NodeIndex),
			)
			centerPlayer.SetNetPeer(trnode.ETRNodeTypeCellLogic,
				trnode.NewNodeInfo(cbRep.LogicInfo.ZoneID,
					cbRep.LogicInfo.NodeType,
					cbRep.LogicInfo.NodeIndex),
			)
			csGlobal.UnlockNickname(cbRep.RoleData.Nickname)
			csGlobal.UnlockUserCreate(req.UserID)
			trframe.SendReplyMessage(okCode,
				cbRep,
				env)
		} else {
			loghlp.Errorf("login role fail!!!")
			trframe.SendReplyErrorMessage(
				okCode,
				msgData,
				env,
			)
		}
	}

	trframe.ForwardZoneMessage(protocol.ESMsgClassPlayer,
		protocol.ESMsgPlayerLoginGame,
		req,
		trnode.ETRNodeTypeCellLogic,
		int32(logicIndex),
		cbDo,
		trframe.MakeMsgEnv2(tmsgCtx.Session, tmsgCtx.NetMessage),
	)
	return protocol.ECodeSuccess, rep, iframe.EHandlePending
}

// 离线
func HandleESMsgPlayerDisconnect(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbserver.ESMsgPlayerDisconnectReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, pbtools.MakeErrorParams("pberror"), iframe.EHandleContent
	}
	rep := &pbserver.ESMsgPlayerDisconnectRep{}
	loghlp.Infof("player(%d) disconnect", tmsgCtx.NetMessage.SecondHead.ID)
	// logicPlayer := servCenter.GetCenterGlobal().FindPlayer(tmsgCtx.NetMessage.SecondHead.ID)
	// if logicPlayer != nil {
	// 	servLogic.GetLogicGlobal().HandlePlayerOffline(logicPlayer)
	// }
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}
