/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 14:24:50
 * @LastEditTime: 2022-09-29 18:22:26
 * @FilePath: \trcell\app\servlogic\servlogicmain\servlogic_player.go
 */
package servlogicmain

import (
	"fmt"
	"time"
	"trcell/app/servlogic/servlogicmain/playergoods"
	"trcell/app/servlogic/servlogicmain/playermd"
	"trcell/app/servlogic/servlogicmain/playermd/pbag"
	"trcell/app/servlogic/servlogicmain/playermd/pcoin"
	"trcell/app/servlogic/servlogicmain/playermd/pequip"
	"trcell/app/servlogic/servlogicmain/playermd/pextra"
	"trcell/pkg/loghlp"
	"trcell/pkg/ormdef"
	"trcell/pkg/pb/pbclient"
	"trcell/pkg/pb/pbframe"
	"trcell/pkg/pb/pbserver"
	"trcell/pkg/protocol"
	"trcell/pkg/tbobj"
	"trcell/pkg/timeutil"
	"trcell/pkg/trframe"
	"trcell/pkg/trframe/trnode"

	"google.golang.org/protobuf/proto"
)

type LogicPlayer struct {
	RoleID           int64
	LastSaveTime     int64 // 最近保存的时间
	LastRealSaveTime int64 // 最近真实保存的时间
	OfflineTime      int64 // 最近离线时间
	IsOnline         bool  // 是否在线
	// 玩家DB数据,必须私有化
	// -----playertable codetag7 begin-----------------
	tbRoleBase  *tbobj.TbRoleBase
	tbRoleEquip *tbobj.TbRoleEquip
	tbRoleExtra *tbobj.TbRoleExtra
	tbRoleCoin  *tbobj.TbRoleCoin
	tbRoleBag   *tbobj.TbRoleBag
	// -----playertable codetag7 end-------------------
	// 玩家物品系统
	goodsMgr *playergoods.LogicGoods

	// private
	netPeers      *trnode.PlayerNetPeer
	lastHeartTime int64 // 最近心跳时间
	moduleList    []playermd.IPlayerModule
}

func NewLogicPlayer(roleID int64) *LogicPlayer {
	player := &LogicPlayer{
		RoleID:     roleID,
		tbRoleBase: tbobj.NewTbRoleBase(),

		// -----playertable codetag10 begin-----------------
		tbRoleEquip: tbobj.NewTbRoleEquip(),
		tbRoleExtra: tbobj.NewTbRoleExtra(),
		tbRoleCoin:  tbobj.NewTbRoleCoin(),
		tbRoleBag:   tbobj.NewTbRoleBag(),
		// -----playertable codetag10 end-------------------
		lastHeartTime: 0,
		moduleList:    make([]playermd.IPlayerModule, playermd.EPlayerModuleMax),
		netPeers:      trnode.NewPlayerNetPeer(),
	}
	player.goodsMgr = playergoods.NewLogicGoods(player)
	// 初始化需要存储数据的业务模块
	player.InitModules()
	return player
}
func (player *LogicPlayer) InitModules() {
	{
		// extra
		oneModule := pextra.NewPlayerExtraData(player)
		player.tbRoleExtra.BindModule(oneModule)
		player.AddModule(playermd.EPlayerModuleExtra, oneModule)
	}
	{
		// 装备
		oneModule := pequip.NewLogicEquip(player)
		player.tbRoleEquip.BindModule(oneModule)
		player.AddModule(playermd.EPlayerModuleEquip, oneModule)
	}
	{
		// 货币
		oneModule := pcoin.NewPlayerCoin(player)
		player.tbRoleCoin.BindModule(oneModule)
		player.AddModule(playermd.EPlayerModuleCoin, oneModule)
	}
	{
		// 背包
		oneModule := pbag.NewPlayerBag(player)
		player.tbRoleBag.BindModule(oneModule)
		player.AddModule(playermd.EPlayerModuleBag, oneModule)
	}
}

