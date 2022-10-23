/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 17:55:36
 * @LastEditTime: 2022-09-19 17:57:01
 * @FilePath: \trcell\pkg\trframe\trnode\tcelllog\tnode_celllog.go
 */
package tcelllog

import (
	"trcell/pkg/loghlp"
	"trcell/pkg/trframe/iframe"
	"trcell/pkg/trframe/trnode"
)

// log 节点
type FrameNodeCellLog struct {
	tframeObj iframe.ITRFrame
	nodeIndex int32
}

func New(frameObj iframe.ITRFrame, index int32) *FrameNodeCellLog {
	return &FrameNodeCellLog{
		tframeObj: frameObj,
		nodeIndex: index,
	}
}

func (frameNode *FrameNodeCellLog) RunStepCheck(curTimeMs int64) bool {
	loghlp.Info("frame node run step check")
	return true
}

func (frameNode *FrameNodeCellLog) RunStepInit(curTimeMs int64) bool {
	loghlp.Info("frame node run step init")
	return frameNode.InitConnectServer()
}
func (frameNode *FrameNodeCellLog) RunStepPreRun(curTimeMs int64) bool {
	loghlp.Info("frame node run step preRun")
	return true
}
func (frameNode *FrameNodeCellLog) RunStepRun(curTimeMs int64) bool {

	return true
}
func (frameNode *FrameNodeCellLog) RunStepStop(curTimeMs int64) bool {
	loghlp.Info("frame node run step stop")
	return true
}
func (frameNode *FrameNodeCellLog) RunStepEnd(curTimeMs int64) bool {
	loghlp.Info("frame node run step end")
	return true
}
func (frameNode *FrameNodeCellLog) NodeType() int32 {
	return trnode.ETRNodeTypeCellLog
}
func (frameNode *FrameNodeCellLog) NodeIndex() int32 {
	return frameNode.nodeIndex
}
