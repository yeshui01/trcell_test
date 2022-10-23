/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 13:51:32
 * @LastEditTime: 2022-09-19 13:51:36
 * @FilePath: \trcell\app\servcenter\iservcenter\servcenter_interface.go
 */
package iservcenter

import "trcell/app/servcenter/servcentermain"

type IServCenter interface {
	GetCenterGlobal() *servcentermain.ServCenterGlobal
}
