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
	// 可继续补充用户信息修改、密码修改等接口
}
