/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-16 14:18:13
 * @LastEditTime: 2022-09-28 16:12:20
 * @FilePath: \trcell\pkg\configdata\csvconfig.go
 */
package configdata

var configData *ConfigData

type ConfigData struct {
	csvModules *CSVModuleMgr
	// custom config
	typeDrawWordsConfig map[int32]*DrawTypeWordsCfg
	taskConfig          map[int32]*TaskCfg
}

func newConfigData(csvPath string) *ConfigData {
	return &ConfigData{
		csvModules:          NewCSVModuleMgr(csvPath),
		typeDrawWordsConfig: make(map[int32]*DrawTypeWordsCfg),
		taskConfig:          make(map[int32]*TaskCfg),
	}
}

func (cfg *ConfigData) LoadConfig() {
	cfg.csvModules.LoadAll()
	// custom
	// task
	cfg.InitTaskCfg()
	// draw guess
	{
		wordsList := cfg.GetDrawGuessCfgList()
		for _, v := range wordsList {
			if typeList, ok := cfg.typeDrawWordsConfig[v.WordType]; ok {
				cfg.typeDrawWordsConfig[v.WordType].WordsList = append(typeList.WordsList, v)
			} else {
				typeList := &DrawTypeWordsCfg{
					WordsList: nil,
				}
				typeList.WordsList = append(typeList.WordsList, v)
				cfg.typeDrawWordsConfig[v.WordType] = typeList
			}
		}
	}
}

func InitConfigData(csvPath string) {
	configData = newConfigData(csvPath)
	configData.LoadConfig()
}

func Instance() *ConfigData {
	return configData
}
