/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 11:14:00
 * @LastEditTime: 2022-09-19 11:19:08
 * @FilePath: \trcell\app\servdata\iservdata\servdata_interface.go
 */
package iservdata

import (
	"trcell/app/servdata/servdatamain"

	"gorm.io/gorm"
)

type IServData interface {
	GetGameDB() *gorm.DB
	GetDataGlobal() *servdatamain.ServDataGlobal
}