// 加载数据
func (player *LogicPlayer) LoadData(pbGameRoleData *pbserver.GameRoleData) {
	for _, tbData := range pbGameRoleData.RoleTables {
		switch tbData.TableID {
		case ormdef.ETableRoleBase:
			{
				player.tbRoleBase.FromBytes(tbData.Data)
				player.RoleID = player.tbRoleBase.GetRoleID()
				break
			}
			// -----playertable codetag4 begin-----------------
		case ormdef.ETableRoleEquip:
			{
				player.tbRoleEquip.FromBytes(tbData.Data)
				break
			}
		case ormdef.ETableRoleExtra:
			{
				player.tbRoleExtra.FromBytes(tbData.Data)
				break
			}
		case ormdef.ETableRoleCoin:
			{
				player.tbRoleCoin.FromBytes(tbData.Data)
				break
			}
		case ormdef.ETableRoleBag:
			{
				player.tbRoleBag.FromBytes(tbData.Data)
				break
			}
			// -----playertable codetag4 end-------------------
		default:
			{
				loghlp.Errorf("unknown table id:%d", tbData.TableID)
			}
		}
	}
	player.LastSaveTime = time.Now().Unix()
	// 各个模块反序列化数据
}

// 各个模块Dump数据
func (player *LogicPlayer) UpdAllModuleData() {
	// 各个模块序列化数据
	for _, md := range player.moduleList {
		md.UpdToDB()
	}
}

// 获取角色返回给客户端的数据
func (player *LogicPlayer) ToClientRoleInfo() *pbclient.RoleInfo {
	return &pbclient.RoleInfo{
		RoleID:   player.tbRoleBase.GetRoleID(),
		Level:    player.tbRoleBase.GetLevel(),
		Nickname: player.tbRoleBase.GetRoleName(),
	}
}

func (player *LogicPlayer) GetBaseData() *tbobj.TbRoleBase {
	return player.tbRoleBase
}

func (player *LogicPlayer) Online() {
	loghlp.Infof("player(%d) online", player.RoleID)
	player.GetBaseData().SetLoginTime(timeutil.NowTime())
	player.UpdateHeartTime(timeutil.NowTime())
	player.IsOnline = true
	// test code,送一个道具
	{
		equipMD, errEquip := player.GetModule(playermd.EPlayerModuleEquip)
		if errEquip == nil {
			equipModule, ok := equipMD.(*pequip.LogicEquip)
			if ok {
				if equipModule.EquipList.Size() == 0 {
					loghlp.Debugf("player(%d) init first equip data", player.RoleID)
					equipModule.InitFirstEquip()
				}
			}
		}
		// 生成100个uid
		extraMD, errExtra := player.GetModule(playermd.EPlayerModuleExtra)
		if errExtra == nil {
			extraModule, ok := extraMD.(*pextra.PlayerExtraData)
			if ok {
				for k := 0; k < 100; k++ {
					loghlp.Debugf("player(%d) genscopeuid:%d", player.RoleID, extraModule.GenScopeUid())
				}
			}
		}
		// 增加道具
		coinMD, errCoin := player.GetModule(playermd.EPlayerModuleCoin)
		if errCoin == nil {
			coinModule, ok := coinMD.(*pcoin.PlayerCoin)
			if ok {
				coinModule.AddNum(10001, 1)
				coinModule.AddNum(10002, 1)
			}
		}
	}
}

func (player *LogicPlayer) Offline() {
	loghlp.Infof("player(%d) offline", player.RoleID)
	player.OfflineTime = timeutil.NowTime()
	player.GetBaseData().SetOfflineTime(player.OfflineTime)
	player.IsOnline = false
	player.GetBaseData().SetOfflineTime(timeutil.NowTime())
	player.UpdAllModuleData()
}

func (player *LogicPlayer) SecUpdate(curTime int64) {

}

func (player *LogicPlayer) UpdateHeartTime(heartTime int64) {
	player.lastHeartTime = heartTime
}

func (player *LogicPlayer) GetHeartTime() int64 {
	return player.lastHeartTime
}

