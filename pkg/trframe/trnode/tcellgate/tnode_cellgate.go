/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-16 18:10:33
 * @LastEditTime: 2022-09-16 18:10:41
 * @FilePath: \trcell\pkg\trframe\trnode\tcellgate\tnode_cellgate.go
 */
package tcellgate

import (
	"trcell/pkg/loghlp"
	"trcell/pkg/trframe/iframe"
	"trcell/pkg/trframe/trnode"
)

// gate 节点
type FrameNodeCellGate struct {
	tframeObj iframe.ITRFrame
	nodeIndex int32
}

func New(frameObj iframe.ITRFrame, index int32) *FrameNodeCellGate {
	return &FrameNodeCellGate{
		tframeObj: frameObj,
		nodeIndex: index,
	}
}

func (frameNode *FrameNodeCellGate) RunStepCheck(curTimeMs int64) bool {
	loghlp.Info("frame node run step check")
	return true
}

func (frameNode *FrameNodeCellGate) RunStepInit(curTimeMs int64) bool {
	loghlp.Info("frame node run step init")
	return frameNode.InitConnectServer()
}
func (frameNode *FrameNodeCellGate) RunStepPreRun(curTimeMs int64) bool {
	loghlp.Info("frame node run step preRun")
	return true
}
func (frameNode *FrameNodeCellGate) RunStepRun(curTimeMs int64) bool {

	return true
}
func (frameNode *FrameNodeCellGate) RunStepStop(curTimeMs int64) bool {
	loghlp.Info("frame node run step stop")
	return true
}
func (frameNode *FrameNodeCellGate) RunStepEnd(curTimeMs int64) bool {
	loghlp.Info("frame node run step end")
	return true
}
func (frameNode *FrameNodeCellGate) NodeType() int32 {
	return trnode.ETRNodeTypeCellGate
}
func (frameNode *FrameNodeCellGate) NodeIndex() int32 {
	return frameNode.nodeIndex
}
