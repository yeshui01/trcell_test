/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-10-14 13:35:22
 * @LastEditTime: 2022-10-14 13:35:27
 * @FilePath: \trcell\app\servtrans\servtransmain\servtrans_global.go
 */
package servtransmain

import "trcell/pkg/loghlp"

type ServTransGlobal struct {
	lastUpdateTime int64
}

func NewServTransGlobal() *ServTransGlobal {
	return &ServTransGlobal{
		lastUpdateTime: 0,
	}
}
func (servGlobal *ServTransGlobal) FrameUpdate(curTimeMs int64) {
	if curTimeMs/1000 > servGlobal.lastUpdateTime {
		servGlobal.lastUpdateTime = curTimeMs / 1000
		servGlobal.Update(servGlobal.lastUpdateTime)
	}
}

func (servGlobal *ServTransGlobal) Update(curTime int64) {
	loghlp.Debugf("transGlobal Update:%d", curTime)
}
