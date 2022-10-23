/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 11:48:48
 * @LastEditTime: 2022-10-10 17:55:53
 * @FilePath: \trcell\app\servcenter\cellserv_center.go
 */
package servcenter

import (
	"math/rand"
	"time"
	"trcell/app/servcenter/servcenterhandler"
	"trcell/app/servcenter/servcentermain"
)

type CellServCenter struct {
	servCenterGlobal *servcentermain.ServCenterGlobal
}

func NewCellServCenter() *CellServCenter {
	rand.Seed(time.Now().UnixNano())
	s := &CellServCenter{
		servCenterGlobal: servcentermain.NewServCenterGlobal(),
	}
	s.RegisterMsgHandler()
	servcenterhandler.InitServCenterObj(s)
	return s
}

func (serv *CellServCenter) GetCenterGlobal() *servcentermain.ServCenterGlobal {
	return serv.servCenterGlobal
}
func (serv *CellServCenter) FrameRun(curTimeMs int64) {
	serv.servCenterGlobal.FrameUpdate(curTimeMs)
}
