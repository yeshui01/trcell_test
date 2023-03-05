/**
 * @author: [mknight]
 * @email : [824338670@qq.com]
 * @create:	2023-01-31 15:23:23
 * @modify:	2023-01-31 15:23:23
 * @desc  :	[多个消息调用回调处理]
 */
package trframe

import (
	"trcell/pkg/loghlp"
	"trcell/pkg/timeutil"
	"trcell/pkg/trframe/iframe"

	"google.golang.org/protobuf/proto"
)

type MultiCallReqPBMaker func(callInsHolder interface{}) (msgClass, msgType int32, pbMsg proto.Message)

// call step
type MultiMsgStep struct {
	ReqPBMsgMaker MultiCallReqPBMaker             // 要发送的消息生成器 callInsHolder用来指定MultiMsgCall
	RemoteNodeUID int64                           // 远程目标节点
	StepCallback  iframe.MultiMsgStepCallbackFunc // 用户的步骤回调接口
	BeginTime     int64                           // 开始执行的时间戳,单位:毫秒
	// 消息结果记录
	MsgClass int32  // 此阶段发送的消息class
	MsgType  int32  // 此阶段发送的消息Type
	OKCode   int32  // 回调code
	MsgData  []byte // 回调消息数据
}

// multi call
type MultiMsgCall struct {
	MultiMsgID  int64 // id
	CallSteps   []*MultiMsgStep
	FinalHandle iframe.MultiMsgCallbackFunc // 最终处理, callInsHolder用来指定MultiMsgCall
	StepIdx     int32                       // 步骤索引,0开始,指示当前在哪一个步骤
	BeginTime   int64                       // 开始时间,单位:ms
}

func newMultiMsgCall(id int64, finalDo iframe.MultiMsgCallbackFunc) *MultiMsgCall {
	return &MultiMsgCall{
		MultiMsgID:  id,
		CallSteps:   make([]*MultiMsgStep, 0),
		FinalHandle: finalDo,
		StepIdx:     0,
	}
}

func (mmc *MultiMsgCall) AddStep(pbMaker MultiCallReqPBMaker, nodeUID int64, stepCb iframe.MultiMsgStepCallbackFunc) {
	stepCall := &MultiMsgStep{
		ReqPBMsgMaker: pbMaker,
		RemoteNodeUID: nodeUID,
		StepCallback:  stepCb,
		BeginTime:     0,
	}
	mmc.CallSteps = append(mmc.CallSteps, stepCall)
}

// -------------- mgr ---------------------------------
type MultiMsgCallMgr struct {
	multiIdx      int64
	multiCallList map[int64]*MultiMsgCall
	frameCore     *TRFrame
	lastUpdTime   int64 // 单位: 秒
}

func newMultiMsgCallMgr(frameObj *TRFrame) *MultiMsgCallMgr {
	return &MultiMsgCallMgr{
		frameCore:     frameObj,
		multiIdx:      0,
		multiCallList: make(map[int64]*MultiMsgCall),
		lastUpdTime:   0,
	}
}

func (mmcMgr *MultiMsgCallMgr) genMultiID() int64 {
	if mmcMgr.multiIdx >= 0x00ffffffffffffff {
		mmcMgr.multiIdx = 0
	}
	mmcMgr.multiIdx++
	return mmcMgr.multiIdx
}

func (mmcMgr *MultiMsgCallMgr) CreateMultiMsgCall(finalDo iframe.MultiMsgCallbackFunc) *MultiMsgCall {
	multiID := mmcMgr.genMultiID()
	mmc := newMultiMsgCall(multiID, finalDo)
	mmcMgr.multiCallList[multiID] = mmc

	return mmc
}

// 异步多步骤串行化调用
func (mmcMgr *MultiMsgCallMgr) AsyncSerialCall(mmc *MultiMsgCall, env *iframe.TRRemoteMsgEnv) bool {
	if mmc == nil {
		return false
	}
	mmc.StepIdx = 0
	mmc.BeginTime = timeutil.NowTimeMs()
	return mmcMgr.asyncSerialStepCall(mmc, env)
}

// 单个步骤调用
func (mmcMgr *MultiMsgCallMgr) asyncSerialStepCall(mmc *MultiMsgCall, env *iframe.TRRemoteMsgEnv) bool {
	if mmc.StepIdx >= int32(len(mmc.CallSteps)) {
		return false
	}
	callStep := mmc.CallSteps[mmc.StepIdx]
	if callStep == nil {
		return false
	}
	if callStep.ReqPBMsgMaker == nil || callStep.StepCallback == nil {
		panic("callstep param error!!!")
	}

	msgClass, msgType, pbReq := callStep.ReqPBMsgMaker(mmc)

	// 请求记录
	callStep.MsgClass = msgClass
	callStep.MsgType = msgType
	callStep.BeginTime = timeutil.NowTimeMs()

	nodeUID := callStep.RemoteNodeUID
	frameObj := mmcMgr.frameCore
	cbInner := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
		oneStep := mmc.CallSteps[mmc.StepIdx]
		oneStep.OKCode = okCode
		oneStep.MsgData = msgData
		if !oneStep.StepCallback(mmc, okCode, msgData, env) {
			// 返回了false,不用后续处理,直接finalHandle
			mmc.FinalHandle(mmc, env)
			delete(mmcMgr.multiCallList, mmc.MultiMsgID)
			return
		}
		mmc.StepIdx++
		if mmc.StepIdx >= int32(len(mmc.CallSteps)) {
			// 所有的调用都处理完毕,执行最终处理
			mmc.FinalHandle(mmc, env)
			delete(mmcMgr.multiCallList, mmc.MultiMsgID)
		} else {
			// 下一阶段处理
			mmcMgr.asyncSerialStepCall(mmc, env)
		}
	}
	return frameObj.ForwardNodePBMessageByNodeUid(msgClass, msgType, pbReq, nodeUID, cbInner, env)
}

func (mmcMgr *MultiMsgCallMgr) update(curTimeMs int64) {
	if curTimeMs-mmcMgr.lastUpdTime < 5000 {
		return
	}
	mmcMgr.lastUpdTime = curTimeMs
	// 检测超时
	for k, v := range mmcMgr.multiCallList {
		if curTimeMs-v.BeginTime >= 15000 {
			errStep := v.CallSteps[v.StepIdx]
			loghlp.Errorf("remote multi call timeout, step(%d), msg(%d_%d) multiId:%d,beginTime(%d)", v.StepIdx, errStep.MsgClass, errStep.MsgType, k, v.BeginTime)
			delete(mmcMgr.multiCallList, k)
		}
	}
}
