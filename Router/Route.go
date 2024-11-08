package Router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go-edgetts/edgetts"
	"time"
)

// SetupRouter 初始化并配置gin框架的路由器
// 返回值 *gin.Engine: 返回配置好的gin路由器实例,可以进一步挂载其他路由或进行服务器监听
func SetupRouter() *gin.Engine {
	// 设置GIN运行模式为DebugMode
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default() // 创建默认配置的gin路由器

	// 设置 CORS 配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 可以写入具体的允许域名，例如：http://example.com
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,           // 是否允许跨域请求携带 Cookie
		MaxAge:           12 * time.Hour, // 预检请求的缓存时间
	}))

	// 为路由器启用恢复中间件,可以在发生panic时恢复
	r.Use(gin.Recovery())
	r.POST("/api/sendGec", edgetts.GecController) //接收Gec结果
	r.GET("/api/getGec", edgetts.GetGec)          //提供外部调用的Gec接口
	return r                                      // 返回配置好的路由器实例
}
