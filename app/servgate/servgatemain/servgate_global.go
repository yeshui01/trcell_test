/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 14:54:40
 * @LastEditTime: 2022-09-19 14:57:48
 * @FilePath: \trcell\app\servgate\servgatemain\servgate_global.go
 */
package servgatemain

import "trcell/pkg/loghlp"

type ServGateGlobal struct {
	lastUpdateTime int64
}

func NewServGateGlobal() *ServGateGlobal {
	return &ServGateGlobal{}
}

func (servGlobal *ServGateGlobal) FrameUpdate(curTimeMs int64) {
	if curTimeMs/1000 > servGlobal.lastUpdateTime {
		servGlobal.lastUpdateTime = curTimeMs / 1000
		servGlobal.Update(servGlobal.lastUpdateTime)
	}
}

func (servGlobal *ServGateGlobal) Update(curTime int64) {
	loghlp.Debugf("servGate update:%d", curTime)
}
