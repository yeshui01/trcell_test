/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-28 11:24:41
 * @LastEditTime: 2022-10-17 10:40:11
 * @FilePath: \trcell\pkg\utils\id_generator.go
 */
package utils

import (
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const maxnodeID = 4095
const maxSeq = 4095

// 雪花算法生成64位唯一id nodeID 12位| seq 12位| 相对时间 40位
type IDGenerator struct {
	nodeID int64

	//seq 12bit 4096
	seq int64

	timeBase    int64
	lastGenTime int64
	mu          sync.Mutex // 多线程下,生成id时加锁保护一下
}

func NewIDGenerator(nodeID int64) (*IDGenerator, error) {

	if nodeID < 0 || nodeID > maxnodeID {
		return nil, fmt.Errorf("create id generator failed, nodeID=%d", nodeID)
	}
	g := &IDGenerator{
		nodeID:      nodeID,
		seq:         0,
		lastGenTime: 0,
		timeBase:    1663948800000, // 毫秒 -> 2022-09-24 00:00:00
	}
	return g, nil
}

func getCurMilliSecond() int64 {
	return time.Now().UnixNano() / 1e6
}

/**
 * @description: 基于雪花算法,线程安全,支持每毫秒4096个id生成
 * @return {*}
 */
func (g *IDGenerator) GetID() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	cur := getCurMilliSecond()
	if g.lastGenTime == 0 {
		g.lastGenTime = cur
	}
	if cur < g.lastGenTime {
		time.Sleep(time.Second)
		logrus.Errorf("system time error,generato id fail")
		return 0
		// panic("system time error")
	}
	//如果不是同一毫秒，重置seq
	if cur != g.lastGenTime {
		g.lastGenTime = cur
		g.seq = 0
	}
	if g.seq >= maxSeq {
		//如果seq用完了，等待一毫秒
		time.Sleep(time.Millisecond)
		g.lastGenTime = g.lastGenTime + 1
		g.seq = 0
	}
	//如果seq没有用完，直接返回id
	id := g.makeID()
	g.seq++
	return id
}

func (g *IDGenerator) makeID() int64 {
	return (g.lastGenTime-g.timeBase)<<24 | g.nodeID<<12 | g.seq
}

var idGenerator *IDGenerator

func InitIDGenerator(nodeID int64) error {
	var err error
	idGenerator, err = NewIDGenerator(nodeID)
	return err
}
func GetID() int64 {
	return idGenerator.GetID()
}
