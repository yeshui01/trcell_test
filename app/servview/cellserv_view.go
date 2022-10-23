/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 17:20:21
 * @LastEditTime: 2022-09-19 17:29:55
 * @FilePath: \trcell\app\servview\cellserv_view.go
 */
package servview

import (
	"trcell/app/servview/servviewhandler"
	"trcell/app/servview/servviewmain"
)

type CellServView struct {
	servViewGlobal *servviewmain.ServViewGlobal
}

func NewCellServView() *CellServView {
	s := &CellServView{
		servViewGlobal: servviewmain.NewServViewGlobal(),
	}
	servviewhandler.InitServViewObj(s)
	s.RegisterMsgHandler()
	return s
}

func (serv *CellServView) GetViewGlobal() *servviewmain.ServViewGlobal {
	return serv.servViewGlobal
}

func (serv *CellServView) FrameRun(curTimeMs int64) {
	serv.servViewGlobal.FrameUpdate(curTimeMs)
}
