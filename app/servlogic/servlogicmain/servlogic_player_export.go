/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-29 11:47:07
 * @LastEditTime: 2022-09-29 18:33:25
 * @FilePath: \trcell\app\servlogic\servlogicmain\servlogic_player_export.go
 */
package servlogicmain

import (
	"trcell/app/servlogic/servlogicmain/playergoods"
	"trcell/app/servlogic/servlogicmain/playermd"
	"trcell/app/servlogic/servlogicmain/playermd/pbag"
	"trcell/app/servlogic/servlogicmain/playermd/pcoin"
	"trcell/app/servlogic/servlogicmain/playermd/pextra"
)

func (player *LogicPlayer) GetGoodsMgr() *playergoods.LogicGoods {
	return player.goodsMgr
}

func (player *LogicPlayer) GetModuleCoin() *pcoin.PlayerCoin {
	coinMD, err := player.GetModule(playermd.EPlayerModuleCoin)
	if err == nil {
		playerCoin := coinMD.(*pcoin.PlayerCoin)
		return playerCoin
	}
	return nil
}
func (player *LogicPlayer) GetModuleCoinReadOnly() *pcoin.PlayerCoin {
	coinMD, err := player.GetModuleReadOnly(playermd.EPlayerModuleCoin)
	if err == nil {
		playerCoin := coinMD.(*pcoin.PlayerCoin)
		return playerCoin
	}
	return nil
}
func (player *LogicPlayer) GetModuleBag() *pbag.PlayerBag {
	bagMD, err := player.GetModule(playermd.EPlayerModuleBag)
	if err == nil {
		playerBag := bagMD.(*pbag.PlayerBag)
		return playerBag
	}
	return nil
}
func (player *LogicPlayer) GetModuleBagReadOnly() *pbag.PlayerBag {
	bagMD, err := player.GetModuleReadOnly(playermd.EPlayerModuleBag)
	if err == nil {
		playerBag := bagMD.(*pbag.PlayerBag)
		return playerBag
	}
	return nil
}
func (player *LogicPlayer) GenScopeUid() int64 {
	extraMD, errExtra := player.GetModule(playermd.EPlayerModuleExtra)
	if errExtra == nil {
		extraModule, ok := extraMD.(*pextra.PlayerExtraData)
		if ok {
			return extraModule.GenScopeUid()
		}
	}
	return 0
}
