/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-20 11:25:49
 * @LastEditTime: 2022-09-20 15:17:53
 * @FilePath: \trcell\app\cellrobot\robot_main.go
 */
package cellrobot

import (
	"trcell/app/cellrobot/robotcore"
)

func NewAiRobot(robotName string, aiType int32) *robotcore.CellRobot {
	r := robotcore.NewCellRobot(robotName)
	// r.TargetRoomID = targetRoomID
	// switch aiType {
	// case robotai.ERobotAiUndercover:
	// 	{
	// 		//r.AiInstance = robotai.NewAiUndercover(r)
	// 		break
	// 	}
	// }
	initRegisterRobotHandle(r)
	return r
}

// 注册基本消息处理
func initRegisterRobotHandle(robotObj *robotcore.CellRobot) {
	// game
	// robotObj.RegMsgHandle(protocol.ECMsgClassGame, protocol.ECMsgGamePushPlayerReadyStatus, robothandler.RobotHandleECMsgGamePushPlayerReadyStatusNotify)
}
