/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 11:12:25
 * @LastEditTime: 2022-09-29 13:49:36
 * @FilePath: \trcell\app\servdata\servdatamain\servdata_player.go
 */
package servdatamain

import "trcell/pkg/tbobj"

type DataPlayer struct {
	RoleID int64
	// -----playertable codetag8 begin-----------------
	DataTbRoleBase  *tbobj.TbRoleBase
	DataTbRoleEquip *tbobj.TbRoleEquip
	DataTbRoleExtra *tbobj.TbRoleExtra
	DataTbRoleCoin  *tbobj.TbRoleCoin
	DataTbRoleBag   *tbobj.TbRoleBag
	// -----playertable codetag8 end-------------------
	VisitTime int64
}

func NewDataPlayer() *DataPlayer {
	return &DataPlayer{}
}
