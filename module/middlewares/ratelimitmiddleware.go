package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	rateLimit         = 5
	rateResetDuration = time.Minute
	userRequests      = make(map[string]int)
	mu                sync.Mutex
)

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIP := c.ClientIP()

		mu.Lock()
		if count, exists := userRequests[userIP]; exists {
			if count >= rateLimit {
				c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded, try again later"})
				mu.Unlock()
				c.Abort()
				return
			}
			userRequests[userIP]++
		} else {
			userRequests[userIP] = 1
		}
		mu.Unlock()
		c.Next()
		go func() {
			time.Sleep(rateResetDuration)
			mu.Lock()
			userRequests[userIP]--
			if userRequests[userIP] <= 0 {
				delete(userRequests, userIP)
			}
			mu.Unlock()
		}()
	}
}
