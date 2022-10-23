/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 11:10:29
 * @LastEditTime: 2022-09-30 11:53:11
 * @FilePath: \trcell\app\servdata\servdatamain\servdata_global.go
 */
package servdatamain

import (
	"time"
	"trcell/pkg/loghlp"
	"trcell/pkg/tbobj"
	"trcell/pkg/timeutil"

	"gorm.io/gorm"
)

type DataDBJob struct {
	DoJob func()
}

type ServDataGlobal struct {
	gameDB         *gorm.DB
	dataPlayers    map[int64]*DataPlayer
	lastUpdateTime int64
	dbJobCh        chan *DataDBJob
	dbJobStop      bool
	// 全局数据表
	DataTbCsGlobal *tbobj.TbCsGlobal
}

func NewServDataGlobal() *ServDataGlobal {
	return &ServDataGlobal{
		gameDB:         nil,
		dataPlayers:    make(map[int64]*DataPlayer),
		dbJobCh:        make(chan *DataDBJob, 2048),
		dbJobStop:      false,
		DataTbCsGlobal: tbobj.NewTbCsGlobal(),
	}
}

func (servGlobal *ServDataGlobal) GetGameDB() *gorm.DB {
	return servGlobal.gameDB
}
func (servGlobal *ServDataGlobal) InitGameDB(gameDB *gorm.DB) {
	servGlobal.gameDB = gameDB
	servGlobal.StartDBRun()
}

func (servGlobal *ServDataGlobal) FindDataPlayer(roleID int64) *DataPlayer {
	if p, ok := servGlobal.dataPlayers[roleID]; ok {
		p.VisitTime = timeutil.NowTime()
		return p
	}
	return nil
}
func (servGlobal *ServDataGlobal) AddDataPlayer(roleID int64, player *DataPlayer) {
	if _, ok := servGlobal.dataPlayers[roleID]; ok {
		return
	}
	player.VisitTime = timeutil.NowTime()
	servGlobal.dataPlayers[roleID] = player
}
func (servGlobal *ServDataGlobal) FrameUpdate(curTimeMs int64) {
	servGlobal.Update(curTimeMs / 1000)
}
func (servGlobal *ServDataGlobal) Update(curTime int64) {
	if curTime <= servGlobal.lastUpdateTime {
		return
	}
	loghlp.Debugf("dataGlobal Update:%d", curTime)
	servGlobal.lastUpdateTime = curTime
	var delList []int64
	for _, v := range servGlobal.dataPlayers {
		if v.VisitTime == 0 {
			v.VisitTime = curTime
		}
		if curTime-v.VisitTime >= 1200 {
			delList = append(delList, v.RoleID)
		}
	}
	for _, v := range delList {
		delete(servGlobal.dataPlayers, v)
	}

}

// DB更新线程
func (servGlobal *ServDataGlobal) StartDBRun() {
	go func() {
		for {
			select {
			case dbOpt, ok := <-servGlobal.dbJobCh:
				{
					if ok {
						loghlp.Debugf("do db job")
						dbOpt.DoJob()
					} else {
						servGlobal.dbJobStop = true
					}
					break
				}
			}
			if servGlobal.dbJobStop {
				break
			}
		}
		loghlp.Info("exit db run")
	}()
}
func (servGlobal *ServDataGlobal) StopDBRun() {
	close(servGlobal.dbJobCh)
	// 等待结束
	for {
		if servGlobal.dbJobStop {
			time.Sleep(time.Second)
			break
		} else {
			time.Sleep(time.Second)
		}
	}
}
func (servGlobal *ServDataGlobal) PostDBJob(dbJob *DataDBJob) {
	servGlobal.dbJobCh <- dbJob
}
