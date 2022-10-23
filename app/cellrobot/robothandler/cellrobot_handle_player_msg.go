/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-20 11:25:49
 * @LastEditTime: 2022-09-20 11:35:17
 * @FilePath: \trcell\app\cellrobot\robothandler\cellrobot_handle_player_msg.go
 */
package robothandler

import (
	"trcell/app/cellrobot/robotcore"
	"trcell/pkg/evhub"
)

func RobotHandleECMsgRoomPushPlayerEnterNotify(robotIns robotcore.ICellRobot, sMsg *evhub.NetMessage) {
	// notifyMsg := &pbclient.ECMsgRoomPushPlayerEnterNotify{}
	// robotIns.LogRecvMsgInfo(sMsg, notifyMsg)
}

// 玩家离开房间
func RobotHandleECMsgRoomPushPlayerLeaveNotify(robotIns robotcore.ICellRobot, sMsg *evhub.NetMessage) {
	// notifyMsg := &pbclient.ECMsgRoomPushPlayerLeaveNotify{}
	// robotIns.LogRecvMsgInfo(sMsg, notifyMsg)
}

// 玩家房间聊天
func RobotHandleECMsgRoomPushChatNotify(robotIns robotcore.ICellRobot, sMsg *evhub.NetMessage) {
	// notifyMsg := &pbclient.ECMsgRoomPushChatNotify{}
	// robotIns.LogRecvMsgInfo(sMsg, notifyMsg)
}
