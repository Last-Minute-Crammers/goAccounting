package router

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func RegisterSocialRoutes() {
	if NoTourist == nil {
		return
	}
	
	socialGroup := NoTourist.Group("/social")
	{
		// 获取好友列表
		socialGroup.GET("/friends", func(c *gin.Context) {
			// TODO: 实现获取好友列表逻辑
			c.JSON(http.StatusOK, gin.H{
				"data": []gin.H{
					{"id": 1, "name": "张小明", "spendReduction": 15, "avatar": ""},
					{"id": 2, "name": "李小红", "spendReduction": 25, "avatar": ""},
				},
			})
		})
		
		// 搜索用户
		socialGroup.GET("/search", func(c *gin.Context) {
			// TODO: 实现搜索用户逻辑
			c.JSON(http.StatusOK, gin.H{
				"data": []gin.H{
					{"id": 3, "name": "王小刚", "avatar": ""},
				},
			})
		})
		
		// 添加好友
		socialGroup.POST("/friends/add", func(c *gin.Context) {
			// TODO: 实现添加好友逻辑
			c.JSON(http.StatusOK, gin.H{"message": "好友申请已发送"})
		})
		
		// 获取理财排行榜
		socialGroup.GET("/ranking", func(c *gin.Context) {
			// TODO: 实现获取排行榜逻辑
			c.JSON(http.StatusOK, gin.H{
				"data": []gin.H{
					{"rank": 1, "name": "张小萌", "progress": 98},
					{"rank": 2, "name": "李小明", "progress": 85},
					{"rank": 3, "name": "王小红", "progress": 78},
					{"rank": 4, "name": "你", "progress": 65},
				},
			})
		})
		
		// 获取好友申请列表
		socialGroup.GET("/friend-requests", func(c *gin.Context) {
			// TODO: 实现获取好友申请列表逻辑
			c.JSON(http.StatusOK, gin.H{"data": []gin.H{}})
		})
		
		// 处理好友申请
		socialGroup.POST("/friend-requests/:id/respond", func(c *gin.Context) {
			// TODO: 实现处理好友申请逻辑
			c.JSON(http.StatusOK, gin.H{"message": "好友申请处理成功"})
		})
	}
}
