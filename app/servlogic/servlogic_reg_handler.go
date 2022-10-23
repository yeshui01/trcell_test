/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 14:15:34
 * @LastEditTime: 2022-09-26 18:12:57
 * @FilePath: \trcell\app\servlogic\servlogic_reg_handler.go
 */
package servlogic

import (
	"trcell/app/servlogic/servlogichandler"
	"trcell/pkg/protocol"
	"trcell/pkg/trframe"
)

func (serv *CellServLogic) RegisterMsgHandler() {
	// player
	trframe.RegWorkMsgHandler(
		protocol.ESMsgClassPlayer,
		protocol.ESMsgPlayerCreateRole,
		servlogichandler.HandleESMsgPlayerCreateRole)
	trframe.RegWorkMsgHandler(
		protocol.ECMsgClassPlayer,
		protocol.ECMsgPlayerKeepHeart,
		servlogichandler.HandleECMsgPlayerKeepHeart)
	trframe.RegWorkMsgHandler(
		protocol.ESMsgClassPlayer,
		protocol.ESMsgPlayerLoginGame,
		servlogichandler.HandleESMsgPlayerLoginGame)

	trframe.RegWorkMsgHandler(
		protocol.ESMsgClassPlayer,
		protocol.ESMsgPlayerDisconnect,
		servlogichandler.HandleESMsgPlayerDisconnect)

}
