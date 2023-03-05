/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 11:10:29
 * @LastEditTime: 2022-09-30 11:53:11
 * @FilePath: \trcell\app\servdata\servdatamain\servdata_global.go
 */
package servdatamain

import (
	"sync"
	"time"
	"trcell/pkg/loghlp"
	"trcell/pkg/tbobj"
	"trcell/pkg/timeutil"

	"gorm.io/gorm"
)

type DataDBJob struct {
	DoJob func() bool
}

type ServDataGlobal struct {
	gameDB         *gorm.DB
	dataPlayers    map[int64]*DataPlayer
	lastUpdateTime int64
	dbJobCh        chan *DataDBJob
	dbJobStopCh    chan bool
	dbJobStop      bool
	dbWaitWg       sync.WaitGroup
	// 全局数据表
	DataTbCsGlobal *tbobj.TbCsGlobal
}

func NewServDataGlobal() *ServDataGlobal {
	return &ServDataGlobal{
		gameDB:         nil,
		dataPlayers:    make(map[int64]*DataPlayer),
		dbJobCh:        make(chan *DataDBJob, 2048),
		dbJobStop:      false,
		dbJobStopCh:    make(chan bool),
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
	servGlobal.dbWaitWg.Add(1)
	go func() {
		var lastDBTime int64 = 0
		sectick := time.NewTicker(time.Second)
		for {
			select {
			case dbOpt, ok := <-servGlobal.dbJobCh:
				{
					if ok {
						servGlobal.dbJobStop = false
						loghlp.Debugf("do db job")
						var tryCount int = 0
						for !dbOpt.DoJob() {
							if lastDBTime > 0 { // 说明此时已经开始计时了
								lastDBTime = timeutil.NowTime()
							}
							tryCount++
							time.Sleep(time.Second) // 休眠1秒继续尝试
							loghlp.Errorf("do db job fail, try continue")
							if tryCount >= 10 {
								break
							}
						}
						if lastDBTime > 0 { // 说明此时已经开始计时了
							lastDBTime = timeutil.NowTime()
						}
					}
					break
				}
			case <-sectick.C:
				{
					// 5秒钟还没有任务,那就可以停掉了
					if lastDBTime > 0 && timeutil.NowTime()-lastDBTime >= 5 {
						servGlobal.dbJobStop = true
						sectick.Stop()
					}
					break
				}
			case <-servGlobal.dbJobStopCh:
				{
					// 停止
					lastDBTime = timeutil.NowTime() // 开始计时
					break
				}
			}

			if servGlobal.dbJobStop {
				break
			}
		}
		loghlp.Info("exit db run")
		servGlobal.dbWaitWg.Done()
	}()
}
func (servGlobal *ServDataGlobal) StopDBRun() {
	close(servGlobal.dbJobStopCh)
	// 等待任务执行结束
	servGlobal.dbWaitWg.Wait()
	close(servGlobal.dbJobCh)
}

func (servGlobal *ServDataGlobal) PostDBJob(dbJob *DataDBJob) {
	servGlobal.dbJobCh <- dbJob
}
