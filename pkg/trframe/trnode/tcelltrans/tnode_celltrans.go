/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-10-14 13:44:34
 * @LastEditTime: 2022-10-14 13:44:39
 * @FilePath: \trcell\pkg\trframe\trnode\tcelltrans\tnode_celltrans.go
 */
package tcelltrans

import (
	"trcell/pkg/loghlp"
	"trcell/pkg/trframe/iframe"
	"trcell/pkg/trframe/trnode"
)

// social 节点
type FrameNodeCellTrans struct {
	tframeObj iframe.ITRFrame
	nodeIndex int32
}

func New(frameObj iframe.ITRFrame, index int32) *FrameNodeCellTrans {
	return &FrameNodeCellTrans{
		tframeObj: frameObj,
		nodeIndex: index,
	}
}

func (frameNode *FrameNodeCellTrans) RunStepCheck(curTimeMs int64) bool {
	loghlp.Info("frame node run step check")
	return true
}

func (frameNode *FrameNodeCellTrans) RunStepInit(curTimeMs int64) bool {
	loghlp.Info("frame node run step init")
	return frameNode.InitConnectServer()
}
func (frameNode *FrameNodeCellTrans) RunStepPreRun(curTimeMs int64) bool {
	loghlp.Info("frame node run step preRun")
	return true
}
func (frameNode *FrameNodeCellTrans) RunStepRun(curTimeMs int64) bool {

	return true
}
func (frameNode *FrameNodeCellTrans) RunStepStop(curTimeMs int64) bool {
	loghlp.Info("frame node run step stop")
	return true
}
func (frameNode *FrameNodeCellTrans) RunStepEnd(curTimeMs int64) bool {
	loghlp.Info("frame node run step end")
	return true
}
func (frameNode *FrameNodeCellTrans) NodeType() int32 {
	return trnode.ETRNodeTypeCellTrans
}
func (frameNode *FrameNodeCellTrans) NodeIndex() int32 {
	return frameNode.nodeIndex
}
