/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-28 10:33:38
 * @LastEditTime: 2022-09-28 15:56:22
 * @FilePath: \trcell\cmd\cellserv_global\main.go
 */
package main

import (
	"context"
	"fmt"
	"net"
	"time"
	"trcell/app/servglobal"
	"trcell/app/servglobal/globalconfig"
	"trcell/pkg/loghlp"
	"trcell/pkg/pb/pbrpc"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	rpcaddr = "localhost:50051"
)

type GlobalDataServer struct {
	pbrpc.UnimplementedGlobalDataServer
}

func (globalServ *GlobalDataServer) EchoTest(ctx context.Context, req *pbrpc.EchoReq) (*pbrpc.EchoRep, error) {
	loghlp.Debugf("serv recv EchoTest:%s", req.SendText)
	return &pbrpc.EchoRep{
		SendText: req.SendText,
	}, nil
}

// 客户端定义
func TestGRPCClient() {
	time.Sleep(time.Second * 5) // 5秒钟之后开始测试

	conn, err := grpc.Dial(rpcaddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		loghlp.Errorf("err:%s", err.Error())
		return
	}
	defer conn.Close()
	var globalClient pbrpc.GlobalDataClient = pbrpc.NewGlobalDataClient(conn)
	echoRep, repErr := globalClient.EchoTest(context.Background(), &pbrpc.EchoReq{
		SendText: "hello rpc world",
	})
	if repErr != nil {
		loghlp.Errorf("repError:%s", repErr.Error())
		return
	}
	loghlp.Infof("echo succ:%+v", echoRep)
}
func TestGRPCServ() {
	lisn, err := net.Listen("tcp", rpcaddr)
	if err != nil {
		loghlp.Errorf("listen error:%s", err.Error())
		return
	}
	s := grpc.NewServer()
	pbrpc.RegisterGlobalDataServer(s, &GlobalDataServer{})
	go TestGRPCClient()
	if err := s.Serve(lisn); err != nil {
		loghlp.Errorf("rpcserv err:%s", err.Error())
	}
}

//
func main() {
	loghlp.SetConsoleLogLevel(logrus.DebugLevel)
	loghlp.ActiveConsoleLog()
	loghlp.Debugf("hello cellserv_global")
	pflag.String("configPath", "./", "config file path")
	pflag.String("index", "0", "server index")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	loghlp.Debugf("configPath=%s", viper.GetString("configPath"))
	loghlp.Debugf("serverIndex=%s", viper.GetString("index"))
	configPath := viper.GetString("configPath")
	servIdx := cast.ToInt32(viper.GetString("index"))
	loghlp.Infof("serv global, cfgPath:%s, servIdx:%d", configPath, servIdx)
	// 加载配置
	appConfig := globalconfig.NewServAppConfig()
	err := globalconfig.ReadGlobalConfigFromFile(configPath, appConfig)
	if err != nil {
		panic(fmt.Sprintf("load global config error:%s", err.Error()))
	} else {
		loghlp.Debugf("load global config succ:%v", appConfig)
	}

	servGlobal := servglobal.NewCellServGlobal(appConfig)
	servGlobal.Run()
	//TestGRPCServ()
}
