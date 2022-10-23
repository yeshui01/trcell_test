/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-10-12 10:14:37
 * @LastEditTime: 2022-10-12 11:30:29
 * @FilePath: \trcell\pkg\jobworker\worker_hub.go
 */
package jobworker

import (
	"sync"
	"trcell/pkg/loghlp"
	"trcell/pkg/timeutil"
)

type workerHandler func(workerID int32, oneJob *OneJob)
type WorkerHub struct {
	maxWorkerNum int32
	workerList   []*JobWorker
	jobHandler   workerHandler
	workerWaits  sync.WaitGroup
	idGenerator  int64
	workderIdx   int
}

func NewWorkerHub(workerNumMax int32) *WorkerHub {
	if workerNumMax < 0 {
		workerNumMax = 0
	}
	wkHub := &WorkerHub{
		maxWorkerNum: workerNumMax,
		idGenerator:  0,
		workderIdx:   0,
	}
	return wkHub
}

func (wkHub *WorkerHub) SetupJobHandler(jobHandler workerHandler) {
	wkHub.jobHandler = jobHandler
}

func (wkHub *WorkerHub) Start() {
	if wkHub.maxWorkerNum < 1 {
		return
	}
	for i := 0; i < int(wkHub.maxWorkerNum); i++ {
		oneWorker := NewJobWorker(2)
		oneWorker.workerID = int32(i + 1)
		oneWorker.JobHandler = wkHub.jobHandler
		wkHub.workerWaits.Add(1)
		go func() {
			oneWorker.Run()
			wkHub.workerWaits.Done()
		}()
		wkHub.workerList = append(wkHub.workerList, oneWorker)
	}
}

func (wkHub *WorkerHub) PostJob(jobType int32, jobParam interface{}) {
	wkHub.idGenerator++
	if wkHub.idGenerator >= 9999999999999 {
		wkHub.idGenerator = 1
	}
	newJob := &OneJob{
		JobID:     wkHub.idGenerator,
		JobType:   jobType,
		JobParam:  jobParam,
		StartTime: timeutil.NowTime(),
	}
	if wkHub.maxWorkerNum < 1 {
		if wkHub.jobHandler != nil {
			wkHub.jobHandler(0, newJob)
		} else {
			loghlp.Errorf("wkHub.maxWorkerNum < 1 and wkHub.jobHandler == nil")
		}
	} else {
		idx := wkHub.workderIdx % len(wkHub.workerList)
		wkHub.workerList[idx].JobChs <- newJob
		wkHub.workderIdx++
	}
}
func (wkHub *WorkerHub) Stop() {
	for i := 0; i < int(wkHub.maxWorkerNum); i++ {
		wker := wkHub.workerList[i]
		close(wker.JobChs)
	}
	wkHub.workerWaits.Wait()
}
