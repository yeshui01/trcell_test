/**
 * @author: [mknight]
 * @email : [824338670@qq.com]
 * @create:	2023-02-02 14:02:44
 * @modify:	2023-02-02 14:02:44
 * @desc  :	[接口限流]
 */
package accmiddleware

import (
	"log"
	"net/http"
	"sync/atomic"
	"time"
	"trcell/pkg/timeutil"

	"github.com/gin-gonic/gin"
)

// 接口限流器
//Limiter 限流器对象
type LoginLimiter struct {
	value int64
	max   int64
	ts    int64
}

//NewLoginLimiter 产生一个限流器
func NewLoginLimiter(cnt int64) *LoginLimiter {
	return &LoginLimiter{
		value: 0,
		max:   cnt,
		ts:    time.Now().Unix(),
	}
}

//Ok 是否可以通过
func (l *LoginLimiter) Ok() bool {
	ts := timeutil.NowTime()
	tsOld := atomic.LoadInt64(&l.ts)
	if ts != tsOld {
		atomic.StoreInt64(&l.ts, ts)
		atomic.StoreInt64(&l.value, 1)
		return true
	}
	return atomic.AddInt64(&(l.value), 1) < l.max
}

//SetMax 设置最大限制
func (l *LoginLimiter) SetMax(m int64) {
	l.max = m
}

//MaxAllowed 限流器
func MaxAllowed(limitValue int64) func(c *gin.Context) {
	limiter := NewLoginLimiter(limitValue)
	log.Println("loginlimiter.SetMax:", limitValue)
	// 返回限流逻辑
	return func(c *gin.Context) {
		if !limiter.Ok() {
			c.AbortWithStatus(http.StatusServiceUnavailable) //超过每秒1000，就返回503错误码
			return
		}
		c.Next()
	}
}
