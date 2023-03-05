/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-16 14:18:14
 * @LastEditTime: 2022-10-14 13:42:13
 * @FilePath: \trcell\pkg\trframe\trnode\trframe_node.go
 */
package trnode

import "trcell/pkg/evhub"

// 节点类型
const (
	ETRNodeTypeNone       = 0
	ETRNodeTypeCellRoot   = 1   // 区域根节点
	ETRNodeTypeCellLog    = 2   // 日志节点
	ETRNodeTypeCellData   = 3   // 数据节点
	ETRNodeTypeCellCenter = 4   // 区域中心节点
	ETRNodeTypeCellGame   = 5   // 游戏节点
	ETRNodeTypeCellLogic  = 6   // 玩家逻辑节点
	ETRNodeTypeCellView   = 7   // 视图表现节点
	ETRNodeTypeCellGate   = 8   // 网关节点
	ETRNodeTypeCellSocial = 100 // 社交节点
	ETRNodeTypeCellTrans  = 200 // 传输节点
)

// 节点
type TRNodeInfo struct {
	ZoneID    int32
	NodeType  int32
	NodeIndex int32
	DesInfo   string
}

func NewNodeInfo(zoneID int32, nodeType int32, nodeIdx int32) *TRNodeInfo {
	return &TRNodeInfo{
		ZoneID:    zoneID,
		NodeType:  nodeType,
		NodeIndex: nodeIdx,
	}
}

// 节点实体
type ITRNodeEntity interface {
	GetNodeInfo() *TRNodeInfo
	Equal(zoneID int32, nodeType int32, nodeIndex int32) bool
	SendMsg(msg *evhub.NetMessage) bool
	LastHeartTime() int64
	SetHeartTime(int64)
	IsServerClient() bool
	GetSessionID() int32
}

// 网络端点
type PlayerNetPeer struct {
	netPeers map[int32]*TRNodeInfo
}

func NewPlayerNetPeer() *PlayerNetPeer {
	return &PlayerNetPeer{}
}

func (peerInfo *PlayerNetPeer) GetNode(nodeType int32) *TRNodeInfo {
	if peerInfo.netPeers == nil {
		return nil
	}
	if n, ok := peerInfo.netPeers[nodeType]; ok {
		return n
	}
	return nil
}
func (peerInfo *PlayerNetPeer) SetNode(nodeType int32, nodeInfo *TRNodeInfo) {
	if peerInfo.netPeers == nil {
		peerInfo.netPeers = make(map[int32]*TRNodeInfo)
	}
	peerInfo.netPeers[nodeType] = nodeInfo
}
func (peerInfo *PlayerNetPeer) GetNodeIndex(nodeType int32) int32 {
	if peerInfo.netPeers == nil {
		return 0
	}
	if n, ok := peerInfo.netPeers[nodeType]; ok {
		return n.NodeIndex
	}
	return 0
}

// 根据zoneID, nodeType, nodeIndex生成节点唯一id标识
func GenNodeUid(zoneID int32, nodeType int32, nodeIndex int32) int64 {
	return int64(zoneID)*100000000 + int64(nodeType)*10000 + int64(nodeIndex)
}

// 根据节点唯一id标识解析出zoneID, nodeType, nodeIndex
func GetNodePartIDByNodeUid(nodeUid int64) (int32, int32, int32) {
	zoneID := int32(nodeUid / 100000000)
	nodeType := int32(nodeUid % 100000000 / 10000)
	nodeIndex := int32(nodeUid % 10000)

	return zoneID, nodeType, nodeIndex
}
