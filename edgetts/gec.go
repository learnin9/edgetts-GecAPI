package edgetts

import (
	"github.com/gin-gonic/gin"
	"go-edgetts/logger"
	"net/http"
)

// GecInfo 结构体用于存储与 GEC 相关的信息，包括消息、版本、更新时间和过期时间。
type GecInfo struct {
	SecMsgEC            string `json:"Sec-MS-GEC"`                 // GEC 消息内容
	SecMSGECVersion     string `json:"Sec-MS-GEC-Version"`         // GEC 版本
	LastUpdate          int64  `json:"last_update"`                // 最后更新时间 (时间戳)
	LastUpdateFormatUTC string `json:"last_update_format(UTC +8)"` // 最后更新时间的格式化字符串 (UTC +8)
	Expiration          int64  `json:"expiration"`                 // 过期时间 (时间戳)
}

// GecController 处理 GEC API 的请求并返回欢迎信息
func GecController(c *gin.Context) {
	var gecInfo GecInfo
	err := c.ShouldBindJSON(&gecInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logger.SugarLogger.Infow("Received GEC request", "获取到的Gec数据", gecInfo)
	c.JSON(200, gin.H{
		"code": "200",
		"msg":  "ok",
	})
	if !checkGec(gecInfo.SecMsgEC) {
		logger.SugarLogger.Debugf("Sec-MS-GEC-->%s已更新, 保存到Redis", gecInfo.SecMsgEC)
		err = redisSaveEdgeGec(&gecInfo)
		if err != nil {
			logger.SugarLogger.Debugf("存储在Redis失败: %s", err.Error())
			c.JSON(500, gin.H{
				"code": "500",
				"msg":  "存储在Redis失败",
			})
			return
		}
	}
	return
}

// GetGec 是处理 GEC API 请求的函数，返回欢迎信息。
func GetGec(c *gin.Context) {
	gecInfo, err := redisGetEdgeGec()
	if err != nil {
		logger.SugarLogger.Debugf("Error saving Edge Gec To Redis: %s", err.Error())
		c.JSON(500, gin.H{
			"code": "500",
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(200, gecInfo)
	return
}
