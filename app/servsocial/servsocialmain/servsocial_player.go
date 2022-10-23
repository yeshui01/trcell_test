/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-10-10 17:42:29
 * @LastEditTime: 2022-10-10 17:45:31
 * @FilePath: \trcell\app\servsocial\servsocialmain\servsocial_player.go
 */
package servsocialmain

type SocialPlayer struct {
	RoleID   int64
	Nickname string
	UserID   int64
	Level    int32
}

func NewSocialPlayer(roleID int64) *SocialPlayer {
	return &SocialPlayer{
		RoleID: roleID,
	}
}
