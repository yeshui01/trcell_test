/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-06-15 14:14:17
 * @LastEditTime: 2022-09-23 14:24:50
 * @Brief:当前工作节点
 */
package trframe

import (
	"math/rand"
	"reflect"
	"trcell/pkg/evhub"
	"trcell/pkg/loghlp"
	"trcell/pkg/pb/pbframe"
	"trcell/pkg/pb/pbtools"
	"trcell/pkg/protocol"
	"trcell/pkg/trframe/iframe"
	"trcell/pkg/trframe/trnode"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// 发送消息(本区内)
func (tf *TRFrame) forwardZoneMessage(msgClass int32, msgType int32, msgData []byte, nodeTpye int32, nodeIndex int32, cb iframe.MsgCallbackFunc, env *iframe.TRRemoteMsgEnv) bool {
	sendMsg := MakeInnerMsg(msgClass, msgType, msgData)
	// 找到节点
	frameNode, err := tf.frameNodeMgr.FindNode(tf.frameConfig.ZoneID, nodeTpye, nodeIndex)
	if frameNode == nil || err != nil {
		loghlp.Errorf("forward message fail, not find node(%d,%d)", nodeTpye, nodeIndex)
		return false
	}
	callInfo := tf.remoteMsgMgr.makeCallInfo(msgClass, msgType, cb, env)
	setSecondHead(sendMsg, 0, callInfo.reqID, 0)
	frameNode.SendMsg(sendMsg)
	loghlp.Debugf("forwardMessage,toNode(%s),msg(%d_%d),reqID:%d",
		frameNode.GetNodeInfo().DesInfo,
		msgClass,
		msgType,
		callInfo.reqID,
	)
	return true
}

func (tf *TRFrame) ForwardMessage(msgClass int32, msgType int32, pbMsg protoreflect.ProtoMessage, nodeTpye int32, nodeIndex int32, cb iframe.MsgCallbackFunc, env *iframe.TRRemoteMsgEnv) bool {
	data, err := proto.Marshal(pbMsg)
	if err != nil {
		loghlp.Errorf("forward message fail, proto marsh error:%s", err.Error())
		return false
	}
	return tf.forwardZoneMessage(msgClass, msgType, data, nodeTpye, nodeIndex, cb, env)
}

func (tf *TRFrame) pushMessage(msgClass int32, msgType int32, msgData []byte, nodeTpye int32, nodeIndex int32) bool {
	sendMsg := MakeInnerMsg(msgClass, msgType, msgData)
	// 找到节点
	frameNode, err := tf.frameNodeMgr.FindNode(tf.frameConfig.ZoneID, nodeTpye, nodeIndex)
	if frameNode == nil || err != nil {
		loghlp.Errorf("push message fail, not find node(%d,%d)", nodeTpye, nodeIndex)
		return false
	}
	frameNode.SendMsg(sendMsg)
	loghlp.Debugf("push message,toNode(%s),msg(%d_%d),reqID:%d",
		frameNode.GetNodeInfo().DesInfo,
		msgClass,
		msgType,
	)
	return true
}

func (tf *TRFrame) PushZoneMessage(msgClass int32, msgType int32, pbMsg protoreflect.ProtoMessage, nodeTpye int32, nodeIndex int32) bool {
	data, err := proto.Marshal(pbMsg)
	if err != nil {
		loghlp.Errorf("PushZoneMessage fail, proto marsh error:%s", err.Error())
		return false
	}
	return tf.pushMessage(msgClass, msgType, data, nodeTpye, nodeIndex)
}

// 发送回复消息
func (tf *TRFrame) SendReplyMessage(okCode int32, pbRep proto.Message, env *iframe.TRRemoteMsgEnv) {
	env.SrcMessage.Head.Result = uint16(okCode)
	if pbRep != nil {
		env.SrcMessage.Data, _ = EncodePBMessage(pbRep)
	}
	if env.SrcMessage.SecondHead != nil {
		env.SrcMessage.SecondHead.RepID = env.SrcMessage.SecondHead.ReqID
	}
	session := tf.evHub.GetSession(env.SrcSessionID)
	if session != nil {
		session.Send(env.SrcMessage)
		loghlp.Debugf("SendReplyMessage(%d_%d)[%s]:%+v",
			env.SrcMessage.Head.MsgClass,
			env.SrcMessage.Head.MsgType,
			pbtools.GetFullNameByMessage(pbRep),
			pbRep,
		)
	} else {
		loghlp.Errorf("SendReplyMessage(%d_%d),but not find session(%d)",
			env.SrcMessage.Head.MsgClass,
			env.SrcMessage.Head.MsgType,
			env.SrcSessionID,
		)
	}
}

// 发送回复消息
func (tf *TRFrame) SendReplyMessageData(okCode int32, repData []byte, env *iframe.TRRemoteMsgEnv) bool {
	env.SrcMessage.Head.Result = uint16(okCode)
	if repData != nil {
		env.SrcMessage.Data = repData
	}
	if env.SrcMessage.SecondHead != nil {
		env.SrcMessage.SecondHead.RepID = env.SrcMessage.SecondHead.ReqID
	}
	session := tf.evHub.GetSession(env.SrcSessionID)
	if session != nil {
		session.Send(env.SrcMessage)
		loghlp.Debugf("SendReplyMessageData(%d_%d)",
			env.SrcMessage.Head.MsgClass,
			env.SrcMessage.Head.MsgType,
		)
	} else {
		loghlp.Errorf("SendReplyMessage(%d_%d),but not find session(%d)",
			env.SrcMessage.Head.MsgClass,
			env.SrcMessage.Head.MsgType,
			env.SrcSessionID,
		)
		return false
	}

	return true
}

func (tf *TRFrame) SendReplyErrorMessage(okCode int32, repData []byte, env *iframe.TRRemoteMsgEnv) {
	env.SrcMessage.Head.Result = uint16(okCode)
	if repData != nil {
		env.SrcMessage.Data = repData
	}
	if env.SrcMessage.SecondHead != nil {
		env.SrcMessage.SecondHead.RepID = env.SrcMessage.SecondHead.ReqID
	}
	session := tf.evHub.GetSession(env.SrcSessionID)
	if session != nil {
		session.Send(env.SrcMessage)
		loghlp.Debugf("SendReplyErrorMessage(%d_%d)okcode(%d)",
			env.SrcMessage.Head.MsgClass,
			env.SrcMessage.Head.MsgType,
			okCode,
		)
	} else {
		loghlp.Errorf("SendReplyMessage(%d_%d),but not find session(%d)",
			env.SrcMessage.Head.MsgClass,
			env.SrcMessage.Head.MsgType,
			env.SrcSessionID,
		)
	}
}

// 通用接口-发送消息到某个节点
func (tf *TRFrame) ForwardNodePBMessageByNodeUid(msgClass int32, msgType int32, pbMsg proto.Message, nodeUid int64, cb iframe.MsgCallbackFunc, env *iframe.TRRemoteMsgEnv) bool {
	msgData, err := EncodePBMessage(pbMsg)
	if err != nil {
		return false
	}
	return tf.ForwardNodeMessageByNodeUid(msgClass, msgType, msgData, nodeUid, cb, env)
}

// 通用接口-发送消息到某个节点
func (tf *TRFrame) ForwardNodeMessageByNodeUid(msgClass int32, msgType int32, msgData []byte, nodeUid int64, cb iframe.MsgCallbackFunc, env *iframe.TRRemoteMsgEnv) bool {
	zoneID, nodeType, nodeIndex := trnode.GetNodePartIDByNodeUid(nodeUid)
	frameNode, err := tf.frameNodeMgr.FindNode(zoneID, nodeType, nodeIndex)
	if frameNode != nil && err == nil {
		sendMsg := MakeInnerMsg(msgClass, msgType, msgData)
		callInfo := tf.remoteMsgMgr.makeCallInfo(msgClass, msgType, cb, env)
		setSecondHead(sendMsg, 0, callInfo.reqID, 0)
		return frameNode.SendMsg(sendMsg)
	}
	// loghlp.Errorf("not find node(%d,%d,%d)", zoneID, nodeType, nodeIndex)
	// 当前非trans节点
	transMsg := &pbframe.EFrameMsgTransMsgReq{
		MsgClass:    msgClass,
		MsgType:     msgType,
		MsgData:     msgData,
		DestNodeUID: nodeUid,
	}
	transMsgData, errTrans := EncodePBMessage(transMsg)
	if errTrans != nil {
		loghlp.Errorf("encode EFrameMsgTransMsgReq error:%s", errTrans.Error())
		return false
	}
	sendMsg := MakeInnerMsg(protocol.EMsgClassFrame, protocol.EFrameMsgTransMsg, transMsgData)
	curNodeType := tf.curWorkNode.NodeType()
	if curNodeType == trnode.ETRNodeTypeCellTrans || curNodeType != trnode.ETRNodeTypeCellRoot {
		// 转发到目标zone的root节点
		zoneNode, err := tf.frameNodeMgr.FindNode(zoneID, trnode.ETRNodeTypeCellRoot, 0)
		if err == nil {
			callInfo := tf.remoteMsgMgr.makeCallInfo(msgClass, msgType, cb, env)
			setSecondHead(sendMsg, 0, callInfo.reqID, 0)
			return zoneNode.SendMsg(sendMsg)
		} else {
			loghlp.Errorf("not find zone root node(%d,%d,%d)", zoneID, nodeType, nodeIndex)
			return false
		}
	} else if curNodeType == trnode.ETRNodeTypeCellRoot {
		// 转发到目标节点
		if zoneID == tf.frameConfig.ZoneID {
			// 本区其他节点
			return tf.forwardZoneMessage(msgClass, msgType, msgData, nodeType, nodeIndex, cb, env)
		} else {
			// 发送到trans节点
			transList := tf.frameNodeMgr.GetNodeListByType(trnode.ETRNodeTypeCellTrans)
			if len(transList) > 0 {
				// 随机取一个
				transIndex := rand.Intn(len(transList))
				transNode := transList[transIndex]
				callInfo := tf.remoteMsgMgr.makeCallInfo(msgClass, msgType, cb, env)
				setSecondHead(sendMsg, 0, callInfo.reqID, 0)
				transNode.SendMsg(sendMsg)
				return true
			} else {
				loghlp.Errorf("not find trans node(%d,%d,%d)", zoneID, nodeType, nodeIndex)
				return false
			}
		}
	}
	return true
}

func MakePBMessage(msgClass int32, msgType int32, pbMsg protoreflect.ProtoMessage, okCode int32) *evhub.NetMessage {
	data, err := proto.Marshal(pbMsg)
	if err == nil {
		evMsg := evhub.MakeMessage(msgClass, msgType, data)
		return evMsg
	} else {
		loghlp.Errorf("make pb message error:%s", err.Error())
	}
	return nil
}

func DecodePBMessage(sMsg *evhub.NetMessage, pbMsg protoreflect.ProtoMessage) bool {
	errParse := proto.Unmarshal(sMsg.Data, pbMsg)
	if errParse != nil {
		loghlp.Errorf("decode PB(%s) error:%s", reflect.TypeOf(pbMsg).String(), errParse.Error())
		return false
	}
	return true
}

func DecodePBMessage2(msgData []byte, pbMsg protoreflect.ProtoMessage) bool {
	errParse := proto.Unmarshal(msgData, pbMsg)
	if errParse != nil {
		loghlp.Errorf("decode PB(%s) error:%s", reflect.TypeOf(pbMsg).String(), errParse.Error())
		return false
	}
	return true
}

func EncodePBMessage(pbMsg proto.Message) ([]byte, error) {
	data, err := proto.Marshal(pbMsg)
	if err != nil {
		loghlp.Errorf("EncodePBMessageError:%s", err.Error())
	}

	return data, err
}

func LogMsgInfo(netMsg *evhub.NetMessage, pbMsg proto.Message) {
	loghlp.Debugf("recv reqmsg(%d_%d)[%s]:%+v",
		netMsg.Head.MsgClass,
		netMsg.Head.MsgType,
		pbtools.GetFullNameByMessage(pbMsg),
		pbMsg,
	)
}

func LogCbMsgInfo(netMsg *evhub.NetMessage, pbMsg proto.Message) {
	loghlp.Debugf("recv cbmsg(%d_%d)[%s]:%+v",
		netMsg.Head.MsgClass,
		netMsg.Head.MsgType,
		pbtools.GetFullNameByMessage(pbMsg),
		pbMsg,
	)
}

// 多节点异步消息接口

/**
 * 功能: 创建异步多节点调用实例
 * @finalDo : 最终的行为操作
 * @return: 多节点异步调用实例
 */
func CreateAsyncMultiCall(finalDo iframe.MultiMsgCallbackFunc) *MultiMsgCall {
	return frameCore.multiMsgCallMgr.CreateMultiMsgCall(finalDo)
}

/**
 * 功能: 开始执行多节点异步串行化操作
 * @mmc : 要开始执行的多节点异步调用实例
 * @env : 远程调用上下文环境
 * @return: true-succ false-fail
 */
func AsyncSerialCall(mmc *MultiMsgCall, env *iframe.TRRemoteMsgEnv) bool {
	return frameCore.multiMsgCallMgr.AsyncSerialCall(mmc, env)
}
