package iservsocial

import "trcell/app/servsocial/servsocialmain"

type IServSocial interface {
	GetSocialGlobal() *servsocialmain.ServSocialGlobal
}
