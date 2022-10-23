/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 18:06:42
 * @LastEditTime: 2022-09-19 18:08:10
 * @FilePath: \trcell\app\servlog\servloghandler\servlog_obj.go
 */
package servloghandler

import "trcell/app/servlog/iservlog"

var (
	servLog iservlog.IServLog
)

func InitServLogObj(iserv iservlog.IServLog) {
	servLog = iserv
}
