/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-16 14:18:13
 * @LastEditTime: 2022-09-30 11:23:12
 * @FilePath: \trcell\pkg\ormdef\table_def.go
 */
package ormdef

const (
	ETableRoleBase = 1 // 玩家基础信息表
	// -----playertable codetag1 begin-----------------
	ETableRoleEquip = 2 // 玩家装备信息表
	ETableRoleExtra = 3 // 玩家额外数据表
	ETableRoleCoin  = 4 // 玩家货币数据表
	ETableRoleBag   = 5 // 玩家背包数据表
	// -----playertable codetag1 end-------------------

	// -----cstable codetag1 begin-----------------
	ETableCsGlobal = 10001
	// -----cstable codetag1 end-------------------
)
