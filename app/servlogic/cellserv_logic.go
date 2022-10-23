/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 14:15:18
 * @LastEditTime: 2022-09-23 15:39:34
 * @FilePath: \trcell\app\servlogic\cellserv_logic.go
 */
package servlogic

import (
	"trcell/app/servlogic/servlogichandler"
	"trcell/app/servlogic/servlogicmain"
)

type CellServLogic struct {
	servLogicGlobal *servlogicmain.ServLogicGlobal
	lastUpdateTime  int64
}

func NewCellServLogic() *CellServLogic {
	s := &CellServLogic{
		servLogicGlobal: servlogicmain.NewServLogicGlobal(),
	}
	servlogichandler.InitServLogicObj(s)
	s.RegisterMsgHandler()
	return s
}

func (serv *CellServLogic) GetLogicGlobal() *servlogicmain.ServLogicGlobal {
	return serv.servLogicGlobal
}

func (serv *CellServLogic) FrameRun(curTimeMs int64) {
	if curTimeMs/1000 > serv.lastUpdateTime {
		serv.lastUpdateTime = curTimeMs / 1000
		serv.servLogicGlobal.Update(serv.lastUpdateTime)
	}
}
