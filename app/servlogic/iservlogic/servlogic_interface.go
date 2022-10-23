/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 14:29:30
 * @LastEditTime: 2022-09-19 14:30:11
 * @FilePath: \trcell\app\servlogic\iservlogic\servlogic_interface.go
 */
package iservlogic

import "trcell/app/servlogic/servlogicmain"

type IServLogic interface {
	GetLogicGlobal() *servlogicmain.ServLogicGlobal
}
