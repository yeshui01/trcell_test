/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 17:14:11
 * @LastEditTime: 2022-09-19 17:16:52
 * @FilePath: \trcell\pkg\trframe\trnode\tcellview\tnode_cellview.go
 */
package tcellview

import (
	"trcell/pkg/loghlp"
	"trcell/pkg/trframe/iframe"
	"trcell/pkg/trframe/trnode"
)

// view 节点
type FrameNodeCellView struct {
	tframeObj iframe.ITRFrame
	nodeIndex int32
}

func New(frameObj iframe.ITRFrame, index int32) *FrameNodeCellView {
	return &FrameNodeCellView{
		tframeObj: frameObj,
		nodeIndex: index,
	}
}

func (frameNode *FrameNodeCellView) RunStepCheck(curTimeMs int64) bool {
	loghlp.Info("frame node run step check")
	return true
}

func (frameNode *FrameNodeCellView) RunStepInit(curTimeMs int64) bool {
	loghlp.Info("frame node run step init")
	return frameNode.InitConnectServer()
}
func (frameNode *FrameNodeCellView) RunStepPreRun(curTimeMs int64) bool {
	loghlp.Info("frame node run step preRun")
	return true
}
func (frameNode *FrameNodeCellView) RunStepRun(curTimeMs int64) bool {

	return true
}
func (frameNode *FrameNodeCellView) RunStepStop(curTimeMs int64) bool {
	loghlp.Info("frame node run step stop")
	return true
}
func (frameNode *FrameNodeCellView) RunStepEnd(curTimeMs int64) bool {
	loghlp.Info("frame node run step end")
	return true
}
func (frameNode *FrameNodeCellView) NodeType() int32 {
	return trnode.ETRNodeTypeCellView
}
func (frameNode *FrameNodeCellView) NodeIndex() int32 {
	return frameNode.nodeIndex
}
