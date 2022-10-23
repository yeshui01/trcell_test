/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-10-11 14:56:23
 * @LastEditTime: 2022-10-11 14:57:31
 * @FilePath: \trcell\cmd\cellserv_social\main.go
 */
package main

import (
	"os"
	"os/signal"
	"trcell/app/servsocial"
	"trcell/pkg/loghlp"
	"trcell/pkg/trframe"
	"trcell/pkg/trframe/trnode"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	loghlp.Infof("cellserv_social main")
	logrus.SetReportCaller(true)
	loghlp.ActiveConsoleLog()
	loghlp.SetConsoleLogLevel(logrus.DebugLevel)
	loghlp.Debugf("trframe cellserv_social start")
	pflag.String("configPath", "./", "config file path")
	pflag.String("index", "0", "server index")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	loghlp.Debugf("configPath=%s", viper.GetString("configPath"))
	loghlp.Debugf("serverIndex=%s", viper.GetString("index"))
	cfgPath := viper.GetString("configPath")
	servIdx := cast.ToInt32(viper.GetString("index"))

	stopSig := make(chan os.Signal)
	signal.Notify(stopSig, os.Interrupt)
	go func() {
		<-stopSig
		trframe.Stop()
	}()
	trframe.Init(cfgPath, trnode.ETRNodeTypeCellSocial, servIdx)
	servSocial := servsocial.NewCellServSocial()
	trframe.RegisterUserFrameRun(func(curTimeMs int64) {
		servSocial.FrameRun(curTimeMs)
	})
	loghlp.ActiveFileLogReportCaller(true)
	trframe.Start()
}
