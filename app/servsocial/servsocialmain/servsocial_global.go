/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-10-10 17:40:57
 * @LastEditTime: 2022-10-10 17:50:03
 * @FilePath: \trcell\app\servsocial\servsocialmain\servsocial_global.go
 */
package servsocialmain

type ServSocialGlobal struct {
	playerList     map[int64]*SocialPlayer
	lastUpdateTime int64
}

func NewServSocialGlobal() *ServSocialGlobal {
	return &ServSocialGlobal{
		playerList: make(map[int64]*SocialPlayer),
	}
}
func (servGlobal *ServSocialGlobal) FrameUpdate(curTimeMs int64) {
	if curTimeMs/1000 > servGlobal.lastUpdateTime {
		servGlobal.lastUpdateTime = curTimeMs / 1000
		servGlobal.Update(servGlobal.lastUpdateTime)
	}
}

func (servGlobal *ServSocialGlobal) Update(curTime int64) {
	//loghlp.Debugf("socialGlobal Update:%d", curTime)
}
