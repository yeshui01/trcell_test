package accrouter

import (
	"trcell/app/account/accmiddleware"

	"github.com/bluele/gcache"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IAccApp interface {
	GetRouter() *gin.Engine
	GetAccountDB() *gorm.DB
	GetCache() gcache.Cache
}

var accApp IAccApp

func SetupAccountApi(app IAccApp) {
	accApp = app

	// 账号
	accApi := accApp.GetRouter().Group("/account")
	accApi.POST("/register", registerAccount)
	accApi.POST("/login", loginAccount)

	// 公告
	noticeApi := accApp.GetRouter().Group("/notice")
	noticeApi.POST("/query", queryNotice)

	// 服务器
	servApi := accApp.GetRouter().Group("/server")
	servApi.POST("/list", getServerList)

	// 限流设置
	accApp.GetRouter().Use(accmiddleware.MaxAllowed(1000)) // 限制每秒最多允许1000个请求
}
