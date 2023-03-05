/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-22 17:04:57
 * @LastEditTime: 2022-09-28 09:48:45
 * @FilePath: \trcell\app\servlogic\servlogichandler\servlogic_player_handler.go
 */
package servlogichandler

import (
	"trcell/app/servlogic/servlogicmain"
	"trcell/pkg/loghlp"
	"trcell/pkg/pb/pbclient"
	"trcell/pkg/pb/pbserver"
	"trcell/pkg/pb/pbtools"
	"trcell/pkg/protocol"
	"trcell/pkg/timeutil"
	"trcell/pkg/trframe"
	"trcell/pkg/trframe/iframe"
	"trcell/pkg/trframe/trnode"
)

// 创角-logic
func HandleESMsgPlayerCreateRole(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbserver.ESMsgPlayerCreateRoleReq{}
	rep := &pbserver.ESMsgPlayerCreateRoleRep{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	if len(req.Nickname) < 1 || req.UserID < 1 {
		return protocol.ECodeParamError,
			pbtools.MakeErrorParams("nickname error"),
			iframe.EHandleContent
	}

	// db创角
	dbList := trframe.GetNodeListByType(trnode.ETRNodeTypeCellData)
	if len(dbList) == 0 {
		return protocol.ECodeDataServException,
			pbtools.MakeErrorParams("dbserv error"),
			iframe.EHandleContent
	}
	logicGlobal := servLogic.GetLogicGlobal()
	dbIndex := int(req.UserID) % len(dbList)
	cbDo := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
		loghlp.Infof("logic create role callback okCode:%d", okCode)
		if okCode == protocol.ECodeSuccess {
			cbRep := &pbserver.ESMsgPlayerCreateRoleRep{}
			if !trframe.DecodePBMessage2(msgData, cbRep) {
				loghlp.Error("decode cbRep error")
				trframe.SendReplyMessage(protocol.ECodePBDecodeError,
					pbtools.MakeErrorParams("ESMsgPlayerCreateRoleRep pberror"),
					env)
				return
			}
			trframe.LogCbMsgInfo(tmsgCtx.NetMessage, cbRep)
			// 创角成功
			logicPlayer := logicGlobal.FindPlayer(cbRep.RoleData.RoleID)
			if logicPlayer != nil {
				loghlp.Errorf("logicplayer(%d) has existed, maybe error!", cbRep.RoleData.RoleID)
			}
			logicPlayer = servlogicmain.NewLogicPlayer(cbRep.RoleData.RoleID)
			logicPlayer.LoadData(cbRep.RoleDetail)
			logicGlobal.AddPlayer(logicPlayer)
			// 上线处理
			logicPlayer.SetNetPeer(trnode.ETRNodeTypeCellGate, &trnode.TRNodeInfo{
				ZoneID:    req.GateInfo.ZoneID,
				NodeType:  req.GateInfo.NodeType,
				NodeIndex: req.GateInfo.NodeIndex,
			})
			logicPlayer.SetNetPeer(trnode.ETRNodeTypeCellData, &trnode.TRNodeInfo{
				ZoneID:    trframe.GetFrameConfig().ZoneID,
				NodeType:  trnode.ETRNodeTypeCellData,
				NodeIndex: int32(dbIndex),
			})
			// 这里后续处理放到下一帧处理
			logicPlayer.IsOnline = true
			logicPlayer.UpdateHeartTime(timeutil.NowTime())
			trframe.AfterMsgJob(func() {
				logicPlayer.UpdateHeartTime(timeutil.NowTime())
				logicGlobal.HandlePlayerOnline(logicPlayer)
			})
			cbRep.RoleData = logicPlayer.ToClientRoleInfo()
			cbRep.RoleDetail = nil // 这个数据可以不用返回了
			cbRep.DataInfo = &pbserver.ServerNodeInfo{
				ZoneID:    trframe.GetFrameConfig().ZoneID,
				NodeType:  trnode.ETRNodeTypeCellData,
				NodeIndex: int32(dbIndex),
			}
			cbRep.LogicInfo = trframe.GetCurPBNodeInfo()
			trframe.SendReplyMessage(okCode,
				cbRep,
				env)
			loghlp.Infof("user(%d)[%s][%s] create role succ finshi!",
				req.UserID,
				req.UserName,
				req.Nickname)
		} else {
			loghlp.Errorf("create role fail!!!")
			trframe.SendReplyErrorMessage(protocol.ECodePBDecodeError,
				msgData,
				env)
		}
	}

	trframe.ForwardZoneMessage(protocol.ESMsgClassPlayer,
		protocol.ESMsgPlayerCreateRole,
		req,
		trnode.ETRNodeTypeCellData,
		int32(dbIndex),
		cbDo,
		trframe.MakeMsgEnv2(tmsgCtx.Session, tmsgCtx.NetMessage),
	)
	return protocol.ECodeSuccess, rep, iframe.EHandlePending
}

