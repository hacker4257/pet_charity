package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hacker4257/pet_charity/pkg/logger"
)

// AccessLog 替代 gin.Logger()，输出结构化请求日志
func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next() // 执行后续 handler

		latency := time.Since(start)
		status := c.Writer.Status()

		fields := []logger.Field{
			logger.Int("status", status),
			logger.Str("method", c.Request.Method),
			logger.Str("path", path),
			logger.Str("query", query),
			logger.Str("ip", c.ClientIP()),
			logger.Dur("latency", latency),
			logger.Int("body_size", c.Writer.Size()),
		}

		// 如果有登录用户，也记录下来
		if userID, exists := c.Get("user_id"); exists {
			fields = append(fields, logger.Any("user_id", userID))
		}

		// 根据状态码选择日志级别
		switch {
		case status >= 500:
			fields = append(fields, logger.Str("error", c.Errors.String()))
			logger.Error("request", fields...)
		case status >= 400:
			logger.Warn("request", fields...)
		default:
			logger.Info("request", fields...)
		}
	}
}
