/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-16 14:18:13
 * @LastEditTime: 2022-09-19 16:55:51
 * @FilePath: \trcell\pkg\appconfig\appconfig_instance.go
 */
package appconfig

import (
	"path"

	"github.com/spf13/viper"
)

var accountCfg *AccountConfig

func Instance() *AccountConfig {
	if accountCfg == nil {
		accountCfg = NewAccountCfg()
	}
	return accountCfg
}

//从本地文件读取配置
func (cfg *AccountConfig) Load(configFilePath string) error {
	fullPath := path.Join(configFilePath, "cell_account.yaml")
	viper.SetConfigFile(fullPath)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = viper.Unmarshal(cfg)
	return err
}
