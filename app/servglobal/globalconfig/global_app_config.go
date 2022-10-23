/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-28 14:15:59
 * @LastEditTime: 2022-09-28 14:45:09
 * @FilePath: \trcell\app\servglobal\globalconfig\global_app_config.go
 */
package globalconfig

import (
	"path"

	"github.com/spf13/viper"
)

// mysql配置
type MysqlCfg struct {
	User   string `yaml:"user"`
	Host   string `yaml:"host"`
	Port   int32  `yaml:"port"`
	Pswd   string `yaml:"pswd"`
	DbName string `yaml:"dbName"`
}

type ServAppConfig struct {
	GlobalDb *MysqlCfg `yaml:"globalDb"`
	RpcAddr  string    `yaml:"rpcAddr"`
	ZoneID   int32     `yaml:"zoneID"`
}

func NewServAppConfig() *ServAppConfig {
	return &ServAppConfig{}
}

//从本地文件读取配置
func ReadGlobalConfigFromFile(configFilePath string, defaultSetting *ServAppConfig) error {
	fullPath := path.Join(configFilePath, "cell_global.yaml")
	viper.SetConfigFile(fullPath)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return viper.Unmarshal(defaultSetting)
}
