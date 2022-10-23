/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-10-14 13:39:49
 * @LastEditTime: 2022-10-14 13:39:53
 * @FilePath: \trcell\app\servtrans\iservtrans\servtrans_interface.go
 */
package iservtrans

import "trcell/app/servtrans/servtransmain"

type IServTrans interface {
	GetTransGlobal() *servtransmain.ServTransGlobal
}
