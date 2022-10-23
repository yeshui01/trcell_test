/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-20 14:24:24
 * @LastEditTime: 2022-09-20 14:25:55
 * @FilePath: \trcell\app\servgate\servgatehandler\servgate_frame_handler.go
 */
package servgatehandler

import (
	"trcell/pkg/evhub"
	"trcell/pkg/loghlp"
	"trcell/pkg/pb/pbframe"
	"trcell/pkg/protocol"
	"trcell/pkg/trframe"
	"trcell/pkg/trframe/iframe"
)

func HandleServerNodeRegister(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	// TODO:
	return protocol.ECodeSuccess, nil, iframe.EHandleContent
}

// 推送消息给客户端
func HandleFramePushMsgToClient(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbframe.EFrameMsgPushMsgToClientReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleNone
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	gateUser := servGate.GetUserManager().GetGateUser(req.RoleID)
	if gateUser != nil {
		pushMessage := evhub.MakeMessage(req.MsgClass, req.MsgType, req.MsgData)
		gateUser.SendMessageToSelf(pushMessage)
		loghlp.Debugf("push player(%d) msg(%d_%d)", req.RoleID, req.MsgClass, req.MsgType)
	}
	return protocol.ECodeSuccess, nil, iframe.EHandleNone
}

// 广播消息给客户端
func HandleFrameBroadcastMsgToClient(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbframe.EFrameMsgBroadcastMsgToClientReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleNone
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	for _, RoleID := range req.RoleList {
		gateUser := servGate.GetUserManager().GetGateUser(RoleID)
		if gateUser != nil {
			pushMessage := evhub.MakeMessage(req.MsgClass, req.MsgType, req.MsgData)
			gateUser.SendMessageToSelf(pushMessage)
			loghlp.Debugf("broadcast push player(%d) msg(%d_%d)", RoleID, req.MsgClass, req.MsgType)
		}
	}
	return protocol.ECodeSuccess, nil, iframe.EHandleNone
}
