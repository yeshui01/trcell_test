/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-16 14:18:13
 * @LastEditTime: 2022-09-19 16:58:31
 * @FilePath: \trcell\pkg\appconfig\appconfig_room_cell.go
 */
package appconfig

import (
	"path"

	"github.com/spf13/viper"
)

// room cell config
// mysql配置
type MysqlCfg struct {
	User   string `yaml:"user"`
	Host   string `yaml:"host"`
	Port   int32  `yaml:"port"`
	Pswd   string `yaml:"pswd"`
	DbName string `yaml:"dbName"`
}

// 账号配置
type AccountCfg struct {
	ListenAddr string `yaml:"listenAddr"`
	LogLevel   int32  `yaml:"logLevel"`
	LogPath    string `yaml:"logPath"`
}

type AccountConfig struct {
	AccountCfgs []*AccountCfg `yaml:"accountCfgs"`
	ZoneID      int32         `yaml:"zoneID"`
	AccountDB   *MysqlCfg     `yaml:"accountDB"`
}

func NewAccountCfg() *AccountConfig {
	return &AccountConfig{}
}

//从本地文件读取配置
func ReadAccountConfigFromFile(configFilePath string, defaultSetting *AccountConfig) error {
	fullPath := path.Join(configFilePath, "cell_account.yaml")
	viper.SetConfigFile(fullPath)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return viper.Unmarshal(defaultSetting)
}
