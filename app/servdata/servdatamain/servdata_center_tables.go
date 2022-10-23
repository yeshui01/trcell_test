/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-30 11:12:07
 * @LastEditTime: 2022-09-30 11:53:24
 * @FilePath: \trcell\app\servdata\servdatamain\servdata_center_tables.go
 */
package servdatamain

import (
	"trcell/pkg/loghlp"
	"trcell/pkg/ormdef"
	"trcell/pkg/pb/pbserver"
	"trcell/pkg/trframe/iframe"
)

// 加载cs相关的数据表-只用与消息处理,其他地方不要调用
func (servGlobal *ServDataGlobal) LoadCenterTables(tableData *pbserver.ServDataTable, tmsgCtx *iframe.TMsgContext) {
	gameDB := servGlobal.gameDB
	{
		// cs_global
		servGlobal.DataTbCsGlobal.SetID(1)
		ormMeta := servGlobal.DataTbCsGlobal.GetOrmMeta()
		err := gameDB.Model(ormMeta).First(ormMeta).Error
		if err == nil {
			loghlp.Debug("load cs_global succ")
			oneTable := &pbserver.DbTableData{
				TableID: ormdef.ETableCsGlobal,
			}
			oneTable.Data, _ = servGlobal.DataTbCsGlobal.ToBytes()
			tableData.ServTables = append(tableData.ServTables, oneTable)
		}
	}
}
