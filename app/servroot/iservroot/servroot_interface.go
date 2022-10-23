/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 10:35:11
 * @LastEditTime: 2022-09-19 10:40:08
 * @FilePath: \trcell\app\servroot\iservroot\servroot_interface.go
 */
package iservroot

import "trcell/app/servroot/servrootmain"

type IServRoot interface {
	GetRootGlobal() *servrootmain.ServRootGlobal
}
