package servloginqueue

import (
	"net"
	"trcell/app/servloginqueue/loginqueue"
	"trcell/app/servloginqueue/loginqueueconfig"
	"trcell/app/servloginqueue/loginqueuerpc"
	"trcell/pkg/loghlp"
	"trcell/pkg/pb/pbrpc"

	"google.golang.org/grpc"
)

type CellServLoginQueue struct {
	rpcServ    *loginqueuerpc.LoginQueueServer
	appConfig  *loginqueueconfig.ServAppConfig
	loginQueue *loginqueue.LoginQueue
}

func NewCellServLoginQueue(appCfg *loginqueueconfig.ServAppConfig) *CellServLoginQueue {
	serv := &CellServLoginQueue{
		appConfig:  appCfg,
		loginQueue: loginqueue.NewLoginQueue(),
	}
	serv.rpcServ = loginqueuerpc.NewLoginQueueServer(serv)
	return serv
}

func (serv *CellServLoginQueue) Run() {
	lisn, err := net.Listen("tcp", serv.appConfig.RpcAddr)
	if err != nil {
		loghlp.Errorf("rpclisten error:%s", err.Error())
		return
	}
	worerNumOpts := grpc.NumStreamWorkers(100)
	s := grpc.NewServer(worerNumOpts)
	pbrpc.RegisterLoginQueueBackendServer(s, serv.rpcServ)
	if err := s.Serve(lisn); err != nil {
		loghlp.Errorf("rpcserv error:%s", err.Error())
	}
}

func (serv *CellServLoginQueue) GetLoginQueue() *loginqueue.LoginQueue {
	return serv.loginQueue
}
