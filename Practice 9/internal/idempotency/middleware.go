package idempotency

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Middleware creates an idempotency middleware
func Middleware(store Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodPost {
			c.Next()
			return
		}

		key := c.GetHeader("Idempotency-Key")
		if key == "" {
			fmt.Printf("[%s] [MIDDLEWARE] Missing Idempotency-Key, returning 400\n", time.Now().Format("15:04:05.000"))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Idempotency-Key header required",
			})
			return
		}

		if cached, exists := store.Get(key); exists {
			if cached.Completed {
				fmt.Printf("[%s] [MIDDLEWARE] Key '%s' - CACHE HIT (completed), returning 200 with cached result\n",
					time.Now().Format("15:04:05.000"), key)
				c.JSON(cached.StatusCode, gin.H{
					"status":             "completed",
					"response_body":      string(cached.Body),
					"idempotency_cached": true,
				})
				c.Abort()
				return
			}

			fmt.Printf("[%s] [MIDDLEWARE] Key '%s' - DUPLICATE DETECTED (processing), returning 409 Conflict\n",
				time.Now().Format("15:04:05.000"), key)
			c.JSON(http.StatusConflict, gin.H{
				"error":  "Duplicate request in progress",
				"status": "processing",
			})
			c.Abort()
			return
		}

		fmt.Printf("[%s] [MIDDLEWARE] Key '%s' - NEW REQUEST, starting processing\n",
			time.Now().Format("15:04:05.000"), key)
		if !store.StartProcessing(key) {
			fmt.Printf("[%s] [MIDDLEWARE] Key '%s' - RACE CONDITION (failed to reserve), returning 409 Conflict\n",
				time.Now().Format("15:04:05.000"), key)
			c.JSON(http.StatusConflict, gin.H{
				"error":  "Duplicate request in progress",
				"status": "processing",
			})
			c.Abort()
			return
		}

		writer := &responseWriter{
			ResponseWriter: c.Writer,
			statusCode:     http.StatusOK,
			body:           []byte{},
		}
		c.Writer = writer

		c.Next()

		fmt.Printf("[%s] [MIDDLEWARE] Key '%s' - PROCESSING COMPLETED (status: %d), caching result\n",
			time.Now().Format("15:04:05.000"), key, writer.statusCode)
		store.Finish(key, writer.statusCode, writer.body)
	}
}

// responseWriter wraps gin.ResponseWriter to capture response data
type responseWriter struct {
	gin.ResponseWriter
	statusCode int
	body       []byte
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriter) Write(data []byte) (int, error) {
	w.body = append(w.body, data...)
	return w.ResponseWriter.Write(data)
}
