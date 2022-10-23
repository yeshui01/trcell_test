/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 11:51:20
 * @LastEditTime: 2022-09-23 15:18:42
 * @FilePath: \trcell\app\servcenter\servcentermain\servcenter_global.go
 */
package servcentermain

import (
	"trcell/pkg/timeutil"
	"trcell/pkg/trframe"
)

type ServCenterGlobal struct {
	playerList      map[int64]*CenterPlayer
	nameToPlayer    map[string]*CenterPlayer
	lastUpdateTime  int64
	lockNames       map[string]int64 // 临时锁定的名字
	pendingCreating map[int64]int64  // key:user_id value:create_time
}

func NewServCenterGlobal() *ServCenterGlobal {
	return &ServCenterGlobal{
		playerList:      make(map[int64]*CenterPlayer),
		nameToPlayer:    make(map[string]*CenterPlayer),
		pendingCreating: make(map[int64]int64),
		lockNames:       make(map[string]int64),
		lastUpdateTime:  timeutil.NowTime(),
	}
}

func (servGlobal *ServCenterGlobal) FrameUpdate(curTimeMs int64) {
	if curTimeMs/1000 > servGlobal.lastUpdateTime {
		servGlobal.lastUpdateTime = curTimeMs / 1000
		servGlobal.Update(servGlobal.lastUpdateTime)
	}
}

func (servGlobal *ServCenterGlobal) Update(curTime int64) {
	//loghlp.Debugf("centerGlobal Update:%d", curTime)
}

/**
 * @description: 根据id查找cs player
 * @param {int64} roleID
 * @return {*}
 */
func (servGlobal *ServCenterGlobal) FindCsPlayer(roleID int64) *CenterPlayer {
	if p, ok := servGlobal.playerList[roleID]; ok {
		return p
	}
	return nil
}

/**
 * @description: 根据昵称查找cs player
 * @param {string} nickName
 * @return {*}
 */
func (servGlobal *ServCenterGlobal) FindCsPlayerByName(nickName string) *CenterPlayer {
	if p, ok := servGlobal.nameToPlayer[nickName]; ok {
		return p
	}
	return nil
}

/**
 * @description: 创建一个角色
 * @param {int64} roleID
 * @param {string} nickName
 * @param {int32} level
 * @param {int32} icon
 * @return {*}
 */
func (servGlobal *ServCenterGlobal) CreateCsPlayer(roleID int64,
	nickName string,
	level int32,
	icon int32) *CenterPlayer {
	csPlayer := NewCenterPlayer(roleID, nickName, level)
	servGlobal.playerList[roleID] = csPlayer
	servGlobal.nameToPlayer[nickName] = csPlayer
	return csPlayer
}

func (servGlobal *ServCenterGlobal) LockNickname(nickName string) {
	servGlobal.lockNames[nickName] = trframe.GetFrameSysNowTime()
}

func (servGlobal *ServCenterGlobal) IsNicknameLock(nickName string) bool {
	lockTime, ok := servGlobal.lockNames[nickName]
	if ok {
		if timeutil.NowTime()-lockTime > 120 {
			delete(servGlobal.lockNames, nickName)
			return false
		}
	}
	return ok
}
func (servGlobal *ServCenterGlobal) UnlockNickname(nickName string) {
	delete(servGlobal.lockNames, nickName)
}

func (servGlobal *ServCenterGlobal) LockUserCreateRole(userID int64) {
	servGlobal.pendingCreating[userID] = trframe.GetFrameSysNowTime()
}

func (servGlobal *ServCenterGlobal) IsUserCreateLock(userID int64) bool {
	lockTime, ok := servGlobal.pendingCreating[userID]
	if ok {
		if trframe.GetFrameSysNowTime()-lockTime >= 120 {
			delete(servGlobal.pendingCreating, userID)
			return false
		}
	}
	return ok
}

func (servGlobal *ServCenterGlobal) UnlockUserCreate(userID int64) {
	delete(servGlobal.pendingCreating, userID)
}
