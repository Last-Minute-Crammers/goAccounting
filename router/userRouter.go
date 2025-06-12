package router

import (
	v1 "goAccounting/internal/api/v1"
)

func RegisterUserRoutes() {
	userApi := new(v1.UserApi)
	// 用户主页
	Private.GET("/user/home", userApi.Home)
	// 好友相关
	Private.GET("/user/friend/list", userApi.GetFriendList)
	Private.POST("/user/friend/invitation", userApi.CreateFriendInvitation)
	Private.GET("/user/friend/invitation", userApi.GetFriendInvitationList)
	Private.PUT("/user/friend/invitation/:id/accept", userApi.AcceptFriendInvitation)
	Private.PUT("/user/friend/invitation/:id/refuse", userApi.RefuseFriendInvitation)
	
	// 社交互动功能
	Private.GET("/user/friend/:id/consumption", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Friend consumption distribution", "data": []gin.H{
			{"category": "餐饮", "percentage": 30},
			{"category": "娱乐", "percentage": 25},
			{"category": "交通", "percentage": 15},
		}})
	})
	
	// 好友PK排行榜
	Private.GET("/user/friend/ranking", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Friend savings ranking", "data": []gin.H{
			{"username": "friend1", "savings": 5000, "rank": 1},
			{"username": "friend2", "savings": 3500, "rank": 2},
		}})
	})
	
	// 吐槽功能
	Private.POST("/user/friend/:id/comment", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Comment posted successfully"})
	})
	Private.GET("/user/friend/:id/comments", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Friend comments", "data": []gin.H{
			{"content": "你这个月买书花得有点多哦～", "author": "friend1"},
		}})
	})
	
	// 存钱目标相关
	Private.POST("/user/savings/goal", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Savings goal set successfully"})
	})
	Private.GET("/user/savings/goal", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Current savings goal", "data": gin.H{
			"target": 10000, "current": 6500, "progress": 65,
		}})
	})
	
	// 可继续补充用户信息修改、密码修改等接口
}
