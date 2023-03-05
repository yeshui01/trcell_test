/**
 * @author: [mknight]
 * @email : [824338670@qq.com]
 * @create:	2023-01-18 10:16:02
 * @modify:	2023-01-18 10:16:02
 * @desc  :	[消息队列]
 */
package tmsgqueue

import (
	"container/list"
	"sync"
)

type TMsgQueue struct {
	msgQueue  list.List
	muQueue   sync.Mutex
	condQueue *sync.Cond
}

func NewTMsgQueue() *TMsgQueue {
	q := &TMsgQueue{}

	q.condQueue = sync.NewCond(&q.muQueue)
	return q
}

// 放入数据,这个接口在生产者线程调用
func (tq *TMsgQueue) Push(v any) {
	tq.muQueue.Lock()
	tq.msgQueue.PushBack(v)
	tq.muQueue.Unlock()
	tq.condQueue.Broadcast()
	// tq.condQueue.Signal()
}

// 获取一个元素,如果没有返回nil,这个接口在消费者线程调用
func (tq *TMsgQueue) Pop() any {
	tq.muQueue.Lock()
	defer tq.muQueue.Unlock()

	if tq.msgQueue.Len() < 1 {
		return nil
	}

	return tq.msgQueue.Remove(tq.msgQueue.Front())
}

// 等待获取一个元素,如果没有会阻塞等待,这个接口在消费者线程调用
func (tq *TMsgQueue) WaitPopOne() any {
	tq.muQueue.Lock()
	defer tq.muQueue.Unlock()

	for tq.msgQueue.Len() < 1 {
		tq.condQueue.Wait()
	}
	return tq.msgQueue.Remove(tq.msgQueue.Front())
}
