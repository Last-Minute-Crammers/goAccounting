package router

import (
	"goAccounting/internal/api/v1"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var Engine *gin.Engine

func init() {
	gin.SetMode(gin.DebugMode)
	Engine = gin.Default()
	
	// 配置CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	Engine.Use(cors.New(config))
	
	// 添加日志中间件
	Engine.Use(gin.Logger())
	Engine.Use(gin.Recovery())
	
	// 健康检查路由
	Engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "Server is running"})
	})
	
	// 公共路由组（不需要认证）
	publicV1 := Engine.Group("/public")
	{
		publicApi := v1.PublicApi{}
		publicV1.POST("/user/login", publicApi.Login)
		publicV1.POST("/user/register", publicApi.Register)
	}
	
	// API路由组（需要认证）
	apiV1 := Engine.Group("/user")
	{
		// 用户相关
		userApi := v1.UserApi{}
		apiV1.GET("/home", userApi.Home)
		apiV1.GET("/stats", userApi.GetStats)
		apiV1.PUT("/info", userApi.UpdateInfo)
		apiV1.PUT("/password", userApi.UpdatePassword)
		
		// 分类相关
		categoryApi := v1.CategoryApi{}
		apiV1.GET("/category/list", categoryApi.GetList)
		apiV1.POST("/category/", categoryApi.CreateOne)
		apiV1.GET("/category/:id", categoryApi.GetOne)
		apiV1.PUT("/category/:id", categoryApi.UpdateOne)
		apiV1.DELETE("/category/:id", categoryApi.DeleteOne)
		
		// 交易相关
		transactionApi := v1.TransactionApi{}
		apiV1.POST("/transaction/", transactionApi.CreateOne)
		apiV1.GET("/transaction/list", transactionApi.GetList)
		apiV1.GET("/transaction/:id", transactionApi.GetOne)
		apiV1.PUT("/transaction/:id", transactionApi.UpdateOne)
		apiV1.DELETE("/transaction/:id", transactionApi.DeleteOne)
		apiV1.GET("/transaction/statistic/month", transactionApi.GetMonthStatistic)
		apiV1.GET("/transaction/statistic/category", transactionApi.GetCategoryStats)
		
		// 好友相关
		apiV1.GET("/friend/list", userApi.GetFriendList)
		apiV1.POST("/friend/invitation", userApi.CreateFriendInvitation)
		apiV1.GET("/friend/invitation", userApi.GetFriendInvitationList)
		apiV1.PUT("/friend/invitation/:id/accept", userApi.AcceptFriendInvitation)
		apiV1.PUT("/friend/invitation/:id/refuse", userApi.RefuseFriendInvitation)
	}
}
