/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-28 18:44:01
 * @LastEditTime: 2022-09-29 14:08:41
 * @FilePath: \trcell\app\servlogic\servlogicmain\playermd\pbag\player_bags.go
 */
package pbag

import (
	"trcell/app/servlogic/servlogicmain/playermd"
	"trcell/pkg/loghlp"
	"trcell/pkg/ormdef"
	"trcell/pkg/pb/pbserver"

	"google.golang.org/protobuf/proto"
)

const (
	EBagTypeNone  = 0
	EBagTypeItem  = 1 // 普通道具背包
	EBagTypeEquip = 2 // 装备背包
	EBagTypeHero  = 3 // 英雄背包
)

type PlayerBag struct {
	player  playermd.ILogicPlayer
	BagList map[int32]*GoodsBag //key:bagType
}

func NewPlayerBag(p playermd.ILogicPlayer) *PlayerBag {
	return &PlayerBag{
		player:  p,
		BagList: make(map[int32]*GoodsBag),
	}
}

func (bagOwner *PlayerBag) ToBytes() []byte {
	dbPlayerBag := &pbserver.DBBagList{}
	for k, v := range bagOwner.BagList {
		dbBag := &pbserver.DBGoodsBag{
			BagType: k,
			PosIdx:  v.PosIdx,
			PosList: make(map[int32]*pbserver.DBBagPos),
		}
		for posID, bagPos := range v.PosList {
			dbBagPos := &pbserver.DBBagPos{
				PosID:     posID,
				GoodsType: bagPos.GoodsType,
				ExcelID:   bagPos.ExcelID,
				Num:       bagPos.Num,
				AttachID:  bagPos.AttachID,
			}
			dbBag.PosList[posID] = dbBagPos
		}
		dbPlayerBag.BagList = append(dbPlayerBag.BagList, dbBag)
	}
	saveData, _ := proto.Marshal(dbPlayerBag)
	return saveData
}
func (bagOwner *PlayerBag) FromBytes(binaryData []byte) {
	if binaryData == nil {
		return
	}
	dbPlayerBag := &pbserver.DBBagList{}
	err := proto.Unmarshal(binaryData, dbPlayerBag)
	if err != nil {
		loghlp.Errorf("PlayerBag FromBytes error:%s", err.Error())
		return
	}
	for _, dbBag := range dbPlayerBag.BagList {
		oneBag := NewGoodsBag(dbBag.BagType)
		oneBag.PosIdx = dbBag.PosIdx
		for posID, bagPos := range dbBag.PosList {
			oneBag.PosList[posID] = &BagPos{
				PosID:     posID,
				GoodsType: bagPos.GoodsType,
				Num:       bagPos.Num,
				AttachID:  bagPos.AttachID,
			}
		}
		bagOwner.BagList[oneBag.BagType] = oneBag
	}
}
func (bagOwner *PlayerBag) UpdToDB() {
	bagOwner.player.SetTableUpd(ormdef.ETableRoleBag)
}

func (bagOwner *PlayerBag) GetBag(bagType int32) *GoodsBag {
	if o, ok := bagOwner.BagList[bagType]; ok {
		return o
	}
	return nil
}
