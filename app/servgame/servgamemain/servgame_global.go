package servgamemain

import (
	"trcell/pkg/loghlp"
	"trcell/pkg/timeutil"
)

type ServGameGlobal struct {
	lastUpdateTime int64
}

func NewServGameGlobal() *ServGameGlobal {
	return &ServGameGlobal{
		lastUpdateTime: timeutil.NowTime(),
	}
}
func (serv *ServGameGlobal) FrameUpdate(curTimeMs int64) {
	if curTimeMs/1000 > serv.lastUpdateTime {
		serv.lastUpdateTime = curTimeMs / 1000
		serv.Update(serv.lastUpdateTime)
	}
}
func (serv *ServGameGlobal) Update(curTime int64) {
	loghlp.Debugf("game Update:%d", curTime)
}
