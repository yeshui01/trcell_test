/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-10-19 14:10:18
 * @LastEditTime: 2022-10-19 15:07:19
 * @FilePath: \trcell\cmd\testapp\testworker\main.go
 */
package main

import (
	"fmt"
	"time"
	"trcell/app/servglobal/globalclient"
	"trcell/pkg/jobworker"
	"trcell/pkg/loghlp"
	"trcell/pkg/pb/pbrpc"
	"trcell/pkg/timeutil"

	"github.com/sirupsen/logrus"
)

func main() {
	loghlp.ActiveConsoleLog()
	loghlp.SetConsoleLogLevel(logrus.DebugLevel)
	loghlp.Infof("hello test worker")
	workerHub := jobworker.NewWorkerHub(0)
	workerHub.SetupJobHandler(handleJob)
	workerHub.Start()
	for i := 0; i < 10; i++ {
		workerHub.PostJob(1, fmt.Sprintf("this is job%d", i+1))
	}
	time.Sleep(time.Second * 1)
	workerHub.Stop()

	// 测试生成uid的耗时

	{
		globalClient := globalclient.NewServGlobalClient("localhost:50071")
		err := globalClient.Connect("localhost:50071")
		if err != nil {
			loghlp.Errorf("connect globalServ error:%s", err.Error())
		} else {
			beginUidTime := timeutil.NowTimeMs()
			loghlp.Infof("begin test gen uid time:%d", beginUidTime)
			uidNum := int32(1000)
			repUid, errUid := globalClient.CallGenUid(&pbrpc.GenUIDReq{
				Num: uidNum,
			})
			if errUid != nil {
				loghlp.Errorf("genUid error:%s", errUid.Error())
			} else {
				loghlp.Infof("gen uid succ, uidNum:%d", len(repUid.UIDs))
			}
			endUidTime := timeutil.NowTimeMs()
			loghlp.Infof("end test gen uid time:%d, elapse(%d) ms, total uidNum:%d", endUidTime, (endUidTime - beginUidTime), uidNum)

			// second times

			beginUidTime = timeutil.NowTimeMs()
			repUid, errUid = globalClient.CallGenUid(&pbrpc.GenUIDReq{
				Num: uidNum,
			})
			endUidTime = timeutil.NowTimeMs()
			if errUid != nil {
				loghlp.Errorf("genUid error:%s", errUid.Error())
			} else {
				loghlp.Infof("gen uid succ2, uidNum:%d", len(repUid.UIDs))
				loghlp.Infof("end2 test gen uid elapse(%d) ms, total uidNum:%d", (endUidTime - beginUidTime), len(repUid.UIDs))
			}
		}
	}

	time.Sleep(time.Second * 1)
}

func handleJob(workerID int32, oneJob *jobworker.OneJob) {
	loghlp.Infof("handle job, oneJob.JobID:%d, jobType:%d", oneJob.JobID, oneJob.JobType)
	// 这里测试,echo
	echoStr := oneJob.JobParam.(string)
	loghlp.Infof("worker(%d) echoJobParam:%s", workerID, echoStr)
	// 生成一个rpcuid
	globalClient := globalclient.NewServGlobalClient("localhost:50071")
	err := globalClient.Connect("localhost:50071")
	if err != nil {
		loghlp.Errorf("connect globalServ error:%s", err.Error())
		return
	}
	repUid, errUid := globalClient.CallGenUid(&pbrpc.GenUIDReq{
		Num: 2,
	})
	if errUid != nil {
		loghlp.Errorf("genUid error:%s", errUid.Error())
		return
	}
	loghlp.Infof("genuid succ:%+v", repUid)
}
