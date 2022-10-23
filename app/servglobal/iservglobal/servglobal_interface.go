/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-28 14:11:27
 * @LastEditTime: 2022-10-17 10:18:36
 * @FilePath: \trcell\app\servglobal\iservglobal\servglobal_interface.go
 */
package iservglobal

import (
	"trcell/pkg/utils"

	"gorm.io/gorm"
)

type IServGlobal interface {
	GetGlobalDB() *gorm.DB
	GetIDGenerator() *utils.IDGenerator
}
