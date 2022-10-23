/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-26 16:37:18
 * @LastEditTime: 2022-09-26 17:29:06
 * @FilePath: \trcell\app\servlogic\servlogicmain\playermd\pextra\scope_uid.go
 */
package pextra

import (
	"trcell/pkg/timeutil"
	"trcell/pkg/trframe"
)

// 基于玩家数据范围内的uid
type PlayerScopeUid struct {
	roleID  int64
	idx     int64
	uidPool map[int64]int64
}

func newPlayerScopeUid(roleId int64) *PlayerScopeUid {
	return &PlayerScopeUid{
		roleID:  roleId,
		idx:     0,
		uidPool: make(map[int64]int64),
	}
}

func (obj *PlayerScopeUid) GenScopeUid() int64 {
	var rID int64 = 0
	curTime := trframe.GetFrameSysNowTime()
	for k, v := range obj.uidPool {
		if curTime-v >= 3 {
			rID = k
			break
		}
	}
	if rID == 0 {
		if obj.idx == 0 {
			obj.idx = 1000000
		}
		obj.idx++
		rID = obj.idx
	} else {
		delete(obj.uidPool, rID)
	}
	return rID
}
func (obj *PlayerScopeUid) RecycleScopeUid(uid int64) {
	obj.uidPool[uid] = timeutil.NowTime()
}
