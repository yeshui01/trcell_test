/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-16 14:18:13
 * @LastEditTime: 2022-09-19 16:26:02
 * @FilePath: \trcell\app\account\accrouter\account_model.go
 */
package accrouter

type OrmUser struct {
	UserID       int64  `gorm:"column:user_id;primaryKey"`
	UserName     string `gorm:"column:user_name"`
	Pswd         string `gorm:"column:pswd"`
	RegisterTime int64  `gorm:"column:register_time"`
	Status       int32  `gorm:"column:status"`
	DataZone     int32  `gorm:"column:data_zone"`
}

func (tbuser *OrmUser) TableName() string {
	return "user"
}

type OrmServerList struct {
	ID         int32  `gorm:"column:id;primaryKey"`
	ServerName string `gorm:"column:server_name"`
	GateAddr   string `gorm:"column:gate_addr"`
	Recommend  int32  `gorm:"column:recommend"`
}

func (tb *OrmServerList) TableName() string {
	return "server_list"
}

type OrmCellNotice struct {
	ID      int32  `gorm:"column:id;primaryKey"`
	Content string `gorm:"column:content"`
	UpdTime int64  `gorm:"column:upd_time"`
}

func (tb *OrmCellNotice) TableName() string {
	return "cell_notice"
}
