/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-28 16:24:07
 * @LastEditTime: 2022-09-28 18:00:23
 * @FilePath: \trcell\pkg\gamemd\itask\task_module.go
 */
package itask

import (
	"trcell/pkg/configdata"
	"trcell/pkg/loghlp"

	"github.com/sirupsen/logrus"
)

type TaskModuleBase struct {
	taskOwner    ITaskOwner
	taskList     map[int32]*TaskDetailData
	acceptList   map[int32]*TaskDetailData // 当前接取的任务列表
	completeList map[int32]*TaskDetailData
	targetToTask map[int32]map[int32]*TaskDetailData // key1:targetType key2:taskID
}

func NewTaskModuleBase(owner ITaskOwner) *TaskModuleBase {
	return &TaskModuleBase{
		taskOwner:    owner,
		taskList:     make(map[int32]*TaskDetailData),
		acceptList:   make(map[int32]*TaskDetailData),
		completeList: make(map[int32]*TaskDetailData),
		targetToTask: make(map[int32]map[int32]*TaskDetailData),
	}
}

func (taskMd *TaskModuleBase) GetTaskData(taskID int32) *TaskDetailData {
	if data, ok := taskMd.taskList[taskID]; ok {
		return data
	}
	return nil
}
func (taskMd *TaskModuleBase) InitTaskTarget(taskData *TaskDetailData, taskCfg *configdata.TaskCfg) {
	if len(taskData.TargetList) > len(taskCfg.Targets) {
		taskData.TargetList = taskData.TargetList[0:len(taskCfg.Targets)]
	} else if len(taskData.TargetList) < len(taskCfg.Targets) {
		for i := 0; i < (len(taskCfg.Targets) - len(taskData.TargetList)); i++ {
			taskData.TargetList = append(taskData.TargetList, TargetData{
				Num:        0,
				TargetType: 0,
			})
		}
	}
	for i, v := range taskCfg.Targets {
		taskData.TargetList[i].Num = 0
		taskData.TargetList[i].TargetType = v.TargetType
	}
}

func (taskMd *TaskModuleBase) ChangeTaskStatus(taskID int32, status int32, optEnv *TaskOptEnv) *TaskDetailData {
	taskData := taskMd.GetTaskData(taskID)
	if taskData == nil {
		// 初始化对象
		taskCfg := configdata.Instance().GetTaskCfg(taskID)
		if taskCfg == nil {
			loghlp.Errorf("not find task cfg:%v", taskID)
			return nil
		}
		taskData = NewTaskDetailData(taskCfg.TaskID)
		// 初始化目标参数
		taskMd.InitTaskTarget(taskData, taskCfg)
		taskMd.taskList[taskData.TaskID] = taskData
	}
	taskData.Status = status
	switch status {
	case ETaskStatusAccept:
		{
			delete(taskMd.completeList, taskID)
			taskMd.acceptList[taskID] = taskData
			taskData.ClearNum()
			for _, oneTarget := range taskData.TargetList {
				if taskMap, ok := taskMd.targetToTask[oneTarget.TargetType]; ok {
					taskMap[taskData.TaskID] = taskData
				} else {
					taskMd.targetToTask[oneTarget.TargetType] = make(map[int32]*TaskDetailData)
					taskMd.targetToTask[oneTarget.TargetType][taskData.TaskID] = taskData
				}
			}
			break
		}
	case ETaskStatusComplete:
		{
			delete(taskMd.acceptList, taskID)
			for _, oneTarget := range taskData.TargetList {
				if targetMap, ok := taskMd.targetToTask[oneTarget.TargetType]; ok {
					delete(targetMap, taskID)
					if len(targetMap) < 1 {
						delete(taskMd.targetToTask, oneTarget.TargetType)
					}
				}
				taskMd.completeList[taskID] = taskData
			}
			break
		}
	case ETaskStatusFinish:
		{
			delete(taskMd.acceptList, taskID)
			delete(taskMd.completeList, taskID)
			for _, oneTarget := range taskData.TargetList {
				if targetMap, ok := taskMd.targetToTask[oneTarget.TargetType]; ok {
					delete(targetMap, taskID)
					if len(targetMap) < 1 {
						delete(taskMd.targetToTask, oneTarget.TargetType)
					}
				}
			}
			break
		}
	default:
	}

	return taskData
}
func (taskMd *TaskModuleBase) GetTaskByTarget(targetType int32) []*TaskDetailData {
	var taskList []*TaskDetailData = make([]*TaskDetailData, 0)
	targetMap, ok := taskMd.targetToTask[targetType]
	if ok {
		for _, v := range targetMap {
			taskList = append(taskList, v)
		}
	}
	return taskList
}

// func (taskMd *TaskModuleBase) RecalcTaskProgress(targetCfg *configdata.TaskTarget, optEnv *TaskOptEnv) int64 {
// 	return 0
// }

// 任务事件触发
func (taskMd *TaskModuleBase) OnTaskEvent(targetType int32, num int64, param1 int64, parma2 int64, optEnv *TaskOptEnv, extraParams ...int64) {
	completeTaskList := make([]*TaskDetailData, 0)
	doingTasks := taskMd.GetTaskByTarget(targetType)
	for _, taskData := range doingTasks {
		taskID := taskData.TaskID
		taskCfg := configdata.Instance().GetTaskCfg(int32(taskID))
		if taskCfg == nil {
			continue
		}
		// // 检查是否满足更新条件
		// if !(p.CheckUpdate(player, task.ETaskModuleTest, taskCfg, param1, parma2, optEnv, extraParams...)) {
		// 	continue
		// }
		if taskData.Status != ETaskStatusAccept {
			// 保护一下,理论上,这里应该只出现当前接取的任务
			logrus.Warnf("check task module logic, task_id:%v, status:%v", taskData.TaskID, taskData.Status)
			continue
		}
		for i, targetData := range taskData.TargetList {
			if targetData.TargetType != targetType {
				continue
			}
			newNum := targetData.Num
			// 获取更新逻辑类型
			logicType := GetProgressLogicType(targetData.TargetType, false)
			if logicType == ETaskProgressLogicAdd {
				newNum = newNum + num
			} else if logicType == ETaskProgressLogicUpd {
				if newNum == num {
					// 没有变化
					break
				}
				newNum = num
			} else {
				// 重新取数据计算
				n := taskMd.taskOwner.RecalcTaskProgress(taskCfg.Targets[i], optEnv)
				if n == newNum {
					// 没有变化
					break
				}
				newNum = n
			}
			targetData.Num = newNum
			// 检测任务是否完成
			if taskData.IsComplete(taskCfg) {
				completeTaskList = append(completeTaskList, taskData)
			}
			break
		}
	}
	loghlp.Debugf("completeTaskList:%d", len(completeTaskList))
}
