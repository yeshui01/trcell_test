/**
 * @author: [mknight]
 * @email : [824338670@qq.com]
 * @create:	2023-02-02 09:52:38
 * @modify:	2023-02-02 09:52:38
 * @desc  :	[description]
 */
package loginqueueconfig

import (
	"fmt"
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

func (cfg *MysqlCfg) GenDsnStr() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Pswd,
		cfg.Host,
		cfg.Port,
		cfg.DbName)

	return dsn
}

type ServAppConfig struct {
	// GlobalDb *MysqlCfg `yaml:"globalDb"`
	RpcAddr string `yaml:"rpcAddr"`
	ZoneID  int32  `yaml:"zoneID"`
}

func NewServAppConfig() *ServAppConfig {
	return &ServAppConfig{}
}

//从本地文件读取配置
func ReadLoginQueueConfigFromFile(configFilePath string, defaultSetting *ServAppConfig) error {
	fullPath := path.Join(configFilePath, "cell_loginqueue.yaml")
	viper.SetConfigFile(fullPath)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return viper.Unmarshal(defaultSetting)
}
