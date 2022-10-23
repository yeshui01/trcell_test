package tframedispatcher

import (
	"trcell/pkg/evhub"
	"trcell/pkg/loghlp"
	"trcell/pkg/pb/pbtools"
	"trcell/pkg/protocol"
	"trcell/pkg/trframe/iframe"
	"trcell/pkg/trframe/trnode"

	"google.golang.org/protobuf/reflect/protoreflect"
)

type FrameMsgDispatcher struct {
	frameInstance iframe.ITRFrame
	handlerList   map[int32]iframe.FrameMsgHandler
}

func (dsp *FrameMsgDispatcher) RegisterMsgHandler(msgClass int32, msgType int32, msgHander iframe.FrameMsgHandler) {
	handleKey := dsp.genHandlerKey(msgClass, msgType)
	dsp.handlerList[handleKey] = msgHander
}

func (dsp *FrameMsgDispatcher) Dispatch(session iframe.ISession, msg *evhub.NetMessage, customData interface{}) bool {
	handleKey := dsp.genHandlerKey(int32(msg.Head.MsgClass), int32(msg.Head.MsgType))
	msgHandler, ok := dsp.handlerList[handleKey]
	if ok {
		loghlp.Debugf("DispatchMsg(%d_%d)",
			msg.Head.MsgClass,
			msg.Head.MsgType,
		)
		tmsgCtx := &iframe.TMsgContext{
			FrameInstance: dsp.frameInstance,
			Session:       session,
			NetMessage:    msg,
			CustomData:    customData,
		}
		okCode, retData, rt := msgHandler(tmsgCtx)
		if rt == iframe.EHandleContent {
			msg.Head.Result = uint16(okCode)
			if okCode == protocol.ECodeSuccess {
				pbRep, ok := retData.(protoreflect.ProtoMessage)
				if ok {
					loghlp.Debugf("reply repmsg(%d_%d),isok(%d)[%s]:%+v",
						msg.Head.MsgClass,
						msg.Head.MsgType,
						okCode,
						pbtools.GetFullNameByMessage(pbRep),
						pbRep,
					)
					// pb
					data, err := iframe.ToPbData(pbRep)
					if err == nil {
						msg.Data = data
					}
					if msg.Head.HasSecond > 0 {
						msg.SecondHead.RepID = msg.SecondHead.ReqID
					}
				} else {
					loghlp.Errorf("retData to pb message(%d,%d) fail!!!!!", msg.Head.MsgClass, msg.Head.MsgType)
					msg.Head.Result = protocol.ECodeSysError
					loghlp.Debugf("reply repmsg2(%d_%d),isok(%d)[%s]",
						msg.Head.MsgClass,
						msg.Head.MsgType,
						okCode,
						pbtools.GetFullNameByMessage(pbRep),
					)
					if msg.Head.HasSecond > 0 {
						msg.SecondHead.RepID = msg.SecondHead.ReqID
					}
				}
			} else {
				pbRep, ok := retData.(protoreflect.ProtoMessage)
				if ok {
					loghlp.Debugf("reply exception repmsg(%d_%d),isok(%d)[%s]:%+v",
						msg.Head.MsgClass,
						msg.Head.MsgType,
						okCode,
						pbtools.GetFullNameByMessage(pbRep),
						pbRep,
					)
					// pb
					data, err := iframe.ToPbData(pbRep)
					if err == nil {
						msg.Data = data
					}
					if msg.Head.HasSecond > 0 {
						msg.SecondHead.RepID = msg.SecondHead.ReqID
					}
				} else {
					//msg.Head.Result = protocol.ECodeSysError
					loghlp.Debugf("reply exception repmsg2(%d_%d),isok(%d)[%s]",
						msg.Head.MsgClass,
						msg.Head.MsgType,
						okCode,
						pbtools.GetFullNameByMessage(pbRep),
					)
					if msg.Head.HasSecond > 0 {
						msg.SecondHead.RepID = msg.SecondHead.ReqID
					}
				}
			}
			session.SendMsg(msg)
		}
	} else {
		if dsp.frameInstance.GetCurNodeType() != trnode.ETRNodeTypeCellGate {
			loghlp.Warnf("msg(%d_%d) dispatch fail, not find handler",
				msg.Head.MsgClass,
				msg.Head.MsgType)
		}
		return false
	}
	return true
}

func (dsp *FrameMsgDispatcher) genHandlerKey(msgClass int32, msgType int32) int32 {
	return msgClass*1000 + msgType
}

func NewFrameMsgDispatcher(frameObj iframe.ITRFrame) *FrameMsgDispatcher {
	return &FrameMsgDispatcher{
		frameInstance: frameObj,
		handlerList:   make(map[int32]iframe.FrameMsgHandler, 0),
	}
}
