package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimit(redisClient *redis.Client, maxRequests int, window time.Duration) gin.HandlerFunc {
	if redisClient == nil {
		return func(c *gin.Context) { c.Next() }
	}

	return func(c *gin.Context) {
		ctx := context.Background()
		key := "ratelimit:" + c.ClientIP()

		count, err := redisClient.Incr(ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}

		if count == 1 {
			_ = redisClient.Expire(ctx, key, window).Err()
		}

		if count > int64(maxRequests) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "too many requests"})
			return
		}

		c.Next()
	}
}
