package router

import (
    v1 "goAccounting/internal/api/v1"
    "goAccounting/router/middleware"
    "github.com/gin-gonic/gin"
)

type AdminRouter struct{}

func (ar *AdminRouter) InitAdminRouter(group *gin.RouterGroup) {
    adminApi := v1.AdminApi{}
    
    // 所有 admin 路由都需要 JWT 认证 + 管理员权限
    adminGroup := group.Group("/admin").
        Use(middleware.JWTAuth()).
        Use(middleware.AdminOnly())
    {
        // 用户管理
        adminGroup.GET("/users", adminApi.ListUsers)
        adminGroup.DELETE("/users", adminApi.DeleteUser)
        adminGroup.PUT("/users/admin", adminApi.UpdateUserAdmin)
        
        // 后续可以扩展交易记录管理
        // adminGroup.GET("/transactions", adminApi.ListTransactions)
        // adminGroup.DELETE("/transactions", adminApi.DeleteTransaction)
    }
}