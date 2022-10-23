/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 17:22:03
 * @LastEditTime: 2022-09-19 17:32:04
 * @FilePath: \trcell\app\servview\servviewmain\servview_global.go
 */
package servviewmain

import (
	"trcell/pkg/loghlp"
	"trcell/pkg/timeutil"
)

type ServViewGlobal struct {
	lastUpdateTime int64
	PlayerList     map[int64]*ViewPlayer
}

func NewServViewGlobal() *ServViewGlobal {
	return &ServViewGlobal{
		PlayerList:     make(map[int64]*ViewPlayer),
		lastUpdateTime: timeutil.NowTime(),
	}
}

func (servGlobal *ServViewGlobal) FrameUpdate(curTimeMs int64) {
	if curTimeMs/1000 > servGlobal.lastUpdateTime {
		servGlobal.lastUpdateTime = curTimeMs / 1000
		servGlobal.Update(servGlobal.lastUpdateTime)
	}
}

func (servGlobal *ServViewGlobal) Update(curTime int64) {
	loghlp.Debugf("servView update:%d", curTime)
}
