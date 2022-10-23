/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-20 11:25:49
 * @LastEditTime: 2022-09-20 14:55:33
 * @FilePath: \trcell\app\cellrobot\robotcore\cellrobot_dispatch_msg.go
 */
package robotcore

func genHandleKey(msgClass int32, msgType int32) int32 {
	return msgClass*1000 + msgType
}

func (robotObj *CellRobot) RegMsgHandle(msgClass int32, msgType int32, handleFunc RobotMsgHandler) {
	handleKey := genHandleKey(msgClass, msgType)
	robotObj.MsgHandlers[handleKey] = handleFunc
}

func (robotObj *CellRobot) DispatchRobotMsg(clientMsg *ClientMessage) {
	handleKey := genHandleKey(int32(clientMsg.Head.MainType), int32(clientMsg.Head.SubType))
	if clientMsg.Head.ID > 0 {
		if callEnv, ok := robotObj.AsyncCall[int32(clientMsg.Head.ID)]; ok {
			robotObj.Debugf("handle robot callback msg(%d_%d),ack_seqid(%d)", clientMsg.Head.MainType, clientMsg.Head.SubType, clientMsg.Head.ID)
			callEnv.CallbackFunc(clientMsg)
			delete(robotObj.AsyncCall, int32(clientMsg.Head.ID))
		} else {
			robotObj.Errorf("not find robot callenv, callback msg(%d_%d),seqid(%d)", clientMsg.Head.MainType, clientMsg.Head.SubType, clientMsg.Head.ID)
		}
		return
	}
	if handleFunc, ok := robotObj.MsgHandlers[handleKey]; ok {
		if robotObj.RealRobotIns != nil {
			handleFunc(robotObj.RealRobotIns, clientMsg)
		} else {
			handleFunc(robotObj, clientMsg)
		}
	} else {
		robotObj.Debugf("not find robot msg(%d_%d) handler", clientMsg.Head.MainType, clientMsg.Head.SubType)
	}
}
