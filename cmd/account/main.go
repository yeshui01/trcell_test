/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-16 14:18:13
 * @LastEditTime: 2022-09-19 17:08:01
 * @FilePath: \trcell\cmd\account\main.go
 */
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	"trcell/app/account"
	"trcell/pkg/appconfig"
	"trcell/pkg/loghlp"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	loghlp.ActiveConsoleLog()
	loghlp.Debugf("cellserv_account start")
	pflag.String("configPath", "./", "config file path")
	pflag.String("index", "100", "server index")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	loghlp.Debugf("configPath=%s", viper.GetString("configPath"))
	loghlp.Debugf("serverIndex=%s", viper.GetString("index"))
	// go func() {
	// 	log.Println(http.ListenAndServe("127.0.0.1:6060", nil))
	// }()
	serverIndex := cast.ToInt(viper.GetString("index"))
	accountCfg := appconfig.NewAccountCfg()
	errConfig := appconfig.ReadAccountConfigFromFile(viper.GetString("configPath"), accountCfg)
	if errConfig != nil {
		loghlp.Errorf("read config error:%s", errConfig.Error())
		return
	}
	jvData, err := json.Marshal(accountCfg)
	if err != nil {
		loghlp.Debugf("err:%s", err.Error())
	} else {
		loghlp.Debugf("configJv:%s", string(jvData))
	}
	// runtime.SetMutexProfileFraction(1)
	// runtime.SetBlockProfileRate(1)
	loghlp.Debugf("accountCfg:%+v", accountCfg)
	loghlp.SetConsoleLogLevel(logrus.Level(accountCfg.AccountCfgs[serverIndex].LogLevel))
	accApp := account.NewAccount(accountCfg)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", accountCfg.AccountDB.User, accountCfg.AccountDB.Pswd, accountCfg.AccountDB.Host, accountCfg.AccountDB.Port, accountCfg.AccountDB.DbName)
	loghlp.Infof("open mysql:%s", dsn)
	accApp.OpenMysqlDB(dsn)

	stopCh := make(chan os.Signal)
	signal.Notify(stopCh, os.Interrupt)
	go func() {
		<-stopCh
		accApp.Stop()
	}()
	accApp.Run(accountCfg.AccountCfgs[serverIndex].ListenAddr, 0)
}
