/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 18:07:23
 * @LastEditTime: 2022-09-19 18:07:29
 * @FilePath: \trcell\app\servlog\iservlog\servlog_interface.go
 */
package iservlog

import "trcell/app/servlog/servlogmain"

type IServLog interface {
	GetLogGlobal() *servlogmain.ServLogGlobal
}
