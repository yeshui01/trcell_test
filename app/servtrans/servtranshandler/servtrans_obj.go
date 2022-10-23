/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-10-14 13:38:29
 * @LastEditTime: 2022-10-14 13:40:48
 * @FilePath: \trcell\app\servtrans\servtranshandler\servtrans_obj.go
 */
package servtranshandler

import "trcell/app/servtrans/iservtrans"

var (
	servTrans iservtrans.IServTrans
)

func InitServTransObj(iserv iservtrans.IServTrans) {
	servTrans = iserv
}
