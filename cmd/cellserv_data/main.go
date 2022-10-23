/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 10:23:04
 * @LastEditTime: 2022-09-23 16:04:33
 * @FilePath: \trcell\cmd\cellserv_data\main.go
 */
package main

import (
	"fmt"
	"os"
	"os/signal"
	"trcell/app/servdata"
	"trcell/pkg/loghlp"
	"trcell/pkg/trframe"
	"trcell/pkg/trframe/trnode"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	loghlp.Infof("cellserv_data main")
	logrus.SetReportCaller(true)
	loghlp.ActiveConsoleLog()
	loghlp.SetConsoleLogLevel(logrus.DebugLevel)
	loghlp.Debugf("trframe cellserv_root start")
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
	trframe.Init(cfgPath, trnode.ETRNodeTypeCellData, servIdx)
	servData := servdata.NewCellServData()
	dbCfg := trframe.GetFrameConfig().CellDataCfgs[servIdx]
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbCfg.GameDb.User, dbCfg.GameDb.Pswd, dbCfg.GameDb.Host, dbCfg.GameDb.Port, dbCfg.GameDb.DbName)
	loghlp.Infof("open gamemysql:%s", dsn)
	servData.OpenMysqlDB(dsn)
	servData.CheckInitGameDB()
	// 检查role_base表
	trframe.RegisterUserFrameRun(func(curTimeMs int64) {
		servData.FrameRun(curTimeMs)
	})
	loghlp.ActiveFileLogReportCaller(true)
	trframe.Start()
}
