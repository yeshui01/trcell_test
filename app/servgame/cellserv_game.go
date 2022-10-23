/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 15:14:05
 * @LastEditTime: 2022-09-19 15:22:25
 * @FilePath: \trcell\app\servgame\cellserv_game.go
 */
package servgame

import (
	"trcell/app/servgame/servgamehandler"
	"trcell/app/servgame/servgamemain"
)

type CellServGame struct {
	servGameGlobal *servgamemain.ServGameGlobal
}

func NewCellServGame() *CellServGame {
	s := &CellServGame{
		servGameGlobal: servgamemain.NewServGameGlobal(),
	}
	s.RegisterMsgHandler()
	servgamehandler.InitServGameObj(s)
	return s
}

func (serv *CellServGame) GetGameGlobal() *servgamemain.ServGameGlobal {
	return serv.servGameGlobal
}

func (serv *CellServGame) FrameRun(curTimeMs int64) {
	serv.servGameGlobal.FrameUpdate(curTimeMs)
}
