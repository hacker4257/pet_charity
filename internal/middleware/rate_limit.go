package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hacker4257/pet_charity/internal/database"
	"github.com/hacker4257/pet_charity/pkg/response"
)

// RateLimit 基于 Redis 的固定窗口限流中间件
// prefix: 限流键前缀（区分不同接口组）
// limit:  窗口内允许的最大请求数
// window: 窗口时长
func RateLimit(prefix string, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf("rate:%s:%s", prefix, c.ClientIP())

		ctx := context.Background()
		count, err := database.RDB.Incr(ctx, key).Result()
		if err != nil {
			// Redis 故障时放行，不阻塞业务
			c.Next()
			return
		}

		// 首次计数时设置过期
		if count == 1 {
			database.RDB.Expire(ctx, key, window)
		}

		// 写入限流响应头
		ttl, _ := database.RDB.TTL(ctx, key).Result()
		c.Header("X-RateLimit-Limit", strconv.Itoa(limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(max(0, limit-int(count))))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(ttl).Unix(), 10))

		if int(count) > limit {
			c.Header("Retry-After", strconv.Itoa(int(ttl.Seconds())))
			response.Fail(c, http.StatusTooManyRequests, 429, "too many requests, please try again later")
			c.Abort()
			return
		}

		c.Next()
	}
}
