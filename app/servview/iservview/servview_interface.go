/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 17:27:44
 * @LastEditTime: 2022-09-19 17:27:50
 * @FilePath: \trcell\app\servview\iservview\servview_interface.go
 */
package iservview

import "trcell/app/servview/servviewmain"

type IServView interface {
	GetViewGlobal() *servviewmain.ServViewGlobal
}
