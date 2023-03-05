/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-30 11:00:46
 * @LastEditTime: 2022-09-30 11:53:37
 * @FilePath: \trcell\app\servdata\servdatahandler\servdata_servdata_handler.go
 */
package servdatahandler

import (
	"trcell/app/servdata/servdatamain"
	"trcell/pkg/loghlp"
	"trcell/pkg/ormdef"
	"trcell/pkg/pb/pbserver"
	"trcell/pkg/pb/pbtools"
	"trcell/pkg/protocol"
	"trcell/pkg/tbobj"
	"trcell/pkg/trframe"
	"trcell/pkg/trframe/iframe"
)

func HandleESMsgServDataLoadTables(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbserver.ESMsgServDataLoadTablesReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, pbtools.MakeErrorParams("pberror"), iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	rep := &pbserver.ESMsgServDataLoadTablesRep{
		DataTable: &pbserver.ServDataTable{},
	}
	//
	dataGlobal := servData.GetDataGlobal()
	if req.LoadOwner == "center" {
		// center load
		dataGlobal.LoadCenterTables(rep.DataTable, tmsgCtx)
	}
	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}

func HandleESMsgServDataSaveTables(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbserver.ESMsgServDataSaveTablesReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, pbtools.MakeErrorParams("pberror"), iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	rep := &pbserver.ESMsgServDataSaveTablesRep{}
	//
	dataGlobal := servData.GetDataGlobal()
	gameDB := dataGlobal.GetGameDB()
	for _, v := range req.DataTable.ServTables {
		switch v.TableID {
		case ormdef.ETableCsGlobal:
			{
				dataGlobal.DataTbCsGlobal.FromBytes(v.Data)
				oneTable := tbobj.NewTbCsGlobal()
				oneTable.FromBytes(v.Data)
				// 发送到db线程更新
				dbJob := func() bool {
					errDB := gameDB.Model(oneTable.GetOrmMeta()).Select("*").Updates(oneTable.GetOrmMeta()).Error
					if errDB != nil {
						loghlp.Errorf("save table error:%s", errDB.Error())
						return false
					}
					return true
				}
				dataGlobal.PostDBJob(&servdatamain.DataDBJob{
					DoJob: dbJob,
				})
				break
			}
		}
	}
	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}
