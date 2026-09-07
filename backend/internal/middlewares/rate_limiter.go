package middlewares

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type Client struct {
	tokens    int
	lastRefil time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	Clients  map[string]*Client
	interval time.Duration
	rate     int
}

func NewRateLimiter(rate int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		Clients:  make(map[string]*Client),
		rate:     rate,
		interval: interval,
	}
}

func (rl *RateLimiter) LimitRate(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				http.Error(w, "invalid client ip", http.StatusInternalServerError)
				return
			}

			rl.mu.Lock()
			defer rl.mu.Unlock()

			c, exists := rl.Clients[ip]
			if !exists {
				rl.Clients[ip] = &Client{
					tokens:    rl.rate - 1,
					lastRefil: time.Now(),
				}
				next.ServeHTTP(w, r)
				return
			}

			now := time.Now()
			if now.Sub(c.lastRefil) > rl.interval {
				c.tokens = rl.rate
				c.lastRefil = now
			}

			if c.tokens <= 0 {
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			c.tokens--

			next.ServeHTTP(w, r)
		})
}
