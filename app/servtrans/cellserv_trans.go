/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-10-14 13:32:35
 * @LastEditTime: 2022-10-14 13:40:31
 * @FilePath: \trcell\app\servtrans\cellserv_trans.go
 */
package servtrans

import (
	"math/rand"
	"time"
	"trcell/app/servtrans/servtranshandler"
	"trcell/app/servtrans/servtransmain"
)

type CellServTrans struct {
	servTransGlobal *servtransmain.ServTransGlobal
}

func NewCellServTrans() *CellServTrans {
	rand.Seed(time.Now().UnixNano())
	s := &CellServTrans{
		servTransGlobal: servtransmain.NewServTransGlobal(),
	}
	s.RegisterMsgHandler()
	servtranshandler.InitServTransObj(s)
	return s
}

func (serv *CellServTrans) GetTransGlobal() *servtransmain.ServTransGlobal {
	return serv.servTransGlobal
}
func (serv *CellServTrans) FrameRun(curTimeMs int64) {
	serv.servTransGlobal.FrameUpdate(curTimeMs)
}
