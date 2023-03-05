package iservloginqueue

import "trcell/app/servloginqueue/loginqueue"

type IServLoginQueue interface {
	GetLoginQueue() *loginqueue.LoginQueue
}
