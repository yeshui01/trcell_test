/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-16 14:18:14
 * @LastEditTime: 2022-09-27 16:49:29
 * @FilePath: \trcell\pkg\protocol\protocol_client.go
 */
package protocol

// 客户端协议类别
const (
	ECMsgClassBase   = 1 // 通用基础类别
	ECMsgClassPlayer = 2 // 玩家类别
)

// ECMsgClassCommon = 1 // 通用类别
const (
	ECMsgBasePushErrorOpt = 1 // 错误操作推送
	//
	ECMsgBaseWrapOpt = 255 // 预留的特殊协议(客户端不用,主要预留给机器人测试)
)

// ECMsgClassPlayer = 2
const (
	ECMsgPlayerLoginGame  = 1 // 登录游戏
	ECMsgPlayerKeepHeart  = 2 // 心跳
	ECMsgPlayerCreateRole = 3 // 创角
)
