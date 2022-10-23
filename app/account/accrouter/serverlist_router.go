package accrouter

import (
	"net/http"
	"strings"
	"trcell/pkg/protocol"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// 获取服务器列表
func getServerList(c *gin.Context) {
	// var req QueryNoticeReq
	// if err := c.BindJSON(&req); err != nil {
	// 	loghlp.Errorf("req bind json fail!!!")
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	accDB := accApp.GetAccountDB()
	var serverListInfo []OrmServerList = make([]OrmServerList, 0)
	errdb := accDB.Find(&serverListInfo).Error
	if errdb != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": protocol.ECodeNotFindNotice,
			"msg":  "not find notice",
		})
		return
	}
	var rep QueryServerListRsp = QueryServerListRsp{
		ServerList: make([]*ServerInfo, 0),
	}
	for _, v := range serverListInfo {
		gateAddrList := strings.Split(v.GateAddr, "|")
		chooseAddr := ""
		choosePort := int32(0)
		if len(gateAddrList) > 0 {
			ipAndPort := strings.Split(gateAddrList[0], ":")
			if len(ipAndPort) > 1 {
				chooseAddr = ipAndPort[0]
				choosePort = int32(cast.ToInt(ipAndPort[1]))
			}
		}
		rep.ServerList = append(rep.ServerList, &ServerInfo{
			ServerID:   v.ID,
			ServerName: v.ServerName,
			Addr:       chooseAddr,
			Port:       choosePort,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "login success",
		"data": rep,
	})
}
