/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-26 13:49:57
 * @LastEditTime: 2022-09-29 18:35:01
 * @FilePath: \trcell\app\servlogic\servlogicmain\playermd\pequip\player_equip.go
 */
package pequip

import (
	"trcell/app/servlogic/servlogicmain/playermd"
	"trcell/pkg/loghlp"
	"trcell/pkg/ormdef"
	"trcell/pkg/pb/pbserver"
	"trcell/pkg/ustd"

	"google.golang.org/protobuf/proto"
)

type EquipData struct {
	EquipID int64
}

type LogicEquip struct {
	player    playermd.ILogicPlayer
	EquipList *ustd.Map[int64, *EquipData]
}

func NewLogicEquip(playerObj playermd.ILogicPlayer) *LogicEquip {
	return &LogicEquip{
		player:    playerObj,
		EquipList: ustd.NewMap[int64, *EquipData](),
	}
}

func (mdObj *LogicEquip) ToBytes() []byte {
	dbLogicEquip := &pbserver.DBPlayerEquip{
		EquipList: make(map[int64]*pbserver.DBEquipOne),
	}
	mdObj.EquipList.ForEach(func(k int64, v *EquipData) {
		oneEquip := &pbserver.DBEquipOne{
			EquipID: k,
		}
		dbLogicEquip.EquipList[k] = oneEquip
	})

	pbdata, err := proto.Marshal(dbLogicEquip)
	if err != nil {
		loghlp.Errorf("player(%d) logic equip ToBytes error:%s", mdObj.player.GetRoleID(), err.Error())
	}
	return pbdata
}
func (mdObj *LogicEquip) FromBytes(binaryData []byte) {
	dbLogicEquip := &pbserver.DBPlayerEquip{}
	proto.Unmarshal(binaryData, dbLogicEquip)
	for k, v := range dbLogicEquip.EquipList {
		oneEquip := &EquipData{
			EquipID: v.EquipID,
		}
		mdObj.EquipList.Set(k, oneEquip)
	}
	loghlp.Debugf("player(%d) equipmodule from bytes, len:%d", mdObj.player.GetRoleID(), mdObj.EquipList.Size())
}

func (mdObj *LogicEquip) UpdToDB() {
	mdObj.player.SetTableUpd(ormdef.ETableRoleEquip)
}

func (mdObj *LogicEquip) InitFirstEquip() {
	newEquip := &EquipData{
		EquipID: mdObj.player.GenScopeUid(),
	}
	mdObj.EquipList.Set(newEquip.EquipID, newEquip)
}
