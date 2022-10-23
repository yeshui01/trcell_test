/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 15:05:17
 * @LastEditTime: 2022-09-19 15:05:24
 * @FilePath: \trcell\pkg\trframe\trnode\tcellgame\tnode_cellgame.go
 */
package tcellgame

import (
	"trcell/pkg/loghlp"
	"trcell/pkg/trframe/iframe"
	"trcell/pkg/trframe/trnode"
)

// game 节点
type FrameNodeCellGame struct {
	tframeObj iframe.ITRFrame
	nodeIndex int32
}

func New(frameObj iframe.ITRFrame, index int32) *FrameNodeCellGame {
	return &FrameNodeCellGame{
		tframeObj: frameObj,
		nodeIndex: index,
	}
}

func (frameNode *FrameNodeCellGame) RunStepCheck(curTimeMs int64) bool {
	loghlp.Info("frame node run step check")
	return true
}

func (frameNode *FrameNodeCellGame) RunStepInit(curTimeMs int64) bool {
	loghlp.Info("frame node run step init")
	return frameNode.InitConnectServer()
}
func (frameNode *FrameNodeCellGame) RunStepPreRun(curTimeMs int64) bool {
	loghlp.Info("frame node run step preRun")
	return true
}
func (frameNode *FrameNodeCellGame) RunStepRun(curTimeMs int64) bool {

	return true
}
func (frameNode *FrameNodeCellGame) RunStepStop(curTimeMs int64) bool {
	loghlp.Info("frame node run step stop")
	return true
}
func (frameNode *FrameNodeCellGame) RunStepEnd(curTimeMs int64) bool {
	loghlp.Info("frame node run step end")
	return true
}
func (frameNode *FrameNodeCellGame) NodeType() int32 {
	return trnode.ETRNodeTypeCellGame
}
func (frameNode *FrameNodeCellGame) NodeIndex() int32 {
	return frameNode.nodeIndex
}