// 登录游戏-logic
func HandleESMsgPlayerLoginGame(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbserver.ESMsgPlayerLoginGameReq{}
	rep := &pbserver.ESMsgPlayerLoginGameRep{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, pbtools.MakeErrorParams("pberror"), iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	if req.UserID < 1 {
		return protocol.ECodeParamError,
			pbtools.MakeErrorParams("userid param error"),
			iframe.EHandleContent
	}
	logicGlobal := servLogic.GetLogicGlobal()
	logicPlayer := logicGlobal.FindPlayerByUserID(req.UserID)
	if logicPlayer != nil {
		// 这里后续处理放到下一帧处理
		logicPlayer.IsOnline = true
		logicPlayer.UpdateHeartTime(timeutil.NowTime())
		trframe.AfterMsgJob(func() {
			logicPlayer.UpdateHeartTime(timeutil.NowTime())
			logicGlobal.HandlePlayerOnline(logicPlayer)
		})
		rep.RoleData = logicPlayer.ToClientRoleInfo()
		rep.DataInfo = &pbserver.ServerNodeInfo{
			ZoneID:    trframe.GetFrameConfig().ZoneID,
			NodeType:  trnode.ETRNodeTypeCellData,
			NodeIndex: int32(logicPlayer.GetNetPeerIndex(trnode.ETRNodeTypeCellData)),
		}
		rep.LogicInfo = trframe.GetCurPBNodeInfo()
		loghlp.Debugf("player(%d) login from logic_cache", logicPlayer.RoleID)
		return protocol.ECodeSuccess, rep, iframe.EHandleContent
	}
	// db加载数据
	dbList := trframe.GetNodeListByType(trnode.ETRNodeTypeCellData)
	if len(dbList) == 0 {
		return protocol.ECodeDataServException,
			pbtools.MakeErrorParams("dbserv error"),
			iframe.EHandleContent
	}

	dbIndex := int(req.UserID) % len(dbList)
	cbDo := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
		loghlp.Infof("account(%d) logic load callback okCode:%d", okCode, req.UserID)
		if okCode == protocol.ECodeSuccess {
			cbRep := &pbserver.ESMsgPlayerLoadRoleRep{}
			if !trframe.DecodePBMessage2(msgData, cbRep) {
				loghlp.Error("decode cbRep error")
				trframe.SendReplyMessage(protocol.ECodePBDecodeError,
					pbtools.MakeErrorParams("ESMsgPlayerCreateRoleRep pberror"),
					env)
				return
			}
			trframe.LogCbMsgInfo(tmsgCtx.NetMessage, cbRep)
			// 加载成功
			logicPlayer := logicGlobal.FindPlayer(cbRep.RoleID)
			if logicPlayer == nil {
				logicPlayer = servlogicmain.NewLogicPlayer(cbRep.RoleID)
				logicPlayer.LoadData(cbRep.RoleDetailData)
				logicGlobal.AddPlayer(logicPlayer)
			} else {
				loghlp.Errorf("logicplayer(%d) has existed, maybe error!", cbRep.RoleID)
				//logicPlayer.LoadData(cbRep.RoleDetailData)
			}
			// 上线处理
			logicPlayer.SetNetPeer(trnode.ETRNodeTypeCellGate, &trnode.TRNodeInfo{
				ZoneID:    req.GateInfo.ZoneID,
				NodeType:  req.GateInfo.NodeType,
				NodeIndex: req.GateInfo.NodeIndex,
			})
			logicPlayer.SetNetPeer(trnode.ETRNodeTypeCellData, &trnode.TRNodeInfo{
				ZoneID:    trframe.GetFrameConfig().ZoneID,
				NodeType:  trnode.ETRNodeTypeCellData,
				NodeIndex: int32(dbIndex),
			})
			// 这里后续处理放到下一帧处理
			logicPlayer.IsOnline = true
			logicPlayer.UpdateHeartTime(timeutil.NowTime())
			trframe.AfterMsgJob(func() {
				logicPlayer.UpdateHeartTime(timeutil.NowTime())
				logicGlobal.HandlePlayerOnline(logicPlayer)
			})
			rep.RoleData = logicPlayer.ToClientRoleInfo()
			rep.DataInfo = &pbserver.ServerNodeInfo{
				ZoneID:    trframe.GetFrameConfig().ZoneID,
				NodeType:  trnode.ETRNodeTypeCellData,
				NodeIndex: int32(dbIndex),
			}
			rep.LogicInfo = trframe.GetCurPBNodeInfo()
			trframe.SendReplyMessage(okCode,
				rep,
				env)
			loghlp.Infof("user(%d)[%s][%s] logingame succ finshi!",
				cbRep.RoleID,
				req.Account,
				logicPlayer.GetBaseData().GetRoleName())
		} else {
			loghlp.Errorf("logingame fail!!!")
			trframe.SendReplyErrorMessage(okCode,
				msgData,
				env)
		}
	}
	loadReq := &pbserver.ESMsgPlayerLoadRoleReq{
		RoleID:   req.RoleID,
		IsLogin:  true,
		UserID:   req.UserID,
		UserName: req.Account,
	}
	trframe.ForwardZoneMessage(protocol.ESMsgClassPlayer,
		protocol.ESMsgPlayerLoadRole,
		loadReq,
		trnode.ETRNodeTypeCellData,
		int32(dbIndex),
		cbDo,
		trframe.MakeMsgEnv2(tmsgCtx.Session, tmsgCtx.NetMessage),
	)
	return protocol.ECodeSuccess, rep, iframe.EHandlePending
}

