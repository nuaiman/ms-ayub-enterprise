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
	clients  map[string]*Client
	interval time.Duration
	rate     int
	lastGC   time.Time
}

func NewRateLimiter(rate int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		clients:  make(map[string]*Client),
		rate:     rate,
		interval: interval,
		lastGC:   time.Now(),
	}
}

// gcLocked removes stale clients. Caller must hold rl.mu.
func (rl *RateLimiter) gcLocked(now time.Time) {
	// Run at most once per interval.
	if now.Sub(rl.lastGC) < rl.interval {
		return
	}
	cutoff := now.Add(-2 * rl.interval)
	for ip, c := range rl.clients {
		if c.lastRefil.Before(cutoff) {
			delete(rl.clients, ip)
		}
	}
	rl.lastGC = now
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

			now := time.Now()
			rl.gcLocked(now)

			c, exists := rl.clients[ip]
			if !exists {
				rl.clients[ip] = &Client{
					tokens:    rl.rate - 1,
					lastRefil: now,
				}
				next.ServeHTTP(w, r)
				return
			}

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
