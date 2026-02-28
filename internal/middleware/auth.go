package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hacker4257/pet_charity/internal/config"
	"github.com/hacker4257/pet_charity/pkg/response"
	"github.com/hacker4257/pet_charity/pkg/utils"
)

func JWTAuth(getTokenVersion func(uint) int) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			//websocket
			if token := c.Query("token"); token != "" {
				authHeader = "Bearer " + token
			}
		}
		if authHeader == "" {
			response.Unauthorized(c, "missing token")
			c.Abort()
			return
		}

		//2.格式必须是 "Bearer xxx"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "invalid token format")
			c.Abort()
			return
		}

		//3. 解析token
		claims, err := utils.ParseToken(parts[1], config.Global.JWT.Secret)
		if err != nil {
			response.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		//版本号验证
		if getTokenVersion != nil {
			currentVer := getTokenVersion(claims.UserID)
			if claims.TokenVersion < currentVer {
				response.Unauthorized(c, "token has been revoked")
				c.Abort()
				return
			}
		}
		//4. 把用户信息存到上下文，后续handler可以取
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		//5. 放行
		c.Next()
	}
}

func RequiredRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			response.Unauthorized(c, "missing role info")
			c.Abort()
			return
		}
		roleStr := role.(string)
		for _, r := range roles {
			if roleStr == r {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "insufficient permissions")
		c.Abort()
	}
}
