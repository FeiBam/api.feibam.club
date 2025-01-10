package middleware

import (
	"api-feibam-club/utils"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func SecurityHeaders(ctx *gin.Context) {
	ctx.Header("X-Frame-Options", "DENY")

	// 内容安全策略
	ctx.Header("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';")

	// 防止 XSS 攻击
	ctx.Header("X-XSS-Protection", "1; mode=block")

	// 强制使用 HTTPS
	ctx.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

	// 限制引用来源
	ctx.Header("Referrer-Policy", "strict-origin")

	// 防止 MIME 类型混淆
	ctx.Header("X-Content-Type-Options", "nosniff")

	// 限制权限
	ctx.Header("Permissions-Policy", "geolocation=(), midi=(), sync-xhr=(), microphone=(), camera=(), magnetometer=(), gyroscope=(), fullscreen=(self), payment=()")

	ctx.Header("Server", "gin")

	// 继续后续请求处理
	ctx.Next()
}

func IsLogin(ctx *gin.Context) {
	// 获取 Authorization 头
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(401, gin.H{"error": "Authorization header required"})
		ctx.Abort()
		return
	}

	// 检查 Bearer 前缀
	const prefix = "Bearer "
	if len(authHeader) <= len(prefix) || authHeader[:len(prefix)] != prefix {
		ctx.JSON(401, gin.H{"error": "Invalid authorization header format"})
		ctx.Abort()
		return
	}

	// 提取 token
	tokenString := authHeader[len(prefix):]

	// 验证 token
	mapClaims, err := utils.ValidateJWT(tokenString)
	if err != nil {
		ctx.JSON(401, gin.H{"error": "Invalid token"})
		ctx.Abort()
		return
	}

	// 获取 user_name 字段
	userName, ok := mapClaims["user_name"].(string)
	if !ok {
		ctx.JSON(401, gin.H{"error": "user_name not found in token"})
		ctx.Abort()
		return
	}

	// 从 tokenStore 中检查是否存在
	tokenStore := utils.GetTokenStoreFromContext(ctx)
	tokenInfo, exists := tokenStore.Get(userName)
	if !exists {
		ctx.JSON(401, gin.H{"error": "Token not found!"})
		ctx.Abort()
		return
	}
	if tokenInfo.Token != tokenString {
		ctx.JSON(401, gin.H{"error": "Invalid token"})
		ctx.Abort()
		return
	}

	// 将解析后的 claims 设置到上下文中，供后续处理使用
	ctx.Set("claims", mapClaims)
	ctx.Next()
}

func XResponseTime(ctx *gin.Context) {
	start := time.Now()
	ctx.Next()                                                 // 继续执行其他中间件和最终的处理函数
	duration := time.Since(start)                              // 计算处理时间
	ctx.Header("X-Response-Time", fmt.Sprintf("%v", duration)) // 设置响应头
}
