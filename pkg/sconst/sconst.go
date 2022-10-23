/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-16 14:18:14
 * @LastEditTime: 2022-09-27 15:56:41
 * @FilePath: \trcell\pkg\sconst\sconst.go
 */
package sconst

const (
	AccountCertificationTime = 86400 * 14 // 一周
	// ClientMsgHeadSize           = 12
	ClientMsgHeadSize              = 6
	LogicPlayerIdleTime            = 1    // logic_player闲置保存时间
	LogicPlayerHeartCheckTime      = 60   // logic_player心跳检测时间
	TcpMTUSize                     = 1500 // MTU默认值
	TcpMessageInnerHeadAdapterSize = 64   // ip+tcp头部大小适配
	TcpPackUnit                    = 1436 // 拆包大小
)

// 断线原因
const (
	EPlayerOfflineReasonNormal       = 0 // 正常断线
	EPlayerOfflineReasonKickOut      = 1 // 踢人
	EPlayerOfflineReasonReplaceLogin = 2 // 顶号登录
)

// 玩家登录阶段
const (
	EPlayerLoginStepNone     = 0
	EPlayerLoginStepCreating = 1 // 创建中
	EPlayerLoginStepLogining = 2 // 登录
)
