package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/Yogdunana/yogduoj/backend/internal/config"
	"github.com/Yogdunana/yogduoj/backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiterMiddleware struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rate     rate.Limit
	burst    int
}

func NewRateLimiterMiddleware(cfg *config.RateLimitConfig) *RateLimiterMiddleware {
	rl := &RateLimiterMiddleware{
		visitors: make(map[string]*visitor),
		rate:     rate.Limit(cfg.RequestsPerMinute) / 60.0,
		burst:    cfg.Burst,
	}

	// Cleanup stale visitors periodically
	go rl.cleanupStaleVisitors()

	return rl
}

func (rl *RateLimiterMiddleware) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if v, exists := rl.visitors[ip]; exists {
		v.lastSeen = time.Now()
		return v.limiter
	}

	limiter := rate.NewLimiter(rl.rate, rl.burst)
	rl.visitors[ip] = &visitor{limiter: limiter, lastSeen: time.Now()}
	return limiter
}

func (rl *RateLimiterMiddleware) cleanupStaleVisitors() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiterMiddleware) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := rl.getLimiter(ip)

		if !limiter.Allow() {
			response.Error(c, http.StatusTooManyRequests, "rate limit exceeded")
			c.Abort()
			return
		}

		c.Next()
	}
}
