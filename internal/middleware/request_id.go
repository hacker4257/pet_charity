package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/hacker4257/pet_charity/pkg/utils"
)

// RequestID 为每个请求生成唯一 ID，写入 header 和 context
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-ID")
		if id == "" {
			id = utils.RandomCode(16)
		}
		c.Set("request_id", id)
		c.Header("X-Request-ID", id)
		c.Next()
	}
}
