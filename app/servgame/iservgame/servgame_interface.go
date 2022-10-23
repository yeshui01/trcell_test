package iservgame

import "trcell/app/servgame/servgamemain"

type IServGame interface {
	GetGameGlobal() *servgamemain.ServGameGlobal
}
