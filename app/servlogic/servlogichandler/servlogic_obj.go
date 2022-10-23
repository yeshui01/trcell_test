/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 14:22:49
 * @LastEditTime: 2022-09-19 14:31:49
 * @FilePath: \trcell\app\servlogic\servlogichandler\servlogic_obj.go
 */
package servlogichandler

import (
	"trcell/app/servlogic/iservlogic"
)

var (
	servLogic iservlogic.IServLogic
)

func InitServLogicObj(iserv iservlogic.IServLogic) {
	servLogic = iserv
}
