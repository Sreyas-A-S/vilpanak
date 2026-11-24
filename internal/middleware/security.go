package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Cross-Origin-Resource-Policy", "cross-origin")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self' * data: blob: 'unsafe-inline' 'unsafe-eval'; img-src * data: blob: 'unsafe-inline'; font-src * data:; style-src * 'unsafe-inline'; script-src * 'unsafe-inline' 'unsafe-eval'; connect-src *")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

type clientLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	clients   = make(map[string]*clientLimiter)
	clientsMu sync.Mutex
)

// RateLimiter limits requests per client IP
func RateLimiter(r rate.Limit, b int) gin.HandlerFunc {
	// Cleanup routine every 5 mins
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			clientsMu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 10*time.Minute {
					delete(clients, ip)
				}
			}
			clientsMu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		clientsMu.Lock()
		client, exists := clients[ip]
		if !exists {
			client = &clientLimiter{limiter: rate.NewLimiter(r, b)}
			clients[ip] = client
		}
		client.lastSeen = time.Now()
		limiter := client.limiter
		clientsMu.Unlock()

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please slow down.",
			})
			return
		}
		c.Next()
	}
}
