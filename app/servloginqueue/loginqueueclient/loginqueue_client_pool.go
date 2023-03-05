/**
 * @author: [mknight]
 * @email : [824338670@qq.com]
 * @create:	2023-02-02 11:08:53
 * @modify:	2023-02-02 11:08:53
 * @desc  :	[description]
 */
package loginqueueclient

import (
	"sync"
	"trcell/pkg/loghlp"
)

type LoginQueueClientPool struct {
	poolSize    uint32
	servAddr    string
	clientMu    sync.Mutex
	clientList  []*ServLoginQueueClient
	condClients *sync.Cond
	numClient   int32
}

func NewLoginQueueClientPool(maxSize uint32, rpcAddr string) *LoginQueueClientPool {
	po := &LoginQueueClientPool{
		poolSize:   maxSize,
		servAddr:   rpcAddr,
		clientList: make([]*ServLoginQueueClient, 0),
		numClient:  0,
	}
	po.condClients = sync.NewCond(&(po.clientMu))
	return po
}

func (po *LoginQueueClientPool) GetClient() *ServLoginQueueClient {
	po.clientMu.Lock()
	if len(po.clientList) > 0 {
		c := po.clientList[0]
		po.clientList = po.clientList[1:len(po.clientList)]
		loghlp.Debugf("KeepReuse GetClient from clientList")
		po.clientMu.Unlock()
		return c
	}

	if po.numClient >= int32(po.poolSize) {
		// 说明满了,需要等待回收
		for len(po.clientList) < 1 {
			po.condClients.Wait()
		}
		c := po.clientList[0]
		po.clientList = po.clientList[1:len(po.clientList)]
		loghlp.Debugf("KeepReuse GetClient from clientList")
		po.clientMu.Unlock()
		return c
	}
	c := NewServLoginQueueClient(po.servAddr)
	errConn := c.Connect()
	if errConn != nil {
		loghlp.Errorf("new queue client connect error:%s", errConn.Error())
		po.clientMu.Unlock()
		return nil
	}

	loghlp.Debugf("KeepReuse GetClient New Client")
	po.numClient++
	po.clientMu.Unlock()
	return c
}

func (po *LoginQueueClientPool) KeepReuse(rpcClient *ServLoginQueueClient) {
	po.clientMu.Lock()
	if len(po.clientList) >= int(po.poolSize) {
		rpcClient.Close()
		loghlp.Infof("KeepReuse queue client pool reach max size, close queue client")
	} else {
		po.clientList = append(po.clientList, rpcClient)
		loghlp.Infof("KeepReuse reuse queue client")
	}
	po.condClients.Broadcast() // 通知唤醒
	po.clientMu.Unlock()
}
