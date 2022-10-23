/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-30 11:10:32
 * @LastEditTime: 2022-09-30 11:11:25
 * @FilePath: \trcell\app\servdata\servdatamain\servdata_global_data_tables.go
 */
package servdatamain

import (
	"trcell/pkg/loghlp"
	"trcell/pkg/ormdef"
	"trcell/pkg/pb/pbserver"
	"trcell/pkg/tbobj"
)

func (servGlobal *ServDataGlobal) CheckInitPlayerTables(dataPlayer *DataPlayer) {
	gameDB := servGlobal.gameDB
	var err error
	// -----playertable codetag9 begin-----------------
	{
		// role_equip
		dataPlayer.DataTbRoleEquip = tbobj.NewTbRoleEquip()
		dataPlayer.DataTbRoleEquip.SetRoleID(dataPlayer.RoleID)
		ormMeta := dataPlayer.DataTbRoleEquip.GetOrmMeta()
		err = gameDB.Model(ormMeta).First(ormMeta).Error
		if err != nil {
			loghlp.Warnf("check init player(%d) exception:%s", dataPlayer.RoleID, err.Error())
			// 创建新的
			err = gameDB.Model(ormMeta).Create(ormMeta).Error
			if err == nil {
				loghlp.Infof("init player(%d) table[TbRoleEquip] succ", dataPlayer.RoleID)
			} else {
				loghlp.Errorf("init player(%d) table[TbRoleEquip] error:%s", dataPlayer.RoleID, err.Error())
			}
		}
	}
	{
		// role_extra
		dataPlayer.DataTbRoleExtra = tbobj.NewTbRoleExtra()
		dataPlayer.DataTbRoleExtra.SetRoleID(dataPlayer.RoleID)
		ormMeta := dataPlayer.DataTbRoleExtra.GetOrmMeta()
		err = gameDB.Model(ormMeta).First(ormMeta).Error
		if err != nil {
			loghlp.Warnf("check init player(%d) exception:%s", dataPlayer.RoleID, err.Error())
			// 创建新的
			err = gameDB.Model(ormMeta).Create(ormMeta).Error
			if err == nil {
				loghlp.Infof("init player(%d) table[TbRoleExtra] succ", dataPlayer.RoleID)
			} else {
				loghlp.Errorf("init player(%d) table[TbRoleExtra] error:%s", dataPlayer.RoleID, err.Error())
			}
		}
	}
	{
		// role_coin
		dataPlayer.DataTbRoleCoin = tbobj.NewTbRoleCoin()
		dataPlayer.DataTbRoleCoin.SetRoleID(dataPlayer.RoleID)
		ormMeta := dataPlayer.DataTbRoleCoin.GetOrmMeta()
		err = gameDB.Model(ormMeta).First(ormMeta).Error
		if err != nil {
			loghlp.Warnf("check init player(%d) exception:%s", dataPlayer.RoleID, err.Error())
			// 创建新的
			err = gameDB.Model(ormMeta).Create(ormMeta).Error
			if err == nil {
				loghlp.Infof("init player(%d) table[TbRoleCoin] succ", dataPlayer.RoleID)
			} else {
				loghlp.Errorf("init player(%d) table[TbRoleCoin] error:%s", dataPlayer.RoleID, err.Error())
			}
		}
	}
	{
		// role_bag
		dataPlayer.DataTbRoleBag = tbobj.NewTbRoleBag()
		dataPlayer.DataTbRoleBag.SetRoleID(dataPlayer.RoleID)
		ormMeta := dataPlayer.DataTbRoleBag.GetOrmMeta()
		err = gameDB.Model(ormMeta).First(ormMeta).Error
		if err != nil {
			loghlp.Warnf("check init player(%d) exception:%s", dataPlayer.RoleID, err.Error())
			// 创建新的
			err = gameDB.Model(ormMeta).Create(ormMeta).Error
			if err == nil {
				loghlp.Infof("init player(%d) table[TbRoleBag] succ", dataPlayer.RoleID)
			} else {
				loghlp.Errorf("init player(%d) table[TbRoleBag] error:%s", dataPlayer.RoleID, err.Error())
			}
		}
	}
	// -----playertable codetag9 end-------------------
}
func (servGlobal *ServDataGlobal) FillMsgRoleDetailData(dataPlayer *DataPlayer, pbMsgRoleDetail *pbserver.GameRoleData) {
	// role_base
	{
		oneTable := &pbserver.DbTableData{
			TableID: ormdef.ETableRoleBase,
		}
		oneTable.Data, _ = dataPlayer.DataTbRoleBase.ToBytes()
		pbMsgRoleDetail.RoleTables = append(pbMsgRoleDetail.RoleTables, oneTable)
	}
	// 其他表数据TODO

	// -----playertable codetag5 begin-----------------
	{
		// role_equip
		oneTable := &pbserver.DbTableData{
			TableID: ormdef.ETableRoleEquip,
		}
		oneTable.Data, _ = dataPlayer.DataTbRoleEquip.ToBytes()
		pbMsgRoleDetail.RoleTables = append(pbMsgRoleDetail.RoleTables, oneTable)
	}
	{
		// role_extra
		oneTable := &pbserver.DbTableData{
			TableID: ormdef.ETableRoleExtra,
		}
		oneTable.Data, _ = dataPlayer.DataTbRoleExtra.ToBytes()
		pbMsgRoleDetail.RoleTables = append(pbMsgRoleDetail.RoleTables, oneTable)
	}
	{
		// role_coin
		oneTable := &pbserver.DbTableData{
			TableID: ormdef.ETableRoleCoin,
		}
		oneTable.Data, _ = dataPlayer.DataTbRoleCoin.ToBytes()
		pbMsgRoleDetail.RoleTables = append(pbMsgRoleDetail.RoleTables, oneTable)
	}
	{
		// role_bag
		oneTable := &pbserver.DbTableData{
			TableID: ormdef.ETableRoleBag,
		}
		oneTable.Data, _ = dataPlayer.DataTbRoleBag.ToBytes()
		pbMsgRoleDetail.RoleTables = append(pbMsgRoleDetail.RoleTables, oneTable)
	}
	// -----playertable codetag5 end-------------------
}
