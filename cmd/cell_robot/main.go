/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-20 13:55:30
 * @LastEditTime: 2022-09-23 17:28:35
 * @FilePath: \trcell\cmd\cell_robot\main.go
 */
package main

import (
	"fmt"
	"sync"
	"time"
	"trcell/app/cellrobot"
	"trcell/app/cellrobot/robotai"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(
		&logrus.TextFormatter{
			ForceColors:     false, // 这里必须设置false,否则时间格式显示不正常
			TimestampFormat: "2006-01-02 15:04:05",
			// TimestampFormat: time.RFC3339,
		},
	)
	logrus.SetLevel(logrus.DebugLevel)
	pflag.String("hostAddr", "localhost:15000", "account addr")
	pflag.String("nameId", "0", "name id")
	pflag.String("robotNum", "0", "robot number")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	hostAddr := viper.GetString("hostAddr")
	nameId := cast.ToInt32(viper.GetString("nameId"))
	if nameId == 0 {
		fmt.Print("please enter nameId:")
		fmt.Scanln(&nameId)
		// scanner := bufio.NewScanner(os.Stdin)
	}
	robotNum := cast.ToInt64(viper.GetString("robotNum"))
	if robotNum == 0 {
		fmt.Print("please enter robot num:")
		fmt.Scanln(&robotNum)
		// scanner := bufio.NewScanner(os.Stdin)
	}
	logrus.Debugf("hostAddr=%s", viper.GetString("hostAddr"))
	logrus.Debugf("robotNum=%s", viper.GetString("robotNum"))
	logrus.Info("hello cell robot")
	var robotNamePrifix string = "robot"
	var sg sync.WaitGroup
	for i := 0; i < int(robotNum); i++ {
		sg.Add(1)
		robotName := fmt.Sprintf("%s%d", robotNamePrifix, nameId+int32(i))
		go func() {
			testRobot := cellrobot.NewAiRobot(robotName, robotai.ERobotAiNone)
			if !testRobot.GetServerList(hostAddr) {
				sg.Done()
				return
			}
			if !testRobot.LoginAccount(hostAddr) {
				testRobot.RegisterAccount(hostAddr)
				if testRobot.LoginAccount(hostAddr) {
					if testRobot.ConnectServer(testRobot.GateAddr) == nil {
						testRobot.LoginGame()
						testRobot.Run()
					}
				}
			} else {
				if testRobot.ConnectServer(testRobot.GateAddr) == nil {
					testRobot.LoginGame()
					testRobot.Run()
				}
			}
			sg.Done()
		}()
		time.Sleep(time.Millisecond * 500)
	}
	sg.Wait()
}
