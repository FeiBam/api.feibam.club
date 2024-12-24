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

func IsLogin(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" { // 假设 isValidToken 检查 token 合法性
		c.JSON(401, utils.JsonResponse("ok", 401, "", "Unauthorized", nil))
		c.Abort() // 中断请求链，后续的中间件或控制器不会执行
		return
	}
	c.Next()
}

func XResponseTime(ctx *gin.Context) {
	start := time.Now()
	ctx.Next()                                                 // 继续执行其他中间件和最终的处理函数
	duration := time.Since(start)                              // 计算处理时间
	ctx.Header("X-Response-Time", fmt.Sprintf("%v", duration)) // 设置响应头
}
