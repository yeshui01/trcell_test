/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 11:34:54
 * @LastEditTime: 2022-09-19 11:36:20
 * @FilePath: \trcell\pkg\trframe\trnode\tcellcenter\tnode_cellcenter.go
 */
package tcellcenter

import (
	"trcell/pkg/loghlp"
	"trcell/pkg/trframe/iframe"
	"trcell/pkg/trframe/trnode"
)

// center 节点
type FrameNodeCellCenter struct {
	tframeObj iframe.ITRFrame
	nodeIndex int32
}

func New(frameObj iframe.ITRFrame, index int32) *FrameNodeCellCenter {
	return &FrameNodeCellCenter{
		tframeObj: frameObj,
		nodeIndex: index,
	}
}

func (frameNode *FrameNodeCellCenter) RunStepCheck(curTimeMs int64) bool {
	loghlp.Info("frame node run step check")
	return true
}

func (frameNode *FrameNodeCellCenter) RunStepInit(curTimeMs int64) bool {
	loghlp.Info("frame node run step init")
	return frameNode.InitConnectServer()
}
func (frameNode *FrameNodeCellCenter) RunStepPreRun(curTimeMs int64) bool {
	loghlp.Info("frame node run step preRun")
	return true
}
func (frameNode *FrameNodeCellCenter) RunStepRun(curTimeMs int64) bool {

	return true
}
func (frameNode *FrameNodeCellCenter) RunStepStop(curTimeMs int64) bool {
	loghlp.Info("frame node run step stop")
	return true
}
func (frameNode *FrameNodeCellCenter) RunStepEnd(curTimeMs int64) bool {
	loghlp.Info("frame node run step end")
	return true
}
func (frameNode *FrameNodeCellCenter) NodeType() int32 {
	return trnode.ETRNodeTypeCellCenter
}
func (frameNode *FrameNodeCellCenter) NodeIndex() int32 {
	return frameNode.nodeIndex
}
