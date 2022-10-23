/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-20 11:12:30
 * @LastEditTime: 2022-09-20 11:12:34
 * @FilePath: \trcell\app\servgate\servgate_forward.go
 */
package servgate

// 是否直接转发到View
func (serv *CellServGate) IsForwardToView(msgClass int32, msgType int32) bool {
	msgKey := msgClass*MsgFactor + msgType
	switch msgKey {
	// case protocol.ECMsgClassGame*MsgFactor + protocol.ECMsgGameReadyOpt:
	// 	return true
	}
	return false
}

// 是否直接转发到Center
func (serv *CellServGate) IsForwardToCenter(msgClass int32, msgType int32) bool {
	msgKey := msgClass*MsgFactor + msgType
	switch msgKey {
	// case protocol.ECMsgClassGame*MsgFactor + protocol.ECMsgGameReadyOpt:
	// 	return true
	}
	return false
}

// 是否直接转发到Game
func (serv *CellServGate) IsForwardToGame(msgClass int32, msgType int32) bool {
	msgKey := msgClass*MsgFactor + msgType
	switch msgKey {
	// case protocol.ECMsgClassGame*MsgFactor + protocol.ECMsgGameReadyOpt:
	// 	return true
	}
	return false
}
