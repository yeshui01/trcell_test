/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 13:44:42
 * @LastEditTime: 2022-09-29 15:48:29
 * @FilePath: \trcell\app\servcenter\servcentermain\servcenter_player.go
 */
package servcentermain

import "trcell/pkg/trframe/trnode"

type CenterPlayer struct {
	RoleID int64
	// 数据私有化
	nickName string
	level    int32
	//
	netPeers *trnode.PlayerNetPeer
}

func NewCenterPlayer(roleID int64, nickName string, lv int32) *CenterPlayer {
	return &CenterPlayer{
		RoleID:   roleID,
		nickName: nickName,
		level:    lv,
		netPeers: trnode.NewPlayerNetPeer(),
	}
}

func (player *CenterPlayer) SetNetPeer(nodeType int32, nodeInfo *trnode.TRNodeInfo) {
	player.netPeers.SetNode(nodeType, nodeInfo)
}
func (player *CenterPlayer) GetNetPeer(nodeType int32) *trnode.TRNodeInfo {
	return player.netPeers.GetNode(nodeType)
}
func (player *CenterPlayer) GetNetPeerIndex(nodeType int32) int32 {
	return player.netPeers.GetNodeIndex(nodeType)
}
