/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-10-10 17:39:26
 * @LastEditTime: 2022-10-11 14:58:06
 * @FilePath: \trcell\app\servsocial\cellserv_social.go
 */
package servsocial

import (
	"math/rand"
	"time"
	"trcell/app/servsocial/servsocialhandler"
	"trcell/app/servsocial/servsocialmain"
)

type CellServSocial struct {
	servSocialGlobal *servsocialmain.ServSocialGlobal
}

func NewCellServSocial() *CellServSocial {
	rand.Seed(time.Now().UnixNano())
	s := &CellServSocial{
		servSocialGlobal: servsocialmain.NewServSocialGlobal(),
	}
	s.RegisterMsgHandler()
	servsocialhandler.InitServSocialObj(s)
	return s
}

func (serv *CellServSocial) GetSocialGlobal() *servsocialmain.ServSocialGlobal {
	return serv.servSocialGlobal
}
func (serv *CellServSocial) FrameRun(curTimeMs int64) {
	serv.servSocialGlobal.FrameUpdate(curTimeMs)
}
