/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 18:02:26
 * @LastEditTime: 2022-09-19 18:17:21
 * @FilePath: \trcell\app\servlog\servlogmain\servlog_global.go
 */
package servlogmain

import "trcell/pkg/loghlp"

type ServLogGlobal struct {
	lastUpdateTime int64
}

func NewServLogGlobal() *ServLogGlobal {
	return &ServLogGlobal{}
}

func (servGlobal *ServLogGlobal) FrameUpdate(curTimeMs int64) {
	if curTimeMs/1000 > servGlobal.lastUpdateTime {
		servGlobal.lastUpdateTime = curTimeMs / 1000
		servGlobal.Update(servGlobal.lastUpdateTime)
	}
}

func (servGlobal *ServLogGlobal) Update(curTime int64) {
	loghlp.Debugf("logServ update:%d", curTime)
}
