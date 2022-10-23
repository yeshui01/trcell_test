/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-10-10 17:51:57
 * @LastEditTime: 2022-10-10 17:52:24
 * @FilePath: \trcell\app\servsocial\servsocialhandler\servsocial_obj.go
 */
package servsocialhandler

import "trcell/app/servsocial/iservsocial"

var (
	servSocial iservsocial.IServSocial
)

func InitServSocialObj(iserv iservsocial.IServSocial) {
	servSocial = iserv
}