// 心跳
func HandleECMsgPlayerKeepHeart(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbclient.ECMsgPlayerKeepHeartReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, pbtools.MakeErrorParams("pberror"), iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	rep := &pbclient.ECMsgPlayerKeepHeartRsp{}
	logicGlobal := servLogic.GetLogicGlobal()
	logicPlayer := logicGlobal.FindPlayer(tmsgCtx.NetMessage.SecondHead.ID)
	if logicPlayer == nil {
		loghlp.Errorf("heart, but logic player(%d) not find", tmsgCtx.NetMessage.SecondHead.ID)
		return protocol.ECodeSuccess, rep, iframe.EHandleNone
	}
	loghlp.Debugf("logic player(%d) update heart time", logicPlayer.RoleID)
	logicPlayer.UpdateHeartTime(trframe.GetFrameSysNowTime())
	trframe.PushZoneClientPBMessage(protocol.ECMsgClassPlayer,
		protocol.ECMsgPlayerKeepHeart,
		req,
		trnode.ETRNodeTypeCellData,
		logicPlayer.GetNetPeerIndex(trnode.ETRNodeTypeCellData),
		logicPlayer.GetRoleID(),
	)
	return protocol.ECodeSuccess, rep, iframe.EHandleNone
}

// 离线
func HandleESMsgPlayerDisconnect(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbserver.ESMsgPlayerDisconnectReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, pbtools.MakeErrorParams("pberror"), iframe.EHandleContent
	}
	rep := &pbserver.ESMsgPlayerDisconnectRep{}
	loghlp.Infof("player(%d) disconnect", tmsgCtx.NetMessage.SecondHead.ID)
	logicPlayer := servLogic.GetLogicGlobal().FindPlayer(tmsgCtx.NetMessage.SecondHead.ID)
	if logicPlayer != nil {
		servLogic.GetLogicGlobal().HandlePlayerOffline(logicPlayer)
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}
