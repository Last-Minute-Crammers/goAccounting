package middleware

import (
	"goAccounting/internal/api/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从请求头获取token
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			response.TokenExpired(ctx)
			return
		}
		
		// 检查Bearer前缀
		tokenString := ""
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = authHeader[7:]
		} else {
			response.TokenExpired(ctx)
			return
		}
		
		if tokenString == "" {
			response.TokenExpired(ctx)
			return
		}
		
		// TODO: 验证JWT token
		// 这里需要实现JWT token的验证逻辑
		// claims, err := verifyToken(tokenString)
		// if err != nil {
		//     response.TokenExpired(ctx)
		//     return
		// }
		
		// 暂时模拟用户ID，实际应从JWT claims中获取
		ctx.Set("userID", "mock_user_id")
		ctx.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
}
