/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 14:16:14
 * @LastEditTime: 2022-10-12 18:10:38
 * @FilePath: \trcell\app\servlogic\servlogicmain\servlogic_global.go
 */
package servlogicmain

import (
	"trcell/pkg/loghlp"
	"trcell/pkg/ormdef"
	"trcell/pkg/pb/pbserver"
	"trcell/pkg/protocol"
	"trcell/pkg/sconst"
	"trcell/pkg/tbobj"
	"trcell/pkg/trframe"
	"trcell/pkg/trframe/iframe"
	"trcell/pkg/trframe/trnode"
)

type ServLogicGlobal struct {
	playerList     map[int64]*LogicPlayer
	user2Player    map[int64]*LogicPlayer
	lastUpdateTime int64
	onceJobList    []func() // 单次job列表
	reportTime     int64
}

func NewServLogicGlobal() *ServLogicGlobal {
	return &ServLogicGlobal{
		playerList:  map[int64]*LogicPlayer{},
		user2Player: map[int64]*LogicPlayer{},
	}
}

func (servGlobal *ServLogicGlobal) FrameUpdate(curTimeMs int64) {
	if curTimeMs/1000 > servGlobal.lastUpdateTime {
		servGlobal.lastUpdateTime = curTimeMs / 1000
		servGlobal.Update(servGlobal.lastUpdateTime)
	}
	servGlobal.updateJob(curTimeMs)
}
func (servGlobal *ServLogicGlobal) Update(curTime int64) {
	//loghlp.Debugf("logicGlobal update:%d", curTime)
	if curTime <= servGlobal.lastUpdateTime {
		return
	}
	servGlobal.UpdateSavePlayer(curTime)
	for _, v := range servGlobal.playerList {
		v.SecUpdate(curTime)
	}
	servGlobal.UpdateRemoveIdlePlayer(curTime)
	if curTime-servGlobal.reportTime >= 60 {
		servGlobal.ReportStatus()
		servGlobal.reportTime = curTime
	}
	servGlobal.lastUpdateTime = curTime
}
func (servGlobal *ServLogicGlobal) UpdateSavePlayer(curTime int64) {
	var saveCount int = 0
	for _, player := range servGlobal.playerList {
		if curTime-player.LastSaveTime >= 3 {
			// 先固定30秒保存一次
			servGlobal.SavePlayer(player, false)
			saveCount++
			if saveCount >= 200 {
				// 每次最多保存200个
				break
			}
		}
	}
}
func (servGlobal *ServLogicGlobal) SavePlayer(player *LogicPlayer, forceAll bool) {
	reqDB := &pbserver.ESMsgPlayerSaveRoleReq{
		RoleID: player.RoleID,
	}
	var saveTableNum int32 = 0
	curTime := trframe.GetFrameSysNowTime()
	if player.tbRoleBase.GetDbStatus() != tbobj.DbStatusNone {
		roleBase := &pbserver.DbTableData{
			TableID: ormdef.ETableRoleBase,
		}
		player.tbRoleBase.SetUpdTime(curTime)
		roleBase.Data, _ = player.tbRoleBase.ToBytes()
		reqDB.RoleTables = append(reqDB.RoleTables, roleBase)
		saveTableNum++
		loghlp.Debugf("saveplayer(%d) module[ETableRoleBase]", player.RoleID)
		player.tbRoleBase.ClearDbStatus()
	}
	// -----playertable codetag2 begin-----------------
	if forceAll || player.tbRoleEquip.GetDbStatus() != tbobj.DbStatusNone {
		roleTableObj := &pbserver.DbTableData{
			TableID: ormdef.ETableRoleEquip,
		}
		player.tbRoleEquip.SetUpdTime(curTime)
		roleTableObj.Data, _ = player.tbRoleEquip.ToBytes()
		reqDB.RoleTables = append(reqDB.RoleTables, roleTableObj)
		saveTableNum++
		loghlp.Debugf("saveplayer(%d) module[ETableRoleEquip]", player.RoleID)
		player.tbRoleEquip.ClearDbStatus()
	}
	if forceAll || player.tbRoleExtra.GetDbStatus() != tbobj.DbStatusNone {
		roleTableObj := &pbserver.DbTableData{
			TableID: ormdef.ETableRoleExtra,
		}
		player.tbRoleExtra.SetUpdTime(curTime)
		roleTableObj.Data, _ = player.tbRoleExtra.ToBytes()
		reqDB.RoleTables = append(reqDB.RoleTables, roleTableObj)
		saveTableNum++
		loghlp.Debugf("saveplayer(%d) module[ETableRoleExtra]", player.RoleID)
		player.tbRoleExtra.ClearDbStatus()
	}
	if forceAll || player.tbRoleCoin.GetDbStatus() != tbobj.DbStatusNone {
		roleTableObj := &pbserver.DbTableData{
			TableID: ormdef.ETableRoleCoin,
		}
		player.tbRoleCoin.SetUpdTime(curTime)
		roleTableObj.Data, _ = player.tbRoleCoin.ToBytes()
		reqDB.RoleTables = append(reqDB.RoleTables, roleTableObj)
		saveTableNum++
		loghlp.Debugf("saveplayer(%d) module[ETableRoleCoin]", player.RoleID)
		player.tbRoleCoin.ClearDbStatus()
	}
	if forceAll || player.tbRoleBag.GetDbStatus() != tbobj.DbStatusNone {
		roleTableObj := &pbserver.DbTableData{
			TableID: ormdef.ETableRoleBag,
		}
		player.tbRoleBag.SetUpdTime(curTime)
		roleTableObj.Data, _ = player.tbRoleBag.ToBytes()
		reqDB.RoleTables = append(reqDB.RoleTables, roleTableObj)
		saveTableNum++
		loghlp.Debugf("saveplayer(%d) module[ETableRoleBag]", player.RoleID)
		player.tbRoleBag.ClearDbStatus()
	}
	// -----playertable codetag2 end-------------------
	// 发送到db保存
	// 发送消息

	// 这里的session是frameSession
	if saveTableNum > 0 {
		roleID := player.RoleID
		cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
			loghlp.Infof("save player(%d) callback success,okCode:%d", roleID, okCode)
		}
		var dataIndex int32 = player.GetNetPeerIndex(trnode.ETRNodeTypeCellData)
		cbEnv := trframe.MakeMsgEnv(0, nil)
		trframe.ForwardZoneMessage(
			protocol.ESMsgClassPlayer,
			protocol.ESMsgPlayerSaveRole,
			reqDB,
			trnode.ETRNodeTypeCellData,
			dataIndex,
			cb,
			cbEnv,
		)
		player.LastRealSaveTime = trframe.GetFrameSysNowTime()
	}
	player.LastSaveTime = trframe.GetFrameSysNowTime()
}

