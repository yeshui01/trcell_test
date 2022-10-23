/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 11:08:44
 * @LastEditTime: 2022-09-30 09:59:24
 * @FilePath: \trcell\app\servdata\cellserv_data.go
 */
package servdata

import (
	"fmt"
	"time"
	"trcell/app/servdata/servdatahandler"
	"trcell/app/servdata/servdatamain"
	"trcell/pkg/loghlp"
	"trcell/pkg/trframe"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type CellServData struct {
	servDataGlobal *servdatamain.ServDataGlobal
}

func NewCellServData() *CellServData {
	s := &CellServData{
		servDataGlobal: servdatamain.NewServDataGlobal(),
	}
	s.RegisterMsgHandler()
	servdatahandler.InitServDataObj(s)
	return s
}

// // 打开本地数据库
// func (serv *CellServData) OpenLocalDB(dbFile string) bool {
// 	dbLogger := logger.New(
// 		logrus.New(), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
// 		logger.Config{
// 			SlowThreshold:             time.Second, // 慢 SQL 阈值
// 			LogLevel:                  logger.Info, // 日志级别
// 			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
// 			Colorful:                  false,       // 禁用彩色打印
// 		},
// 	)
// 	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{
// 		Logger: dbLogger,
// 	})
// 	if err != nil {
// 		panic("open dbfile faile")
// 	}
// 	sqlDB, _ := db.DB()
// 	sqlDB.SetMaxIdleConns(2)
// 	sqlDB.SetMaxOpenConns(2)
// 	sqlDB.SetConnMaxLifetime(time.Hour)
// 	serv.servDataGlobal.InitGameDB(db)
// 	return true
// }

// 打开Mysql数据库
func (serv *CellServData) OpenMysqlDB(connStr string) bool {
	dbLogger := logger.New(
		logrus.New(), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,       // 禁用彩色打印
		},
	)
	db, err := gorm.Open(mysql.Open(connStr), &gorm.Config{
		Logger: dbLogger,
	})
	if err != nil {
		loghlp.Errorf("open mysql db error:%s", err.Error())
		panic("open dbfile faile")
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetMaxOpenConns(2)
	sqlDB.SetConnMaxLifetime(time.Hour)
	serv.servDataGlobal.InitGameDB(db)
	return true
}

func (serv *CellServData) GetGameDB() *gorm.DB {
	return serv.servDataGlobal.GetGameDB()
}
func (serv *CellServData) GetDataGlobal() *servdatamain.ServDataGlobal {
	return serv.servDataGlobal
}
func (serv *CellServData) FrameRun(curTimeMs int64) {
	serv.servDataGlobal.FrameUpdate(curTimeMs)
}

func (serv *CellServData) CheckInitGameDB() {
	zoneID := trframe.GetFrameConfig().ZoneID
	var result int64 = 0
	var initID int64 = 1000000*int64(zoneID) + 1
	serv.GetGameDB().Raw("SELECT role_id FROM role_base LIMIT 1").Scan(&result)
	if result == 0 {
		// 初始化
		err := serv.GetGameDB().Exec(fmt.Sprintf("ALTER TABLE role_base AUTO_INCREMENT = %d;", initID)).Error
		if err != nil {
			loghlp.Errorf("init alter table role_base auto_increment fail:%s", err.Error())
			panic(fmt.Sprintf("init alter role_base auto_increment err:%s", err.Error()))
		} else {
			loghlp.Infof("init alter role_base auto_increment succ, initID:%d", initID)
		}
	}
}
