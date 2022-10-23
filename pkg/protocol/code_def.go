/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-16 14:18:14
 * @LastEditTime: 2022-09-23 14:55:16
 * @FilePath: \trcell\pkg\protocol\code_def.go
 */
package protocol

const (
	ECodeSuccess     = 0 // 正常
	ECodeSysError    = 1 // 系统错误
	ECodeParamError  = 2 // 参数错误
	ECodeAsyncHandle = 3 // 异步处理,当前不处理
	// 系统错误定义
	ECodeDBError            = 100 // db错误
	ECodePBDecodeError      = 101 // pb反序列化错误
	ECodeTokenError         = 102 // token解析失败
	ECodeTokenExpire        = 103 // token过期
	ECodeInvalideOperation  = 104 // 无效操作
	ECodeLogicServException = 105 // 逻辑服务器异常
	ECodeDataServException  = 106 // 数据服务器异常
	ECodeNotFindNotice      = 107 // 公告不存在
	// 错误码定义
	ECodeAccNameHasExisted       = 1000 // 账号已经存在
	ECodeAccNotExisted           = 1001 // 账号不存在
	ECodeAccCertificationTimeOut = 1002 // 账号超出最终认证时间
	ECodeAccPasswordError        = 1003 // 账号密码错误
	ECodeRoleNotExisted          = 1004 // 角色不存在
	ECodeRoleHasOnline           = 1005 // 玩家已经在线
	ECodeRoleCreatingLock        = 1006 // 创角锁定中
	ECodeRoleNickNameExisted     = 1007 // 角色昵称已经存在
)