func (servGlobal *ServLogicGlobal) FindPlayer(roleID int64) *LogicPlayer {
	if p, ok := servGlobal.playerList[roleID]; ok {
		return p
	}
	return nil
}
func (servGlobal *ServLogicGlobal) FindPlayerByUserID(userID int64) *LogicPlayer {
	if p, ok := servGlobal.user2Player[userID]; ok {
		return p
	}
	return nil
}
func (servGlobal *ServLogicGlobal) AddPlayer(p *LogicPlayer) {
	servGlobal.playerList[p.RoleID] = p
	servGlobal.user2Player[p.tbRoleBase.GetUserID()] = p
}

func (servGlobal *ServLogicGlobal) HandlePlayerOnline(player *LogicPlayer) {
	loghlp.Infof("HandlePlayerOnline:%d", player.RoleID)
	player.Online()
}
func (servGlobal *ServLogicGlobal) HandlePlayerOffline(player *LogicPlayer) {
	loghlp.Infof("HandlePlayerOffline:%d", player.RoleID)
	player.Offline()
	servGlobal.SavePlayer(player, true)
}
func (servGlobal *ServLogicGlobal) UpdateRemoveIdlePlayer(curTime int64) {
	var idlePlayers []int64
	for k, v := range servGlobal.playerList {
		if !v.IsOnline && curTime-v.OfflineTime >= sconst.LogicPlayerIdleTime {
			loghlp.Debugf("will remove idleplayer(%d) for offline timeout, offlineTime:%d", k, v.OfflineTime)
			idlePlayers = append(idlePlayers, k)
			delete(servGlobal.user2Player, v.tbRoleBase.GetUserID())
			if len(idlePlayers) >= 200 {
				// 一次移除200个
				break
			}
		} else if v.IsOnline {
			if v.GetHeartTime() == 0 {
				v.UpdateHeartTime(curTime)
			} else if curTime-v.GetHeartTime() >= sconst.LogicPlayerHeartCheckTime {
				loghlp.Debugf("will remove idleplayer(%d) for hearttime timeout, hearttime:%d", k, v.GetHeartTime())
				idlePlayers = append(idlePlayers, k)
				delete(servGlobal.user2Player, v.tbRoleBase.GetUserID())
				if len(idlePlayers) >= 200 {
					// 一次移除200个
					break
				}
			}
		}
	}
	for _, v := range idlePlayers {
		loghlp.Infof("remove idle player(%d)", v)

		delete(servGlobal.playerList, v)
	}
}
func (servGlobal *ServLogicGlobal) PostJob(doJob func()) {
	servGlobal.onceJobList = append(servGlobal.onceJobList, doJob)
}

func (servGlobal *ServLogicGlobal) updateJob(curTimeMs int64) {
	if len(servGlobal.onceJobList) > 0 {
		for _, doJob := range servGlobal.onceJobList {
			doJob()
		}
		servGlobal.onceJobList = nil
	}
}

func (servGlobal *ServLogicGlobal) ReportStatus() {
	loghlp.Infof("ReportStatus,playerNum:%d,onceJobList:%d", len(servGlobal.playerList), len(servGlobal.onceJobList))
}
