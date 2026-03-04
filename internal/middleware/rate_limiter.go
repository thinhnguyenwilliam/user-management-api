// user-management-api/internal/middleware/rate_limiter.go
package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	clients map[string]*client
	rps     rate.Limit
	burst   int
	ttl     time.Duration
	cleanup time.Duration
}

func NewRateLimiter(rps int, burst int, ttl time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*client),
		rps:     rate.Limit(rps),
		burst:   burst,
		ttl:     ttl,
		cleanup: time.Minute,
	}

	go rl.cleanupLoop()

	return rl
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := extractClientIP(ctx)

		limiter := rl.getLimiter(ip)

		if !limiter.Allow() {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":   "Too many requests",
				"message": "Rate limit exceeded. Please try again later.",
			})
			return
		}

		ctx.Next()
	}
}

func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	c, exists := rl.clients[ip]
	if !exists {
		limiter := rate.NewLimiter(rl.rps, rl.burst)
		rl.clients[ip] = &client{
			limiter:  limiter,
			lastSeen: time.Now(),
		}
		return limiter
	}

	c.lastSeen = time.Now()
	return c.limiter
}

func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.cleanup)
	defer ticker.Stop()

	for range ticker.C {
		rl.cleanupClients()
	}
}

func (rl *RateLimiter) cleanupClients() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	for ip, c := range rl.clients {
		if time.Since(c.lastSeen) > rl.ttl {
			delete(rl.clients, ip)
		}
	}
}

func extractClientIP(ctx *gin.Context) string {
	ip := ctx.ClientIP()
	if ip != "" {
		return ip
	}

	host, _, err := net.SplitHostPort(ctx.Request.RemoteAddr)
	if err != nil {
		return ctx.Request.RemoteAddr
	}

	return host
}
