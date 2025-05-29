package router

import (
	"net/http"
	"github.com/gin-gonic/gin"
	v1 "goAccounting/internal/api/v1"
)

func RegisterUserRoutes() {
	if NoTourist == nil {
		return
	}
	
	userApi := v1.UserApi{}
	
	userGroup := NoTourist.Group("/user")
	{
		// 用户基础信息
		userGroup.GET("/home", userApi.Home)
		
		// 好友相关
		friendGroup := userGroup.Group("/friend")
		{
			friendGroup.GET("/list", userApi.GetFriendList)
			
			// 好友邀请
			invitationGroup := friendGroup.Group("/invitation")
			{
				invitationGroup.GET("", userApi.GetFriendInvitationList)
				invitationGroup.POST("", userApi.CreateFriendInvitation)
				invitationGroup.PUT("/:id/accept", userApi.AcceptFriendInvitation)
				invitationGroup.PUT("/:id/refuse", userApi.RefuseFriendInvitation)
			}
		}
	}
}
