package middleware

import (
	"goAccounting/initialize"
	"goAccounting/util/jwtTool"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Missing authorization header",
			})
			ctx.Abort()
			return
		}

		tokenString := ""
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = authHeader[7:]
		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Invalid authorization format",
			})
			ctx.Abort()
			return
		}

		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Missing token",
			})
			ctx.Abort()
			return
		}

		// 使用新的解析方法同时获取 userId 和 isAdmin
		userId, isAdmin, err := jwtTool.ParseAdminFromToken(tokenString, []byte(initialize.Config.System.JwtKey))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Invalid token: " + err.Error(),
			})
			ctx.Abort()
			return
		}

		ctx.Set("userId", userId)
		ctx.Set("isAdmin", isAdmin)
		ctx.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		isAdmin, exists := ctx.Get("isAdmin")
		if !exists || !isAdmin.(bool) {
			ctx.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "Admin access required",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// CORS 中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}

// 数据库注入中间件
func InjectDB() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 直接使用全局数据库连接
		ctx.Set("db", initialize.Database)
		ctx.Next()
	}
}
