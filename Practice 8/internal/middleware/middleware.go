package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"practice-8/pkg/utils"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		userClaims, ok := claims.(*utils.Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		if userClaims.Role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

type RateLimiterConfig struct {
	RequestsPerMinute int
}

func RateLimiter(config RateLimiterConfig) gin.HandlerFunc {
	type RequestTracker struct {
		count     int
		resetTime time.Time
	}

	trackers := make(map[string]*RequestTracker)
	var mu sync.Mutex

	return func(c *gin.Context) {
		mu.Lock()
		defer mu.Unlock()

		key := c.ClientIP()
		if claims, exists := c.Get("claims"); exists {
			if userClaims, ok := claims.(*utils.Claims); ok {
				key = fmt.Sprintf("user:%d", userClaims.UserID)
			}
		}

		now := time.Now()
		tracker, exists := trackers[key]

		if !exists || now.After(tracker.resetTime) {
			trackers[key] = &RequestTracker{
				count:     1,
				resetTime: now.Add(1 * time.Minute),
			}
			c.Next()
			return
		}

		if tracker.count >= config.RequestsPerMinute {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}

		tracker.count++
		c.Next()
	}
}
