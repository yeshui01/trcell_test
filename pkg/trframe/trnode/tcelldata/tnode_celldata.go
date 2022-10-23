/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-16 18:17:35
 * @LastEditTime: 2022-09-19 11:24:15
 * @FilePath: \trcell\pkg\trframe\trnode\tcelldata\tnode_celldata.go
 */
package tcelldata

import (
	"trcell/pkg/loghlp"
	"trcell/pkg/trframe/iframe"
	"trcell/pkg/trframe/trnode"
)

// data 节点
type FrameNodeCellData struct {
	tframeObj iframe.ITRFrame
	nodeIndex int32
}

func New(frameObj iframe.ITRFrame, index int32) *FrameNodeCellData {
	return &FrameNodeCellData{
		tframeObj: frameObj,
		nodeIndex: index,
	}
}

func (frameNode *FrameNodeCellData) RunStepCheck(curTimeMs int64) bool {
	loghlp.Info("frame node run step check")
	return true
}

func (frameNode *FrameNodeCellData) RunStepInit(curTimeMs int64) bool {
	loghlp.Info("frame node run step init")
	return frameNode.InitConnectServer()
}
func (frameNode *FrameNodeCellData) RunStepPreRun(curTimeMs int64) bool {
	loghlp.Info("frame node run step preRun")
	return true
}
func (frameNode *FrameNodeCellData) RunStepRun(curTimeMs int64) bool {

	return true
}
func (frameNode *FrameNodeCellData) RunStepStop(curTimeMs int64) bool {
	loghlp.Info("frame node run step stop")
	return true
}
func (frameNode *FrameNodeCellData) RunStepEnd(curTimeMs int64) bool {
	loghlp.Info("frame node run step end")
	return true
}
func (frameNode *FrameNodeCellData) NodeType() int32 {
	return trnode.ETRNodeTypeCellData
}
func (frameNode *FrameNodeCellData) NodeIndex() int32 {
	return frameNode.nodeIndex
}
