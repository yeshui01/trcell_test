package servviewhandler

import (
	"trcell/app/servview/iservview"
)

var (
	servView iservview.IServView
)

func InitServViewObj(iserv iservview.IServView) {
	servView = iserv
}
