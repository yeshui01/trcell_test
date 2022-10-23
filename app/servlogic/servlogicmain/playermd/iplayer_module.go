/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-26 13:23:47
 * @LastEditTime: 2022-09-29 18:35:10
 * @FilePath: \trcell\app\servlogic\servlogicmain\playermd\iplayer_module.go
 */
package playermd

type ILogicPlayer interface {
	GetRoleID() int64
	SetTableUpd(tableID int32)
	GenScopeUid() int64
}
type IPlayerModule interface {
	ToBytes() []byte
	FromBytes(binaryData []byte)
	UpdToDB()
}

// 注意:这里必须是顺序连续
const (
	EPlayerModuleExtra = iota
	EPlayerModuleEquip
	EPlayerModuleCoin
	EPlayerModuleBag
	EPlayerModuleMax
)
