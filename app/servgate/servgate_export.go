/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-20 11:06:48
 * @LastEditTime: 2022-09-20 11:07:07
 * @FilePath: \trcell\app\servgate\servgate_export.go
 */
package servgate

import "trcell/app/servgate/servgatemain"

func (hg *CellServGate) GetUserManager() *servgatemain.HGateUserManager {
	return hg.UserMgr
}
