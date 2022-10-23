/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-16 17:55:43
 * @LastEditTime: 2022-09-19 10:46:31
 * @FilePath: \trcell\pkg\trframe\trnode\tcellroot\tnode_rootserver.go
 */
package tcellroot

import (
	"trcell/pkg/loghlp"
	"trcell/pkg/trframe/iframe"
	"trcell/pkg/trframe/trnode"
)

// root 节点
type FrameNodeRootServer struct {
	tframeObj iframe.ITRFrame
	nodeIndex int32
}

func New(frameObj iframe.ITRFrame, index int32) *FrameNodeRootServer {
	return &FrameNodeRootServer{
		tframeObj: frameObj,
		nodeIndex: index,
	}
}

func (frameNode *FrameNodeRootServer) RunStepCheck(curTimeMs int64) bool {
	loghlp.Info("frame node run step check")
	return true
}

func (frameNode *FrameNodeRootServer) RunStepInit(curTimeMs int64) bool {
	loghlp.Info("frame node run step init")
	return frameNode.InitConnectServer()
}
func (frameNode *FrameNodeRootServer) RunStepPreRun(curTimeMs int64) bool {
	loghlp.Info("frame node run step preRun")
	return true
}
func (frameNode *FrameNodeRootServer) RunStepRun(curTimeMs int64) bool {
	return true
}
func (frameNode *FrameNodeRootServer) RunStepStop(curTimeMs int64) bool {
	loghlp.Info("frame node run step stop")
	return true
}
func (frameNode *FrameNodeRootServer) RunStepEnd(curTimeMs int64) bool {
	loghlp.Info("frame node run step end")
	return true
}
func (frameNode *FrameNodeRootServer) NodeType() int32 {
	return trnode.ETRNodeTypeCellRoot
}
func (frameNode *FrameNodeRootServer) NodeIndex() int32 {
	return frameNode.nodeIndex
}
