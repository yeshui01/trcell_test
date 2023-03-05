package main

import (
	"fmt"
	"time"
	"trcell/app/servloginqueue"
	"trcell/app/servloginqueue/loginqueueclient"
	"trcell/app/servloginqueue/loginqueueconfig"
	"trcell/pkg/loghlp"
	"trcell/pkg/pb/pbrpc"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var queueClientsPool *loginqueueclient.LoginQueueClientPool

func TestLoginQueueClient(servAddr string, accID int32) {
	queueClient := queueClientsPool.GetClient()
	getSeqReq := &pbrpc.LoginSeqNoReq{
		AccountName: fmt.Sprintf("acc%d", accID),
	}

	loginSeqNoRep, err := queueClient.CallGetLoginSeqNo(getSeqReq)
	if err != nil {
		loghlp.Errorf("loginSeqNoRep error:%s", err.Error())
		queueClientsPool.KeepReuse(queueClient)
		return
	}
	loghlp.Infof("user(%s) loginSeqNoRep:%+v", getSeqReq.AccountName, loginSeqNoRep)

	// 获取已经完成的seq
	finishSeqReq := &pbrpc.QueryCurLoginFinishNoReq{}
	finishSeqRep, errF := queueClient.CallGetFinishSeqNo(finishSeqReq)
	if errF != nil {
		loghlp.Errorf("finishSeqRep error:%s", err.Error())
		queueClientsPool.KeepReuse(queueClient)
		return
	}
	if loginSeqNoRep.SeqNo <= finishSeqRep.CurFinishNo+100 {
		// Do Something,可以登录
		incReq := &pbrpc.IncrementLoginFinishReq{
			AccountName: getSeqReq.AccountName,
		}
		queueClient.CallIncrementLoginFinish(incReq)
	}
	// 说明需要排队等候
	queueClientsPool.KeepReuse(queueClient)
}

func main() {
	loghlp.SetConsoleLogLevel(logrus.DebugLevel)
	loghlp.ActiveConsoleLog()
	loghlp.Debugf("hello cellserv_loginqueue")
	pflag.String("configPath", "./", "config file path")
	pflag.String("index", "0", "server index")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	loghlp.Debugf("configPath=%s", viper.GetString("configPath"))
	loghlp.Debugf("serverIndex=%s", viper.GetString("index"))
	configPath := viper.GetString("configPath")
	servIdx := cast.ToInt32(viper.GetString("index"))
	loghlp.Infof("serv loginqueue, cfgPath:%s, servIdx:%d", configPath, servIdx)
	// 加载配置
	appConfig := loginqueueconfig.NewServAppConfig()
	err := loginqueueconfig.ReadLoginQueueConfigFromFile(configPath, appConfig)
	if err != nil {
		panic(fmt.Sprintf("load loginqueue config error:%s", err.Error()))
	} else {
		loghlp.Debugf("load loginqueue config succ:%v", appConfig)
	}

	servLoginQueue := servloginqueue.NewCellServLoginQueue(appConfig)
	queueClientsPool = loginqueueclient.NewLoginQueueClientPool(1, appConfig.RpcAddr)
	for i := 0; i < 3; i++ {
		go func(rpcAddr string, accID int32) {
			time.Sleep(time.Second * 2)
			TestLoginQueueClient(rpcAddr, accID)
		}(appConfig.RpcAddr, int32(i)+1)
	}
	for i := 3; i < 6; i++ {
		go func(rpcAddr string, accID int32) {
			time.Sleep(time.Second * 4)
			TestLoginQueueClient(rpcAddr, accID)
		}(appConfig.RpcAddr, int32(i)+1)
	}
	servLoginQueue.Run()
}
