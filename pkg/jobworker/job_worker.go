/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-10-12 10:13:54
 * @LastEditTime: 2022-10-12 11:30:23
 * @FilePath: \trcell\pkg\jobworker\job_worker.go
 */
package jobworker

import "trcell/pkg/loghlp"

type OneJob struct {
	JobID     int64
	JobType   int32
	JobParam  interface{}
	StartTime int64
}

type JobWorker struct {
	workerID   int32
	JobChs     chan *OneJob
	JobHandler workerHandler
}

func NewJobWorker(chSize int32) *JobWorker {
	if chSize < 1 {
		chSize = 1
	}
	return &JobWorker{
		JobChs: make(chan *OneJob, chSize),
	}
}

func (wk *JobWorker) Run() {
	for toDoJob := range wk.JobChs {
		if wk.JobHandler != nil {
			wk.JobHandler(wk.workerID, toDoJob)
		}
	}
	loghlp.Infof("job worker exist")
}
