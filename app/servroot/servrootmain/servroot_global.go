/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 10:32:45
 * @LastEditTime: 2022-09-19 11:05:43
 * @FilePath: \trcell\app\servroot\servrootmain\servroot_global.go
 */
package servrootmain

import (
	"trcell/pkg/loghlp"
	"trcell/pkg/timeutil"
)

type ServRootGlobal struct {
	lastUpdateTime int64 // 上次更新的时间戳(秒)
	reportTime     int64
}

func NewServRootGlobal() *ServRootGlobal {
	return &ServRootGlobal{
		lastUpdateTime: 0,
		reportTime:     timeutil.NowTime(),
	}
}

func (servGlobal *ServRootGlobal) FrameUpdate(curTimeMs int64) {
	if curTimeMs/1000 > servGlobal.lastUpdateTime {
		servGlobal.lastUpdateTime = curTimeMs / 1000
		servGlobal.Update(servGlobal.lastUpdateTime)
	}
}
func (servGlobal *ServRootGlobal) Update(curTime int64) {
	loghlp.Debugf("rootServ update:%d", curTime)
}
