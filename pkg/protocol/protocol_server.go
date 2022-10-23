/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-16 14:18:14
 * @LastEditTime: 2022-09-30 16:00:43
 * @FilePath: \trcell\pkg\protocol\protocol_server.go
 */
package protocol

// 服务端业务层面消息
const (
	ESMsgClassPlayer   int32 = 300 // 玩家类别
	ESMsgClassServData int32 = 301 // 服务器数据类别
)

// ESMsgClassPlayer int32 = 300 // 玩家类别
const (
	ESMsgPlayerLoadRole   = 1 // 加载角色数据
	ESMsgPlayerSaveRole   = 2 // 保存角色数据
	ESMsgPlayerLoginGame  = 3 // 登录游戏
	ESMsgPlayerDisconnect = 4 // 玩家连接断开
	ESMsgPlayerKickOut    = 5 // 踢掉玩家
	ESMsgPlayerCreateRole = 6 // 玩家创角
)

// ESMsgClassServData int32 = 301 // 服务器数据类别
const (
	ESMsgServDataLoadTables           = 1 // 加载数据表
	ESMsgServDataPushTablesPartial    = 2 // 推送数据表分片数据
	ESMsgServDataSaveTables           = 3 // 保存数据表
	ESMsgServDataLoadTableList        = 4 // 加载数据表(列表)
	ESMsgServDataSaveTableList        = 5 // 保存数据表(列表)
	ESMsgServDataPushTableListPartial = 6 // 推送数据列表 分片数据
)
