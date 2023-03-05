/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 11:49:33
 * @LastEditTime: 2022-09-23 16:57:04
 * @FilePath: \trcell\app\servcenter\servcenter_reg_handler.go
 */
package servcenter

import (
	"trcell/app/servcenter/servcenterhandler"
	"trcell/pkg/protocol"
	"trcell/pkg/trframe"
)

func (serv *CellServCenter) RegisterMsgHandler() {
	// player
	trframe.RegWorkMsgHandler(
		protocol.ESMsgClassPlayer,
		protocol.ESMsgPlayerCreateRole,
		servcenterhandler.HandleESMsgPlayerCreateRole)

	trframe.RegWorkMsgHandler(
		protocol.ESMsgClassPlayer,
		protocol.ESMsgPlayerLoginGame,
		servcenterhandler.HandleESMsgPlayerLoginGame)

	trframe.RegWorkMsgHandler(
		protocol.ESMsgClassPlayer,
		protocol.ESMsgPlayerDisconnect,
		servcenterhandler.HandleESMsgPlayerDisconnect)
}
