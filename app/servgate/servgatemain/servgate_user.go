/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 15:59:32
 * @LastEditTime: 2022-09-28 10:00:15
 * @FilePath: \trcell\app\servgate\servgatemain\servgate_user.go
 */
package servgatemain

import (
	"trcell/pkg/evhub"
	"trcell/pkg/loghlp"
	"trcell/pkg/trframe/trnode"
)

type HGateUser struct {
	userID      int64           // 账号id
	gateConnect *HGateConnction // 对应的连接
	netPeers    *trnode.PlayerNetPeer
}

func NewGateUser(id int64) *HGateUser {
	return &HGateUser{
		userID:   id,
		netPeers: trnode.NewPlayerNetPeer(),
	}
}
func (u *HGateUser) SetNetPeer(nodeType int32, nodeInfo *trnode.TRNodeInfo) {
	u.netPeers.SetNode(nodeType, nodeInfo)
}
func (u *HGateUser) GetNetPeer(nodeType int32) *trnode.TRNodeInfo {
	return u.netPeers.GetNode(nodeType)
}
func (u *HGateUser) GetNetPeerIndex(nodeType int32) int32 {
	return u.netPeers.GetNodeIndex(nodeType)
}

func (u *HGateUser) SetGateConnect(conn *HGateConnction) {
	u.gateConnect = conn
	conn.UserInfo = u
}
func (u *HGateUser) GetGateConnect() *HGateConnction {
	return u.gateConnect
}

func (u *HGateUser) SendMessageToSelf(netMessage *evhub.NetMessage) {
	if u.gateConnect != nil {
		u.gateConnect.SendMsg(netMessage)
	}
}

type HGateUserManager struct {
	userList map[int64]*HGateUser
	// gateServ *HallGate
}

func NewHGateUserManager() *HGateUserManager {
	return &HGateUserManager{
		userList: make(map[int64]*HGateUser),
		// gateServ: hgServ,
	}
}

func (mgr *HGateUserManager) AddGateUser(userID int64, gateUser *HGateUser) {
	mgr.userList[userID] = gateUser
}

func (mgr *HGateUserManager) DelGateUser(userID int64) {
	delete(mgr.userList, userID)
	loghlp.Warnf("DelGateUser:%d", userID)
}
func (mgr *HGateUserManager) GetGateUser(userID int64) *HGateUser {
	if p, ok := mgr.userList[userID]; ok {
		return p
	}
	return nil
}
