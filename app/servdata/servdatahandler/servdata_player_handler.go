/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-22 17:05:50
 * @LastEditTime: 2022-10-17 11:37:15
 * @FilePath: \trcell\app\servdata\servdatahandler\servdata_player_handler.go
 */
package servdatahandler

import (
	"trcell/app/servdata/servdatamain"
	"trcell/pkg/loghlp"
	"trcell/pkg/ormdef"
	"trcell/pkg/pb/pbclient"
	"trcell/pkg/pb/pbserver"
	"trcell/pkg/pb/pbtools"
	"trcell/pkg/protocol"
	"trcell/pkg/tbobj"
	"trcell/pkg/timeutil"
	"trcell/pkg/trframe"
	"trcell/pkg/trframe/iframe"
)

// 创角-data
func HandleESMsgPlayerCreateRole(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbserver.ESMsgPlayerCreateRoleReq{}
	rep := &pbserver.ESMsgPlayerCreateRoleRep{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	dataGlobal := servData.GetDataGlobal()
	gameDB := dataGlobal.GetGameDB()
	// 查询DB
	tbRoleBase := tbobj.NewTbRoleBase()
	tbRoleBase.SetRoleName(req.Nickname)
	errDB := gameDB.Model(tbRoleBase.GetOrmMeta()).Where("role_name=?", req.Nickname).First(tbRoleBase.GetOrmMeta()).Error
	if errDB == nil {
		return protocol.ECodeRoleNickNameExisted,
			pbtools.MakeErrorParams("nickname existed error"),
			iframe.EHandleContent
	}
	// 不存在,创建
	tbRoleBase.SetUserID(req.UserID)
	tbRoleBase.SetRoleName(req.Nickname)
	tbRoleBase.SetCreateTime(timeutil.NowTime())
	tbRoleBase.SetLevel(1)

	errDB = gameDB.Model(tbRoleBase.GetOrmMeta()).Create(tbRoleBase.GetOrmMeta()).Error
	if errDB != nil {
		loghlp.Errorf("create role error:%s", errDB.Error())
		return protocol.ECodeRoleNickNameExisted,
			pbtools.MakeErrorParams("create role dberror"),
			iframe.EHandleContent
	} else {
		loghlp.Infof("create role succ:%+v", tbRoleBase.GetOrmMeta())
	}
	dataPlayer := servdatamain.NewDataPlayer()
	dataPlayer.RoleID = tbRoleBase.GetRoleID()
	dataPlayer.DataTbRoleBase = tbRoleBase
	// ---  其他数据表初始化 begin ----
	// TODO
	servData.GetDataGlobal().CheckInitPlayerTables(dataPlayer)
	// ---  其他数据表初始化 end ----

	dataGlobal.AddDataPlayer(dataPlayer.RoleID, dataPlayer)

	// 返回数据
	rep.RoleData = &pbclient.RoleInfo{
		RoleID: dataPlayer.RoleID,
	}
	rep.RoleDetail = &pbserver.GameRoleData{}
	dataGlobal.FillMsgRoleDetailData(dataPlayer, rep.RoleDetail)
	// ---  其他数据表返回 end ----
	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}

func HandleESMsgPlayerLoadRoleData(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbserver.ESMsgPlayerLoadRoleReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, pbtools.MakeErrorParams("pberror"), iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	rep := &pbserver.ESMsgPlayerLoadRoleRep{}
	//
	dataGlobal := servData.GetDataGlobal()
	gameDB := dataGlobal.GetGameDB()

	dataPlayer := dataGlobal.FindDataPlayer(req.RoleID)
	tbRoleBase := tbobj.NewTbRoleBase()
	if dataPlayer == nil {
		// db加载
		tbRoleBase.SetRoleID(req.RoleID)
		var err error = nil
		if req.RoleID == 0 && req.IsLogin {
			// 根据账号id查找
			err = gameDB.Model(tbRoleBase.GetOrmMeta()).Where("user_id=?", req.UserID).First(tbRoleBase.GetOrmMeta()).Error
		} else {
			// 根据角色id查找
			tbRoleBase.SetRoleID(req.RoleID)
			err = gameDB.Model(tbRoleBase.GetOrmMeta()).First(tbRoleBase.GetOrmMeta()).Error
		}
		if err != nil {
			// 找不到
			return protocol.ECodeRoleNotExisted,
				pbtools.MakeErrorParams("ECodeRoleNotExisted"),
				iframe.EHandleContent
		}
		dataPlayer = servdatamain.NewDataPlayer()
		dataPlayer.RoleID = tbRoleBase.GetRoleID()
		dataPlayer.DataTbRoleBase = tbRoleBase
		dataGlobal.AddDataPlayer(dataPlayer.RoleID, dataPlayer)
		// ---  其他数据表初始化 begin ----
		dataGlobal.CheckInitPlayerTables(dataPlayer)
		// ---  其他数据表初始化 end ----
	}
	// 返回数据
	rep.RoleID = dataPlayer.RoleID
	rep.RoleDetailData = &pbserver.GameRoleData{}
	dataGlobal.FillMsgRoleDetailData(dataPlayer, rep.RoleDetailData)

	// ---  其他数据表返回 end ----
	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}

func HandleESMsgPlayerSaveRoleData(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbserver.ESMsgPlayerSaveRoleReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	gameDB := servData.GetGameDB()
	dbGlobal := servData.GetDataGlobal()
	dataPlayer := dbGlobal.FindDataPlayer(req.RoleID)
	if dataPlayer != nil {
		// 更新到缓存
		for _, v := range req.RoleTables {
			switch v.TableID {
			case ormdef.ETableRoleBase:
				{
					tbRoleBase := tbobj.NewTbRoleBase()
					tbRoleBase.FromBytes(v.Data)
					dataPlayer.DataTbRoleBase.FromBytes(v.Data)
					// 发送到db线程更新
					dbJob := func() {
						gameDB.Model(tbRoleBase.GetOrmMeta()).Select("*").Updates(tbRoleBase.GetOrmMeta())
					}
					dbGlobal.PostDBJob(&servdatamain.DataDBJob{
						DoJob: dbJob,
					})
					break
				}

				// -----playertable codetag6 begin-----------------
			case ormdef.ETableRoleEquip:
				{
					tbSaveObj := tbobj.NewTbRoleEquip()
					tbSaveObj.FromBytes(v.Data)
					dataPlayer.DataTbRoleEquip.FromBytes(v.Data)
					// 发送到db线程更新
					dbJob := func() {
						err := gameDB.Model(tbSaveObj.GetOrmMeta()).Select("*").Updates(tbSaveObj.GetOrmMeta()).Error
						if err != nil {
							loghlp.Errorf("save table error:%s", err.Error())
						}
					}
					dbGlobal.PostDBJob(&servdatamain.DataDBJob{
						DoJob: dbJob,
					})
					break
				}
			case ormdef.ETableRoleExtra:
				{
					tbSaveObj := tbobj.NewTbRoleExtra()
					tbSaveObj.FromBytes(v.Data)
					dataPlayer.DataTbRoleExtra.FromBytes(v.Data)
					// 发送到db线程更新
					dbJob := func() {
						err := gameDB.Model(tbSaveObj.GetOrmMeta()).Select("*").Updates(tbSaveObj.GetOrmMeta()).Error
						if err != nil {
							loghlp.Errorf("save table error:%s", err.Error())
						}
					}
					dbGlobal.PostDBJob(&servdatamain.DataDBJob{
						DoJob: dbJob,
					})
					break
				}
			case ormdef.ETableRoleCoin:
				{
					tbSaveObj := tbobj.NewTbRoleCoin()
					tbSaveObj.FromBytes(v.Data)
					dataPlayer.DataTbRoleCoin.FromBytes(v.Data)
					// 发送到db线程更新
					dbJob := func() {
						err := gameDB.Model(tbSaveObj.GetOrmMeta()).Select("*").Updates(tbSaveObj.GetOrmMeta()).Error
						if err != nil {
							loghlp.Errorf("save table error:%s", err.Error())
						}
					}
					dbGlobal.PostDBJob(&servdatamain.DataDBJob{
						DoJob: dbJob,
					})
					break
				}
			case ormdef.ETableRoleBag:
				{
					tbSaveObj := tbobj.NewTbRoleBag()
					tbSaveObj.FromBytes(v.Data)
					dataPlayer.DataTbRoleBag.FromBytes(v.Data)
					// 发送到db线程更新
					dbJob := func() {
						err := gameDB.Model(tbSaveObj.GetOrmMeta()).Select("*").Updates(tbSaveObj.GetOrmMeta()).Error
						if err != nil {
							loghlp.Errorf("save table error:%s", err.Error())
						}
					}
					dbGlobal.PostDBJob(&servdatamain.DataDBJob{
						DoJob: dbJob,
					})
					break
				}
				// -----playertable codetag6 end-------------------
			default:
				{
					loghlp.Errorf("known table id:%d", v.TableID)
				}
			}
		}
	}
	rep := &pbserver.ESMsgPlayerSaveRoleRep{}
	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}

func HandleECMsgPlayerKeepHeart(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	dbGlobal := servData.GetDataGlobal()
	dataPlayer := dbGlobal.FindDataPlayer(tmsgCtx.NetMessage.SecondHead.ID)
	if dataPlayer == nil {
		loghlp.Warnf("recv heart but not find dataplayer(%d)", tmsgCtx.NetMessage.SecondHead.ID)
	}
	rep := &pbclient.ECMsgPlayerKeepHeartRsp{}
	return protocol.ECodeSuccess, rep, iframe.EHandleNone
}
