/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-28 18:36:10
 * @LastEditTime: 2022-09-29 13:40:38
 * @FilePath: \trcell\app\servlogic\servlogicmain\playermd\pbag\goods_bag.go
 */

package pbag

type BagPos struct {
	PosID     int32 // 格子ID
	GoodsType int32 // 物品类型
	ExcelID   int32 // 物品配置表id
	Num       int64 // 数量
	AttachID  int64 // 关联id
}

type GoodsBag struct {
	BagType int32
	PosIdx  int32 // 格子生成索引
	PosList map[int32]*BagPos
}

func NewGoodsBag(bagType int32) *GoodsBag {
	return &GoodsBag{
		BagType: bagType,
		PosIdx:  0,
		PosList: make(map[int32]*BagPos),
	}
}
func (bagObj *GoodsBag) newBagPos() *BagPos {
	bagObj.PosIdx++
	var posID int32 = bagObj.BagType*10000 + bagObj.PosIdx
	return &BagPos{
		PosID: posID,
		Num:   0,
	}
}
func (bagObj *GoodsBag) getPos(posID int32) *BagPos {
	if pos, ok := bagObj.PosList[posID]; ok {
		return pos
	}
	return nil
}
func (bagObj *GoodsBag) AddGoods(excelID int32, goodsType int32, num int64, attachID int64) *BagPos {
	var goodsPos *BagPos = nil
	var newPos bool = true // 是否需要重新生成格子
	if attachID == 0 {
		// 查找原来的叠加
		for _, v := range bagObj.PosList {
			if v.ExcelID == excelID {
				v.Num = v.Num + num
				goodsPos = v
				newPos = false
				break
			}
		}
	} else {
		newPos = true
	}

	if newPos {
		goodsPos = bagObj.newBagPos()
		bagObj.PosList[goodsPos.PosID] = goodsPos
		goodsPos.ExcelID = excelID
		goodsPos.Num = num
		goodsPos.GoodsType = int32(goodsType)
		goodsPos.AttachID = attachID
	}
	return goodsPos
}

func (bagObj *GoodsBag) CostGoodsByPos(posID int32, costNum int64) bool {
	goodsPos := bagObj.getPos(posID)
	if goodsPos == nil {
		return false
	}
	if goodsPos.Num < costNum {
		return false
	}
	goodsPos.Num = goodsPos.Num - costNum
	return true
}
func (bagObj *GoodsBag) CostGoods(excelID int32, costNum int64) bool {
	var goodsPos *BagPos = nil
	for _, v := range bagObj.PosList {
		if v.ExcelID == excelID {
			goodsPos = v
			break
		}
	}
	if goodsPos == nil {
		return false
	}
	if goodsPos.Num < costNum {
		return false
	}
	goodsPos.Num = goodsPos.Num - costNum
	return true
}

func (bagObj *GoodsBag) GetGoodsNum(excelID int32) int64 {
	var totalNum int64 = 0
	for _, v := range bagObj.PosList {
		if v.ExcelID == excelID {
			totalNum = totalNum + v.Num
			break
		}
	}
	return totalNum
}
