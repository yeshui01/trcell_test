/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 14:02:42
 * @LastEditTime: 2022-09-19 14:33:57
 * @FilePath: \trcell\pkg\trframe\trnode\tcelllogic\tnode_celllogic.go
 */
package tcelllogic

import (
	"trcell/pkg/loghlp"
	"trcell/pkg/trframe/iframe"
	"trcell/pkg/trframe/trnode"
)

// logic 节点
type FrameNodeCellLogic struct {
	tframeObj iframe.ITRFrame
	nodeIndex int32
}

func New(frameObj iframe.ITRFrame, index int32) *FrameNodeCellLogic {
	return &FrameNodeCellLogic{
		tframeObj: frameObj,
		nodeIndex: index,
	}
}

func (frameNode *FrameNodeCellLogic) RunStepCheck(curTimeMs int64) bool {
	loghlp.Info("frame node run step check")
	return true
}

func (frameNode *FrameNodeCellLogic) RunStepInit(curTimeMs int64) bool {
	loghlp.Info("frame node run step init")
	return frameNode.InitConnectServer()
}
func (frameNode *FrameNodeCellLogic) RunStepPreRun(curTimeMs int64) bool {
	loghlp.Info("frame node run step preRun")
	return true
}
func (frameNode *FrameNodeCellLogic) RunStepRun(curTimeMs int64) bool {

	return true
}
func (frameNode *FrameNodeCellLogic) RunStepStop(curTimeMs int64) bool {
	loghlp.Info("frame node run step stop")
	return true
}
func (frameNode *FrameNodeCellLogic) RunStepEnd(curTimeMs int64) bool {
	loghlp.Info("frame node run step end")
	return true
}
func (frameNode *FrameNodeCellLogic) NodeType() int32 {
	return trnode.ETRNodeTypeCellLogic
}
func (frameNode *FrameNodeCellLogic) NodeIndex() int32 {
	return frameNode.nodeIndex
}
