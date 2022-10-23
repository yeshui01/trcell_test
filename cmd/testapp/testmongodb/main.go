/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-10-12 15:23:00
 * @LastEditTime: 2022-10-12 17:03:25
 * @FilePath: \trcell\cmd\testapp\testmongodb\main.go
 */
package main

import (
	"context"
	"time"
	"trcell/pkg/loghlp"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	loghlp.ActiveConsoleLog()
	loghlp.SetConsoleLogLevel(logrus.DebugLevel)
	loghlp.Infof("test mongo db")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	monClientOpts := options.Client().ApplyURI("mongodb://localhost:27017")
	//var err error
	c, err := mongo.Connect(context.TODO(), monClientOpts)
	if err != nil {
		loghlp.Errorf("connect mongo db error:%s", err.Error())
		return
	}
	// 关闭连接
	defer c.Disconnect(context.TODO())
	err = c.Ping(context.TODO(), nil)
	if err != nil {
		loghlp.Errorf("mongo db ping error:%s", err.Error())
		return
	}
	loghlp.Infof("connect to mongdb succ!!!")
	// 定义一个文档
	type TbStudent struct {
		Name string
		Age  int32
		Data []byte
	}
	docCol := c.Database("testdb").Collection("student")
	// 插入
	insertResult, errInsert := docCol.InsertOne(context.TODO(), TbStudent{Name: "Tombinary", Age: 10, Data: []byte("hello world")})
	if errInsert != nil {
		loghlp.Errorf("insert doc error:%s", errInsert.Error())
		return
	}
	objectID := insertResult.InsertedID.(primitive.ObjectID)
	loghlp.Infof("insert one doc data, insertID:%+v", objectID)
	loghlp.Infof("objectID = %s", objectID.Hex())
	// 更新
	update := bson.D{{"$set", bson.D{{Key: "name", Value: "Tom61"}, {Key: "age", Value: "119"}}}}
	ur, errUpd := docCol.UpdateMany(context.TODO(), bson.D{{Key: "name", Value: "Tom6"}}, update)
	if errUpd != nil {
		loghlp.Errorf("upd error:%s", errUpd.Error())
		return
	}
	loghlp.Infof("upd succ:%d", ur.MatchedCount)
	// 查询
	// Find 查询所有
	// FindOne 查询单个
	fr, errFind := docCol.Find(context.TODO(), bson.D{{Key: "name", Value: "Tom7"}})
	if errFind != nil {
		loghlp.Errorf("find data error:%s", errFind.Error())
		return
	}
	for fr.Next(context.Background()) {
		var result bson.D
		errTmp := fr.Decode(&result)
		if errTmp != nil {
			loghlp.Errorf("fr.Decode error:%s", errTmp.Error())
			continue
		}
		loghlp.Infof("frOneResult:%s", result.Map()["name"])
	}
	// 根据主键查找单个
	findId, _ := primitive.ObjectIDFromHex("63466e93f9fa473dde6c9de2")
	fsr := docCol.FindOne(context.Background(), bson.D{{Key: "_id", Value: findId}})
	if fsr == nil {
		loghlp.Errorf("find single data error")
		return
	}
	{
		var result bson.D
		errTmp := fsr.Decode(&result)
		if errTmp != nil {
			loghlp.Errorf("fsr.Decode error:%s", errTmp.Error())
		} else {
			loghlp.Infof("find single data by key succ:%s", result.Map()["name"])
		}
	}
	// 删除
	// 删除一个名字为Tom8的数据
	delR, errDel := docCol.DeleteOne(ctx, bson.D{{Key: "name", Value: "Tom8"}})
	if errDel != nil {
		loghlp.Errorf("delete data error:%s", errDel.Error())
	} else {
		loghlp.Infof("delete single data succ, delCount:%d", delR.DeletedCount)
	}
}
