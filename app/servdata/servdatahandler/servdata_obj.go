/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 11:22:06
 * @LastEditTime: 2022-09-19 11:22:17
 * @FilePath: \trcell\app\servdata\servdatahandler\servdata_obj.go
 */
package servdatahandler

import (
	"trcell/app/servdata/iservdata"
)

var (
	servData iservdata.IServData
)

func InitServDataObj(iserv iservdata.IServData) {
	servData = iserv
}
