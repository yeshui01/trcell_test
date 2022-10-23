/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-10-10 17:59:43
 * @LastEditTime: 2022-10-10 17:59:50
 * @FilePath: \trcell\pkg\trframe\trnode\tcellsocial\tnode_cellsocial.go
 */
package tcellsocial

import (
	"trcell/pkg/loghlp"
	"trcell/pkg/trframe/iframe"
	"trcell/pkg/trframe/trnode"
)

// social 节点
type FrameNodeCellSocial struct {
	tframeObj iframe.ITRFrame
	nodeIndex int32
}

func New(frameObj iframe.ITRFrame, index int32) *FrameNodeCellSocial {
	return &FrameNodeCellSocial{
		tframeObj: frameObj,
		nodeIndex: index,
	}
}

func (frameNode *FrameNodeCellSocial) RunStepCheck(curTimeMs int64) bool {
	loghlp.Info("frame node run step check")
	return true
}

func (frameNode *FrameNodeCellSocial) RunStepInit(curTimeMs int64) bool {
	loghlp.Info("frame node run step init")
	return frameNode.InitConnectServer()
}
func (frameNode *FrameNodeCellSocial) RunStepPreRun(curTimeMs int64) bool {
	loghlp.Info("frame node run step preRun")
	return true
}
func (frameNode *FrameNodeCellSocial) RunStepRun(curTimeMs int64) bool {

	return true
}
func (frameNode *FrameNodeCellSocial) RunStepStop(curTimeMs int64) bool {
	loghlp.Info("frame node run step stop")
	return true
}
func (frameNode *FrameNodeCellSocial) RunStepEnd(curTimeMs int64) bool {
	loghlp.Info("frame node run step end")
	return true
}
func (frameNode *FrameNodeCellSocial) NodeType() int32 {
	return trnode.ETRNodeTypeCellSocial
}
func (frameNode *FrameNodeCellSocial) NodeIndex() int32 {
	return frameNode.nodeIndex
}
