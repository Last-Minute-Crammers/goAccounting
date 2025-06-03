package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取 Authorization header
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Missing authorization header",
			})
			ctx.Abort()
			return
		}

		// 检查 Bearer token 格式
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

		// TODO: 在这里添加 JWT token 验证逻辑
		// 暂时跳过验证，实际项目中需要验证 JWT token
		// if !validateJWTToken(tokenString) {
		//     ctx.JSON(http.StatusUnauthorized, gin.H{
		//         "code":    401,
		//         "message": "Invalid token",
		//     })
		//     ctx.Abort()
		//     return
		// }

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
