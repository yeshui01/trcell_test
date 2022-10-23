/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-28 13:31:32
 * @LastEditTime: 2022-10-17 10:19:15
 * @FilePath: \trcell\app\servglobal\cellserv_global.go
 */
package servglobal

import (
	"fmt"
	"net"
	"time"
	"trcell/app/servglobal/globalconfig"
	"trcell/app/servglobal/globalrpc"
	"trcell/pkg/loghlp"
	"trcell/pkg/pb/pbrpc"
	"trcell/pkg/utils"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type CellServGlobal struct {
	globalDB   *gorm.DB
	appConfig  *globalconfig.ServAppConfig
	rpcServ    *globalrpc.GlobalDataServer
	idGenrator *utils.IDGenerator
}

func NewCellServGlobal(appCfg *globalconfig.ServAppConfig) *CellServGlobal {
	serv := &CellServGlobal{
		appConfig: appCfg,
	}
	var err error
	serv.rpcServ = globalrpc.NewGlobalDataServer(serv)
	serv.idGenrator, err = utils.NewIDGenerator(int64(appCfg.ZoneID))
	if err != nil {
		loghlp.Errorf("newCellServGlobalError:%s", err.Error())
		return nil
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		appCfg.GlobalDb.User,
		appCfg.GlobalDb.Pswd,
		appCfg.GlobalDb.Host,
		appCfg.GlobalDb.Port,
		appCfg.GlobalDb.DbName)
	loghlp.Infof("open globalmysql:%s", dsn)
	serv.OpenMysqlDB(dsn)
	return serv
}

func (serv *CellServGlobal) GetGlobalDB() *gorm.DB {
	return serv.globalDB
}

// 打开Mysql数据库
func (serv *CellServGlobal) OpenMysqlDB(connStr string) bool {
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
		panic("open mysql faile")
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetMaxOpenConns(32)
	sqlDB.SetConnMaxLifetime(time.Hour)
	serv.globalDB = db
	return true
}

func (serv *CellServGlobal) Run() {
	lisn, err := net.Listen("tcp", serv.appConfig.RpcAddr)
	if err != nil {
		loghlp.Errorf("rpclisten error:%s", err.Error())
		return
	}
	worerNumOpts := grpc.NumStreamWorkers(1)
	s := grpc.NewServer(worerNumOpts)
	pbrpc.RegisterGlobalDataServer(s, serv.rpcServ)
	if err := s.Serve(lisn); err != nil {
		loghlp.Errorf("rpcserv error:%s", err.Error())
	}
}
func (serv *CellServGlobal) GetIDGenerator() *utils.IDGenerator {
	return serv.idGenrator
}
