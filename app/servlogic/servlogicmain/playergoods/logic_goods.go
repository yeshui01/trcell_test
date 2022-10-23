/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-28 18:28:51
 * @LastEditTime: 2022-09-29 15:34:40
 * @FilePath: \trcell\app\servlogic\servlogicmain\playergoods\logic_goods.go
 */
package playergoods

import (
	"trcell/app/servlogic/servlogicmain/playermd"
	"trcell/app/servlogic/servlogicmain/playermd/pbag"
	"trcell/app/servlogic/servlogicmain/playermd/pcoin"
	"trcell/pkg/gamemd/goods"
	"trcell/pkg/loghlp"
)

type IGoodsPlayer interface {
	GetModule(moduleID int32) (playermd.IPlayerModule, error)
	GetRoleID() int64
	GetModuleCoin() *pcoin.PlayerCoin
	GetModuleCoinReadOnly() *pcoin.PlayerCoin
	GetModuleBag() *pbag.PlayerBag
	GetModuleBagReadOnly() *pbag.PlayerBag
}

type LogicGoods struct {
	player IGoodsPlayer
}

func NewLogicGoods(p IGoodsPlayer) *LogicGoods {
	return &LogicGoods{
		player: p,
	}
}

type AddGoodsParam struct {
}
type CostGoodsParam struct {
}

func GetBagTypeByGoodsType(goodsType int32) int32 {
	var bagType = pbag.EBagTypeNone
	switch goodsType {
	case goods.EGoodsTypeItem:
		{
			bagType = pbag.EBagTypeItem
			break
		}
	case goods.EGoodsTypeEquip:
		{
			bagType = pbag.EBagTypeEquip
			break
		}
	case goods.EGoodsTypeHero:
		{
			bagType = pbag.EBagTypeHero
			break
		}
	default:
		{
		}
	}
	return int32(bagType)
}

func (goodsMgr *LogicGoods) AddGoods(goodsID int32, num int64, goodsFrom int32, optParam *AddGoodsParam) bool {
	var addRet bool = false
	goodsType := goods.GetGoodsType(goodsID)
	if goodsType == goods.EGoodsTypeCoin {
		coinModule := goodsMgr.player.GetModuleCoin()
		if coinModule == nil {
			return false
		}
		coinModule.AddNum(goodsID, num)
		return addRet
	}

	playerBag := goodsMgr.player.GetModuleBag()
	if playerBag == nil {
		return false
	}
	var bagType = GetBagTypeByGoodsType(goodsType)
	thisBag := playerBag.GetBag(bagType)
	if thisBag == nil {
		return false
	}
	thisBag.AddGoods(goodsID, goodsType, num, 0)
	return addRet
}
func (goodsMgr *LogicGoods) GetGoodsNum(goodsID int32) int64 {
	goodsType := goods.GetGoodsType(goodsID)
	if goodsType == goods.EGoodsTypeCoin {
		coinModule := goodsMgr.player.GetModuleCoinReadOnly()
		if coinModule == nil {
			return 0
		}
		return coinModule.GetNum(goodsID)
	}
	var goodsNum int64 = 0
	playerBag := goodsMgr.player.GetModuleBagReadOnly()
	if playerBag == nil {
		return goodsNum
	}
	var bagType = GetBagTypeByGoodsType(goodsType)
	thisBag := playerBag.GetBag(bagType)
	if thisBag == nil {
		loghlp.Errorf("getGoodsNumError, not find bag,bagType(%d),goodsID(%d)", bagType, goodsID)
		return goodsNum
	}
	goodsNum = thisBag.GetGoodsNum(goodsID)
	return goodsNum
}
func (goodsMgr *LogicGoods) CostGoods(goodsID int32, num int64, goodsFrom int32, optParam *CostGoodsParam) bool {
	goodsType := goods.GetGoodsType(goodsID)
	if goodsType == goods.EGoodsTypeCoin {
		coinModule := goodsMgr.player.GetModuleCoin()
		if coinModule == nil {
			return false
		}
		if coinModule.GetNum(goodsID) < num {
			return false
		}
		coinModule.CostNum(goodsID, num)
		return true
	}
	// 常规背包道具
	playerBag := goodsMgr.player.GetModuleBag()
	if playerBag == nil {
		return false
	}
	var bagType = GetBagTypeByGoodsType(goodsType)
	thisBag := playerBag.GetBag(bagType)
	if thisBag == nil {
		loghlp.Errorf("costGoodsError, not find bag,bagType(%d),goodsID(%d)", bagType, goodsID)
		return false
	}
	return thisBag.CostGoods(goodsID, num)
}
