/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-16 14:18:14
 * @LastEditTime: 2022-09-16 17:51:37
 * @FilePath: \trcell\pkg\trframe\trframe_handler.go
 */
package trframe

import (
	"trcell/pkg/loghlp"
	"trcell/pkg/pb/pbframe"
	"trcell/pkg/protocol"
	"trcell/pkg/trframe/iframe"
	"trcell/pkg/trframe/trnode"

	"google.golang.org/protobuf/proto"
)

// 协议处理-注册节点信息
func handleRegisterNodeInfo(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbframe.FrameMsgRegisterServerInfoReq{}
	rep := &pbframe.FrameMsgRegisterServerInfoRep{}
	proto.Unmarshal(tmsgCtx.NetMessage.Data, req)
	frameSession := tmsgCtx.Session.(*FrameSession)
	frameSession.nodeType = req.NodeType
	// 关联节点信息
	frameSession.nodeInfo = &trnode.TRNodeInfo{
		ZoneID:    req.ZoneID,
		NodeType:  req.NodeType,
		NodeIndex: req.NodeIndex,
		DesInfo:   req.NodeDes,
	}
	frameCore.frameNodeMgr.AddNode(frameSession)
	loghlp.Infof("handleRegisterNodeInfo,req:%+v", req)
	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}

// 协议处理-转发节点信息
func handleTransNodeMsg(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbframe.EFrameMsgTransMsgReq{}
	rep := &pbframe.EFrameMsgTransMsgRep{}
	err := proto.Unmarshal(tmsgCtx.NetMessage.Data, req)
	if err != nil {
		return protocol.ECodeSysError, rep, iframe.EHandleContent
	}
	loghlp.Infof("handleTransNodeMsg,req:%+v", req)
	frameSession := tmsgCtx.Session.(*FrameSession)
	srcSessionID := frameSession.GetSessionID()
	cbEnv := MakeMsgEnv(srcSessionID, tmsgCtx.NetMessage)
	// zoneID, nodeType, nodeIndex := trnode.GetNodePartIDByNodeUid(req.DestNodeUID)
	cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
		loghlp.Infof("trans success,okCode:%d", okCode)
		SendReplyMessageData(okCode, msgData, env)
	}

	forwardRet := ForwardNodeMessage(
		req.MsgClass,
		req.MsgType,
		req.MsgData,
		req.DestNodeUID,
		cb,
		cbEnv)

	if !forwardRet {
		return protocol.ECodeNotFindServNode, rep, iframe.EHandleContent
	}

	return protocol.ECodeAsyncHandle, nil, iframe.EHandlePending
}
