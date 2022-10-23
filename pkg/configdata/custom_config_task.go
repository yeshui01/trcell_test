/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-28 16:05:47
 * @LastEditTime: 2022-09-28 17:16:16
 * @FilePath: \trcell\pkg\configdata\custom_config_task.go
 */
package configdata

type TaskTarget struct {
	TargetType int32
	Param1     int32
	Param2     int32
	Param3     int32
	TotalNum   int64
}

type TaskCfg struct {
	TaskID  int32
	Targets []*TaskTarget
}

func NewTaskCfg(taskID int32) *TaskCfg {
	return &TaskCfg{
		TaskID:  taskID,
		Targets: make([]*TaskTarget, 1), // 默认就一个目标
	}
}
func (cfg *ConfigData) GetTaskCfg(taskID int32) *TaskCfg {
	if taskCfg, ok := cfg.taskConfig[taskID]; ok {
		return taskCfg
	}
	return nil
}
