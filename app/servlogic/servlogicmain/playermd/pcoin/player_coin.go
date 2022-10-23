/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-28 19:18:45
 * @LastEditTime: 2022-09-29 13:37:01
 * @FilePath: \trcell\app\servlogic\servlogicmain\playermd\pcoin\player_coin.go
 */
package pcoin

import (
	"encoding/json"
	"trcell/app/servlogic/servlogicmain/playermd"
	"trcell/pkg/ormdef"
	"trcell/pkg/pb/pbserver"
)

type PlayerCoin struct {
	player   playermd.ILogicPlayer
	CoinList map[int32]int64
}

func NewPlayerCoin(p playermd.ILogicPlayer) *PlayerCoin {
	return &PlayerCoin{
		player:   p,
		CoinList: make(map[int32]int64),
	}
}

func (saveObj *PlayerCoin) ToBytes() []byte {
	dbRoleCoin := &pbserver.DBRoleCoin{}
	dbRoleCoin.CoinList = make(map[int32]int64)
	for k, v := range saveObj.CoinList {
		dbRoleCoin.CoinList[k] = v
	}
	data, _ := json.Marshal(dbRoleCoin)
	return data
}
func (saveObj *PlayerCoin) FromBytes(binaryData []byte) {
	// json
	if binaryData != nil {
		dbRoleCoin := &pbserver.DBRoleCoin{}
		json.Unmarshal(binaryData, dbRoleCoin)
		for k, v := range dbRoleCoin.CoinList {
			saveObj.CoinList[k] = v
		}
	}
}

func (saveObj *PlayerCoin) UpdToDB() {
	saveObj.player.SetTableUpd(ormdef.ETableRoleCoin)
}

func (saveObj *PlayerCoin) GetNum(coinID int32) int64 {
	if n, ok := saveObj.CoinList[coinID]; ok {
		return n
	}
	return 0
}
func (saveObj *PlayerCoin) AddNum(coinID int32, num int64) int64 {
	if n, ok := saveObj.CoinList[coinID]; ok {
		saveObj.CoinList[coinID] = n + num
		return n + num
	}
	saveObj.CoinList[coinID] = num
	return num
}

func (saveObj *PlayerCoin) CostNum(coinID int32, num int64) (int64, bool) {
	if n, ok := saveObj.CoinList[coinID]; ok {
		newNum := n - num
		if newNum < 0 {
			newNum = 0
		}
		saveObj.CoinList[coinID] = newNum
		return newNum, true
	}
	return 0, false
}
