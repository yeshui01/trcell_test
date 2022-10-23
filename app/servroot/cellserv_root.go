/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-16 14:18:13
 * @LastEditTime: 2022-09-19 10:42:12
 * @FilePath: \trcell\app\servroot\cellserv_root.go
 */
package servroot

import (
	"trcell/app/servroot/servroothandler"
	"trcell/app/servroot/servrootmain"
)

type CellServRoot struct {
	servRootGlobal *servrootmain.ServRootGlobal
}

func (serv *CellServRoot) GetRootGlobal() *servrootmain.ServRootGlobal {
	return serv.servRootGlobal
}

func NewCellServRoot() *CellServRoot {
	s := &CellServRoot{
		servRootGlobal: servrootmain.NewServRootGlobal(),
	}
	servroothandler.InitServRootObj(s)
	return s
}

func (serv *CellServRoot) FrameRun(curTimeMs int64) {
	serv.servRootGlobal.FrameUpdate(curTimeMs)
}
