package middleware

import (
	"net/http"
	"strings"
	"goAccounting/util/jwtTool"
	"goAccounting/initialize"

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

		userId, err := jwtTool.ParseUserIdFromToken(tokenString, []byte(jwtTool.SignKey))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Invalid token: " + err.Error(),
			})
			ctx.Abort()
			return
		}
		// 这里的 "userId" 要和 GetUserId 用的是同一个 key，区分大小写
		ctx.Set("userId", userId)

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
