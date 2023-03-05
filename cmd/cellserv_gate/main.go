/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 15:01:44
 * @LastEditTime: 2022-09-19 15:53:40
 * @FilePath: \trcell\cmd\cellserv_gate\main.go
 */
package main

import (
	"net"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"trcell/app/servgate"
	"trcell/pkg/loghlp"
	"trcell/pkg/protocol"
	"trcell/pkg/tcpserver"
	"trcell/pkg/trframe"
	"trcell/pkg/trframe/trnode"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	loghlp.Infof("cellserv_gate main")
	logrus.SetReportCaller(true)
	loghlp.ActiveConsoleLog()
	loghlp.SetConsoleLogLevel(logrus.DebugLevel)
	loghlp.Debugf("trframe cellserv_gate start")
	pflag.String("configPath", "./", "config file path")
	pflag.String("index", "0", "server index")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	loghlp.Debugf("configPath=%s", viper.GetString("configPath"))
	loghlp.Debugf("serverIndex=%s", viper.GetString("index"))
	cfgPath := viper.GetString("configPath")
	servIdx := cast.ToInt32(viper.GetString("index"))
	trframe.Init(cfgPath, trnode.ETRNodeTypeCellGate, servIdx)
	tcpServe := tcpserver.NewTcpServer()
	stopSig := make(chan os.Signal)
	stopCh := make(chan bool)
	signal.Notify(stopSig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		if errP := recover(); errP != nil {
			loghlp.Errorf("gate panic:", errP)
			return
		}
	}()
	go func() {
		<-stopSig
		stopCh <- true
	}()
	tcpServe.SetupConnCallback(func(tcpConn net.Conn, err error) {
		trframe.PostUserCommand(protocol.CellCmdClassTcpsocket, protocol.CmdTypeTcpsocketConnect, tcpConn)
	})
	go func() {
		tcpServe.Run(trframe.GetFrameConfig().CellGateCfgs[servIdx].TcpListenAddr, stopCh)
		trframe.Stop()
	}()
	signal.Notify(stopSig, os.Interrupt)
	servGate := servgate.NewCellServGate()
	trframe.RegUserCommandHandler(func(frameCmd *trframe.TRFrameCommand) {
		loghlp.Infof("recv usercmd(%d_%d),type:%+v",
			frameCmd.UserCmd.GetCmdClass(),
			frameCmd.UserCmd.GetCmdType(),
			reflect.TypeOf(frameCmd.UserCmd.GetCmdData()).String(),
		)
		// 处理命令
		servGate.HandleCommand(frameCmd)
	})
	trframe.RegisterUserFrameRun(func(curTimeMs int64) {
		servGate.FrameRun(curTimeMs)
	})
	loghlp.ActiveFileLogReportCaller(true)
	trframe.Start()
}
