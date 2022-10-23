/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-26 16:49:26
 * @LastEditTime: 2022-09-26 17:18:29
 * @FilePath: \trcell\app\servlogic\servlogicmain\playermd\pextra\player_extra.go
 */
package pextra

import (
	"trcell/app/servlogic/servlogicmain/playermd"
	"trcell/pkg/loghlp"
	"trcell/pkg/ormdef"
	"trcell/pkg/pb/pbserver"

	"google.golang.org/protobuf/proto"
)

type PlayerExtraData struct {
	player   playermd.ILogicPlayer
	scopeUid *PlayerScopeUid
}

func NewPlayerExtraData(p playermd.ILogicPlayer) *PlayerExtraData {
	exData := &PlayerExtraData{
		player:   p,
		scopeUid: newPlayerScopeUid(p.GetRoleID()),
	}
	return exData
}

func (mdObj *PlayerExtraData) ToBytes() []byte {
	dbExtra := &pbserver.DBRoleExtraData{
		ScopeUid: &pbserver.DBRoleScopeUid{
			UidPool: map[int64]int64{},
		},
	}
	dbExtra.ScopeUid.Idx = mdObj.scopeUid.idx
	for k, v := range mdObj.scopeUid.uidPool {
		dbExtra.ScopeUid.UidPool[k] = v
	}
	pbData, err := proto.Marshal(dbExtra)
	if err != nil {
		loghlp.Errorf("extra data marshal error:%s", err.Error())
	}
	return pbData
}

func (mdObj *PlayerExtraData) FromBytes(binaryData []byte) {
	dbExtra := &pbserver.DBRoleExtraData{}
	proto.Unmarshal(binaryData, dbExtra)
	if dbExtra.ScopeUid != nil {
		mdObj.scopeUid.idx = dbExtra.ScopeUid.Idx
		for k, v := range dbExtra.ScopeUid.UidPool {
			mdObj.scopeUid.uidPool[k] = v
		}
	}
}

func (mdObj *PlayerExtraData) UpdToDB() {
	mdObj.player.SetTableUpd(ormdef.ETableRoleExtra)
}

func (mdObj *PlayerExtraData) GenScopeUid() int64 {
	return mdObj.scopeUid.GenScopeUid()
}
