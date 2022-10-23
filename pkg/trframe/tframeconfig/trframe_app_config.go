/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-16 14:18:14
 * @LastEditTime: 2022-10-14 13:42:56
 * @FilePath: \trcell\pkg\trframe\tframeconfig\trframe_app_config.go
 */
package tframeconfig

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

// 账号配置
type AccountCfg struct {
	ListenAddr string    `yaml:"listenAddr"`
	LogLevel   int32     `yaml:"logLevel"`
	LogPath    string    `yaml:"logPath"`
	AccountDb  *MysqlCfg `yaml:"accountDb"`
}

// cells
// cellroot
type CellRootCfg struct {
	ListenAddr string `yaml:"listenAddr"`
	LogLevel   int32  `yaml:"logLevel"`
	ListenMode string `yaml:"listenMode"`
	LogPath    string `yaml:"logPath"`
	CsvPath    string `yaml:"csvPath"`
}

// cellgate
type CellGateCfg struct {
	ListenAddr    string `yaml:"listenAddr"`
	LogLevel      int32  `yaml:"logLevel"`
	ListenMode    string `yaml:"listenMode"`
	WsListenAddr  string `yaml:"wsListenAddr"`
	TcpListenAddr string `yaml:"tcpListenAddr"`
	LogPath       string `yaml:"logPath"`
}

// celldata
type CellDataCfg struct {
	ListenAddr   string   `yaml:"listenAddr"`
	LogLevel     int32    `yaml:"logLevel"`
	ListenMode   string   `yaml:"listenMode"`
	WsListenAddr string   `yaml:"wsListenAddr"`
	LogPath      string   `yaml:"logPath"`
	GameDb       MysqlCfg `yaml:"gameDb"`
}

// cellcenter
type CellCenterCfg struct {
	ListenAddr string `yaml:"listenAddr"`
	LogLevel   int32  `yaml:"logLevel"`
	ListenMode string `yaml:"listenMode"`
	LogPath    string `yaml:"logPath"`
	CsvPah     string `yaml:"csvPath"`
}

// celllogic
type CellLogicCfg struct {
	ListenAddr string `yaml:"listenAddr"`
	LogLevel   int32  `yaml:"logLevel"`
	ListenMode string `yaml:"listenMode"`
	LogPath    string `yaml:"logPath"`
	CsvPah     string `yaml:"csvPath"`
}

// cellgame
type CellGameCfg struct {
	ListenAddr string `yaml:"listenAddr"`
	LogLevel   int32  `yaml:"logLevel"`
	ListenMode string `yaml:"listenMode"`
	LogPath    string `yaml:"logPath"`
	CsvPah     string `yaml:"csvPath"`
}

// cellview
type CellViewCfg struct {
	ListenAddr string `yaml:"listenAddr"`
	LogLevel   int32  `yaml:"logLevel"`
	ListenMode string `yaml:"listenMode"`
	LogPath    string `yaml:"logPath"`
	CsvPah     string `yaml:"csvPath"`
}

// cellLog
type CellLogCfg struct {
	ListenAddr string `yaml:"listenAddr"`
	LogLevel   int32  `yaml:"logLevel"`
	ListenMode string `yaml:"listenMode"`
	LogPath    string `yaml:"logPath"`
}

// cellsocial
type CellSocialCfg struct {
	ListenAddr string `yaml:"listenAddr"`
	LogLevel   int32  `yaml:"logLevel"`
	ListenMode string `yaml:"listenMode"`
	LogPath    string `yaml:"logPath"`
}

// celltrans
type CellTransCfg struct {
	ListenAddr string `yaml:"listenAddr"`
	LogLevel   int32  `yaml:"logLevel"`
	ListenMode string `yaml:"listenMode"`
	LogPath    string `yaml:"logPath"`
}
type FrameConfig struct {
	AccountCfgs []*AccountCfg `yaml:"accountCfgs"`
	ZoneID      int32         `yaml:"zoneID"`
	//
	CellRootCfgs   []*CellRootCfg   `yaml:"cellRootCfgs"`
	CellGateCfgs   []*CellGateCfg   `yaml:"cellGateCfgs"`
	CellDataCfgs   []*CellDataCfg   `yaml:"cellDataCfgs"`
	CellCenterCfgs []*CellCenterCfg `yaml:"cellCenterCfgs"`
	CellLogicCfgs  []*CellLogicCfg  `yaml:"cellLogicCfgs"`
	CellGameCfgs   []*CellGameCfg   `yaml:"cellGameCfgs"`
	CellViewCfgs   []*CellViewCfg   `yaml:"cellViewCfgs"`
	CellLogCfgs    []*CellLogCfg    `yaml:"cellLogCfgs"`
	CellSocialCfgs []*CellSocialCfg `yaml:"cellSocialCfgs"`
	CellTransCfgs  []*CellTransCfg  `yaml:"cellTransCfgs"`
}

func NewFrameConfig() *FrameConfig {
	return &FrameConfig{}
}

//从本地文件读取配置
func ReadFrameConfigFromFile(configFilePath string, defaultSetting *FrameConfig) error {
	fullPath := path.Join(configFilePath, "trframe.yaml")
	viper.SetConfigFile(fullPath)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return viper.Unmarshal(defaultSetting)
}
