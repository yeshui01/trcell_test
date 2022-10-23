/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 15:00:04
 * @LastEditTime: 2022-09-19 15:00:08
 * @FilePath: \trcell\app\servgate\servgatehandler\servgate_obj.go
 */
package servgatehandler

import "trcell/app/servgate/iservgate"

var (
	servGate iservgate.IServGate
)

func InitServGateObj(iserv iservgate.IServGate) {
	servGate = iserv
}
