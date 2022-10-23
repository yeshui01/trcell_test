package accrouter

import (
	"net/http"
	"time"
	"trcell/pkg/crossdef"
	"trcell/pkg/loghlp"
	"trcell/pkg/protocol"
	"trcell/pkg/sconst"

	"github.com/gin-gonic/gin"
)

// 注册
func registerAccount(c *gin.Context) {
	var req AccountRegisterReq
	if err := c.BindJSON(&req); err != nil {
		loghlp.Errorf("req bind json fail!!!")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	accDB := accApp.GetAccountDB()
	var userAccount OrmUser = OrmUser{
		UserName: req.UserName,
	}
	if len(req.UserName) < 5 || len(req.UserName) > 64 {
		c.JSON(http.StatusOK, gin.H{
			"code": protocol.ECodeParamError,
			"msg":  "账号名字格式错误",
		})
		return
	}
	errdb := accDB.Model(userAccount).Where("user_name=?", req.UserName).First(&userAccount).Error
	if errdb == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": protocol.ECodeAccNameHasExisted,
			"msg":  "账号已经存在",
		})
		return
	}
	if userAccount.UserID > 0 {
		// 账号已经存在
		c.JSON(http.StatusOK, gin.H{
			"code": protocol.ECodeAccNameHasExisted,
			"msg":  "db error",
		})
		return
	}
	if len(req.Pswd) < 6 || len(req.Pswd) > 12 {
		// 密码长度不符合
		c.JSON(http.StatusOK, gin.H{
			"code": protocol.ECodeParamError,
			"msg":  "param error",
		})
		return
	}
	// 注册
	userAccount.UserName = req.UserName
	userAccount.RegisterTime = time.Now().Unix()
	userAccount.Pswd = req.Pswd
	userAccount.Status = 0
	var serverInfo OrmServerList = OrmServerList{
		Recommend: 1,
	}
	errHall := accDB.Model(serverInfo).Where("recommend=?", serverInfo.Recommend).First(&serverInfo).Error
	if errHall != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": protocol.ECodeAccNotExisted,
			"msg":  "db error",
		})
		//return
	}
	userAccount.DataZone = serverInfo.ID // 目前默认1
	if userAccount.DataZone == 0 {
		userAccount.DataZone = 1
	}
	errdb = accDB.Create(&userAccount).Error
	if errdb != nil {
		// 账号已经存在
		c.JSON(http.StatusOK, gin.H{
			"code": protocol.ECodeAccNameHasExisted,
			"msg":  "db error",
		})
		return
	}

	// 注册成功
	var rep AccountRegisterRsp
	rep.UserID = userAccount.UserID

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "register success",
		"data": rep,
	})
}

// 登录
func loginAccount(c *gin.Context) {
	var req AccountLoginReq
	if err := c.BindJSON(&req); err != nil {
		loghlp.Errorf("req bind json fail!!!")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	accDB := accApp.GetAccountDB()
	var userAccount OrmUser = OrmUser{
		UserName: req.UserName,
	}
	// 先从缓存取
	dataCache := accApp.GetCache()
	val, errCache := dataCache.Get(req.UserName)
	if errCache == nil {
		userAccount = val.(OrmUser)
		loghlp.Debugf("get username(%s) info from cache!", req.UserName)
	} else {
		errdb := accDB.Model(userAccount).Where("user_name=?", req.UserName).First(&userAccount).Error
		if errdb != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": protocol.ECodeAccNotExisted,
				"msg":  "db error",
			})
			return
		}
	}
	// 密码验证
	if req.Pswd != userAccount.Pswd {
		loghlp.Errorf("check user(%s) password fail,pswd:%s, input pswd:%s", req.UserName, userAccount.Pswd, req.Pswd)
		c.JSON(http.StatusOK, gin.H{
			"code": protocol.ECodeAccPasswordError,
			"msg":  "密码错误",
		})
		return
	}
	nowTime := time.Now().Unix()
	// 生成token
	token, errToken := genJwtToken(userAccount.UserID, userAccount.UserName, userAccount.DataZone)
	if errToken != nil {
		// 系统错误
		c.JSON(http.StatusOK, gin.H{
			"code": protocol.ECodeSysError,
			"msg":  "sys error",
		})
		return
	}
	var serverInfo OrmServerList = OrmServerList{
		ID: userAccount.DataZone,
	}
	errZone := accDB.Model(serverInfo).Where("id=?", userAccount.DataZone).First(&serverInfo).Error
	if errZone != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": protocol.ECodeAccNotExisted,
			"msg":  "db error",
		})
		return
	}
	var rep AccountLoginRsp = AccountLoginRsp{
		Token:      token,
		ServerAddr: serverInfo.GateAddr,
		RestTime:   int32((userAccount.RegisterTime + int64(sconst.AccountCertificationTime)) - nowTime),
	}

	// 解析token
	jwtObj := crossdef.NewJWT()
	jwtObj.SetSignKey(crossdef.SignKey)
	claimData, errJwt := jwtObj.ParseToken(token)
	if errJwt == nil {
		loghlp.Infof("parse token success:%+v", claimData)
	} else {
		loghlp.Errorf("parse token fail:%s", errJwt.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "login success",
		"data": rep,
	})
	// 缓存
	dataCache.SetWithExpire(req.UserName, userAccount, time.Second*600)
}

// 查询公告
func queryNotice(c *gin.Context) {
	var req QueryNoticeReq
	if err := c.BindJSON(&req); err != nil {
		loghlp.Errorf("req bind json fail!!!")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	accDB := accApp.GetAccountDB()
	var cellNotice OrmCellNotice = OrmCellNotice{}
	errdb := accDB.Model(cellNotice).First(&cellNotice).Error
	if errdb != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": protocol.ECodeNotFindNotice,
			"msg":  "not find notice",
		})
		return
	}
	var rep QueryNoticeRsp = QueryNoticeRsp{
		ID:     cellNotice.ID,
		Notice: cellNotice.Content,
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "login success",
		"data": rep,
	})
}
