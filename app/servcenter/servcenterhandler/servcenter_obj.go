/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 11:54:01
 * @LastEditTime: 2022-10-10 17:55:49
 * @FilePath: \trcell\app\servcenter\servcenterhandler\servcenter_obj.go
 */
package servcenterhandler

import (
	"trcell/app/servcenter/iservcenter"
)

var (
	servCenter iservcenter.IServCenter
)

func InitServCenterObj(iserv iservcenter.IServCenter) {
	servCenter = iserv
}
