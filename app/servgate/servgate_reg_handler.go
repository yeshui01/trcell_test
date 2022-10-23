/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 14:59:26
 * @LastEditTime: 2022-09-23 16:52:05
 * @FilePath: \trcell\app\servgate\servgate_reg_handler.go
 */
package servgate

import (
	"trcell/app/servgate/servgatehandler"
	"trcell/pkg/protocol"
	"trcell/pkg/trframe"
)

func (serv *CellServGate) RegisterMsgHandler() {
	// frame
	trframe.RegWorkMsgHandler(
		protocol.EMsgClassFrame,
		protocol.EFrameMsgPushMsgToClient,
		servgatehandler.HandleFramePushMsgToClient,
	)
	trframe.RegWorkMsgHandler(
		protocol.EMsgClassFrame,
		protocol.EFrameMsgBroadcastMsgToClient,
		servgatehandler.HandleFrameBroadcastMsgToClient,
	)
	// player
	trframe.RegWorkMsgHandler(
		protocol.ECMsgClassPlayer,
		protocol.ECMsgPlayerCreateRole,
		servgatehandler.HandleECMsgPlayerCreateRole)

	trframe.RegWorkMsgHandler(
		protocol.ECMsgClassPlayer,
		protocol.ECMsgPlayerKeepHeart,
		servgatehandler.HandlePlayerHeart,
	)
	trframe.RegWorkMsgHandler(
		protocol.ECMsgClassPlayer,
		protocol.ECMsgPlayerLoginGame,
		servgatehandler.HandleECMsgPlayerLoginGame,
	)
	// trframe.RegWorkMsgHandler(
	// 	protocol.ESMsgClassPlayer,
	// 	protocol.ESMsgPlayerKickOut,
	// 	servgatehandler.HandlePlayerKickout,
	// )

}
