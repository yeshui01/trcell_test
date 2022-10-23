/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-28 16:15:30
 * @LastEditTime: 2022-09-28 17:59:46
 * @FilePath: \trcell\pkg\gamemd\itask\task_interface.go
 */
package itask

import "trcell/pkg/configdata"

// 任务状态
const (
	ETaskStatusNone     = 0
	ETaskStatusAccept   = 1 // 接取状态
	ETaskStatusComplete = 2 // 进度条件完成状态
	ETaskStatusFinish   = 3 // 结束
)

// 进度更新逻辑类型
const (
	ETaskProgressLogicAdd    = 1 // 增加
	ETaskProgressLogicRecalc = 2 // 重新计算
	ETaskProgressLogicUpd    = 3 // 覆盖更新
)

type TargetData struct {
	Num        int64 // 当前进度
	TargetType int32
}

type TaskDetailData struct {
	TaskID     int32
	TargetList []TargetData
	Status     int32
}

func NewTaskDetailData(taskID int32) *TaskDetailData {
	return &TaskDetailData{
		TargetList: make([]TargetData, 1), // 目前就一个
		TaskID:     taskID,
	}
}

// 清理进度
func (obj *TaskDetailData) ClearNum() {
	for i := 0; i < len(obj.TargetList); i++ {
		obj.TargetList[i].Num = 0
		obj.TargetList[i].TargetType = 0
	}
}
func (obj *TaskDetailData) IsComplete(taskCfg *configdata.TaskCfg) bool {
	var reachNum int32 = 0
	for i := 0; i < len(obj.TargetList) && i < len(taskCfg.Targets); i++ {
		if GetCompareLogic(taskCfg.Targets[i].TargetType) {
			if obj.TargetList[i].Num >= taskCfg.Targets[i].TotalNum {
				reachNum++
			}
		} else {
			if obj.TargetList[i].Num <= taskCfg.Targets[i].TotalNum {
				reachNum++
			}
		}
	}
	return reachNum >= int32(len(taskCfg.Targets))
}

type TaskOptEnv struct {
}

type ITaskOwner interface {
	// GetTaskData(taskID int32) *TargetData
	// ChangeTaskStatus(taskID int32, status int32, param1 int32, param2 int32, param3 int64, param4 int64, taskEnv *TaskOptEnv) *TaskOptEnv
	RecalcTaskProgress(taskTarget *configdata.TaskTarget, optEnv *TaskOptEnv) int64
}

// 获取逻辑类型,true 进度值>=需要的值,false反之
func GetCompareLogic(targetType int32) bool {
	// return targetType != task.ETaskTargetArenaRank
	return false
}
func GetProgressLogicType(targetType int32, isPermanentAdd bool) int32 {
	if isPermanentAdd {
		return ETaskProgressLogicAdd
	}
	switch targetType {
	case ETaskTargetNone:
		return ETaskProgressLogicAdd
	default:
		break
	}
	return ETaskProgressLogicAdd
}