func (player *LogicPlayer) GetRoleID() int64 {
	return player.RoleID
}

func (player *LogicPlayer) SetTableUpd(tableID int32) {
	var updTime int64 = trframe.GetFrameSysNowTime()
	switch tableID {
	case ormdef.ETableRoleBase:
		{
			player.tbRoleBase.SetUpdTime(updTime)
			break
		}
		// -----playertable codetag3 begin-----------------
	case ormdef.ETableRoleEquip:
		{
			player.tbRoleEquip.SetUpdTime(updTime)
			break
		}
	case ormdef.ETableRoleExtra:
		{
			player.tbRoleExtra.SetUpdTime(updTime)
			break
		}
	case ormdef.ETableRoleCoin:
		{
			player.tbRoleCoin.SetUpdTime(updTime)
			break
		}
	case ormdef.ETableRoleBag:
		{
			player.tbRoleBag.SetUpdTime(updTime)
			break
		}
		// -----playertable codetag3 end-------------------
	default:
		{
			loghlp.Errorf("unknown table id:%d", tableID)
		}
	}
}

func (player *LogicPlayer) AddModule(moduleID int32, ipmd playermd.IPlayerModule) {
	if moduleID < playermd.EPlayerModuleExtra || moduleID >= playermd.EPlayerModuleMax {
		loghlp.Errorf("player addmodule(%d) fail, moduleid is out of range", moduleID)
		return
	}
	player.moduleList[moduleID] = ipmd
}

func (player *LogicPlayer) GetModule(moduleID int32) (playermd.IPlayerModule, error) {
	if moduleID < playermd.EPlayerModuleExtra || moduleID >= playermd.EPlayerModuleMax {
		return nil, fmt.Errorf("moduleID(%d) out of range", moduleID)
	}
	playerModule := player.moduleList[moduleID]
	if playerModule != nil {
		playerModule.UpdToDB() // 访问过的,默认为需要保存
	} else {
		return nil, fmt.Errorf("moduleID(%d) is nil", moduleID)
	}
	return playerModule, nil
}
func (player *LogicPlayer) GetModuleReadOnly(moduleID int32) (playermd.IPlayerModule, error) {
	if moduleID < playermd.EPlayerModuleExtra || moduleID >= playermd.EPlayerModuleMax {
		return nil, fmt.Errorf("moduleID(%d) out of range", moduleID)
	}
	playerModule := player.moduleList[moduleID]
	if playerModule == nil {
		return nil, fmt.Errorf("moduleID(%d) is nil", moduleID)
	}
	return playerModule, nil
}
func (player *LogicPlayer) SetNetPeer(nodeType int32, nodeInfo *trnode.TRNodeInfo) {
	player.netPeers.SetNode(nodeType, nodeInfo)
}
func (player *LogicPlayer) GetNetPeer(nodeType int32) *trnode.TRNodeInfo {
	return player.netPeers.GetNode(nodeType)
}
func (player *LogicPlayer) GetNetPeerIndex(nodeType int32) int32 {
	return player.netPeers.GetNodeIndex(nodeType)
}

/**
 * @description: 发送消息给玩家自己对应的客户端
 * @param {int32} msgClass
 * @param {int32} msgType
 * @param {proto.Message} pbMsg
 * @return {*}
 */
func (player *LogicPlayer) SendMsgToSelfClient(msgClass int32, msgType int32, pbMsg proto.Message) {
	msgData, _ := proto.Marshal(pbMsg)
	req := &pbframe.EFrameMsgPushMsgToClientReq{
		MsgClass: msgClass,
		MsgType:  msgType,
		RoleID:   player.RoleID,
		MsgData:  msgData,
	}
	trframe.PushZoneMessage(protocol.EMsgClassFrame,
		protocol.EFrameMsgPushMsgToClient,
		req,
		trnode.ETRNodeTypeCellGate,
		player.GetNetPeerIndex(trnode.ETRNodeTypeCellGate),
	)
}
