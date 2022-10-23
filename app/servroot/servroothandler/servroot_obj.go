/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 10:36:24
 * @LastEditTime: 2022-09-19 10:37:38
 * @FilePath: \trcell\app\servroot\servroothandler\servroot_obj.go
 */
package servroothandler

import "trcell/app/servroot/iservroot"

var (
	servRoot iservroot.IServRoot
)

func InitServRootObj(iserv iservroot.IServRoot) {
	servRoot = iserv
}
