/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 18:01:03
 * @LastEditTime: 2022-09-19 18:10:43
 * @FilePath: \trcell\app\servlog\cellserv_log.go
 */
package servlog

import (
	"trcell/app/servlog/servloghandler"
	"trcell/app/servlog/servlogmain"
)

type CellServLog struct {
	servLogGlobal *servlogmain.ServLogGlobal
}

func NewCellServLog() *CellServLog {
	s := &CellServLog{
		servLogGlobal: servlogmain.NewServLogGlobal(),
	}
	s.RegisterMsgHandler()
	servloghandler.InitServLogObj(s)
	return s
}

func (serv *CellServLog) GetLogGlobal() *servlogmain.ServLogGlobal {
	return serv.servLogGlobal
}

func (serv *CellServLog) FrameRun(curTimeMs int64) {
	serv.servLogGlobal.FrameUpdate(curTimeMs)
}
