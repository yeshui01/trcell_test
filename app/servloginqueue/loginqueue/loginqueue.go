package loginqueue

import (
	"sync"
	"trcell/pkg/timeutil"
)

type LoginQueue struct {
	loginSeq  int32
	finishNum int32
	loginMu   sync.Mutex

	lastFinishTime int64

	userSeqs map[string]int32 // 玩家序列号记录
	userMu   sync.RWMutex
}

func NewLoginQueue() *LoginQueue {
	return &LoginQueue{
		loginSeq:       0,
		finishNum:      0,
		userSeqs:       make(map[string]int32),
		lastFinishTime: 0,
	}
}

func (lq *LoginQueue) GenLoginSeqNo() (int32, int32) {
	lq.loginMu.Lock()
	defer lq.loginMu.Unlock()

	if lq.finishNum >= lq.loginSeq {
		// 重置一下
		lq.loginSeq = 0
		lq.finishNum = 0
	} else {
		nowTime := timeutil.NowTime()
		if nowTime-lq.lastFinishTime >= 90 {
			// 最近的一次是90秒前,那说明可能出现意外情况,中断了,这里开始重置
			for k := range lq.userSeqs {
				delete(lq.userSeqs, k)
			}
			lq.loginSeq = 0
			lq.finishNum = 0
			lq.lastFinishTime = nowTime
		}
	}

	lq.loginSeq++
	return lq.loginSeq, lq.loginSeq - lq.finishNum
}

func (lq *LoginQueue) GetCurLoginSeqNo() int32 {
	lq.loginMu.Lock()
	defer lq.loginMu.Unlock()
	return lq.loginSeq
}

func (lq *LoginQueue) GetFinishNum() int32 {
	lq.loginMu.Lock()
	defer lq.loginMu.Unlock()
	return lq.finishNum
}

func (lq *LoginQueue) AddFinishNum() {
	lq.loginMu.Lock()
	lq.finishNum++
	lq.lastFinishTime = timeutil.NowTime()
	lq.loginMu.Unlock()
}

func (lq *LoginQueue) GetUserSeq(userName string) int32 {
	lq.userMu.RLock()

	if n, ok := lq.userSeqs[userName]; ok {
		lq.userMu.RUnlock()
		return n
	}
	lq.userMu.RUnlock()

	return 0
}

func (lq *LoginQueue) CacheUserSeq(userName string, seqNo int32) bool {
	lq.userMu.Lock()
	if _, ok := lq.userSeqs[userName]; !ok {
		lq.userSeqs[userName] = seqNo
		lq.userMu.Unlock()
		return true
	}
	lq.userMu.Unlock()
	return false
}

func (lq *LoginQueue) DelUserSeq(userName string) {
	lq.userMu.Lock()
	delete(lq.userSeqs, userName)
	lq.userMu.Unlock()
}

func (lq *LoginQueue) GetQueueNum() int32 {
	lq.userMu.Lock()
	n := lq.loginSeq - lq.finishNum
	lq.userMu.Unlock()

	return n
}
