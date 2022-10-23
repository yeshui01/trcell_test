/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 11:10:06
 * @LastEditTime: 2022-09-30 11:32:59
 * @FilePath: \trcell\app\servdata\servdata_reg_handler.go
 */
package servdata

import (
	"trcell/app/servdata/servdatahandler"
	"trcell/pkg/protocol"
	"trcell/pkg/trframe"
)

func (serv *CellServData) RegisterMsgHandler() {
	// player
	trframe.RegWorkMsgHandler(protocol.ESMsgClassPlayer,
		protocol.ESMsgPlayerCreateRole,
		servdatahandler.HandleESMsgPlayerCreateRole)
	trframe.RegWorkMsgHandler(protocol.ESMsgClassPlayer,
		protocol.ESMsgPlayerLoadRole,
		servdatahandler.HandleESMsgPlayerLoadRoleData)
	trframe.RegWorkMsgHandler(protocol.ESMsgClassPlayer,
		protocol.ESMsgPlayerSaveRole,
		servdatahandler.HandleESMsgPlayerSaveRoleData)

	trframe.RegWorkMsgHandler(protocol.ECMsgClassPlayer,
		protocol.ECMsgPlayerKeepHeart,
		servdatahandler.HandleECMsgPlayerKeepHeart)

	trframe.RegWorkMsgHandler(protocol.ESMsgClassServData,
		protocol.ESMsgServDataLoadTables,
		servdatahandler.HandleESMsgServDataLoadTables)
}
