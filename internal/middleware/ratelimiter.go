package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// client tracks a per-IP rate limiter and the last time that IP was seen.
type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiter is a per-IP token-bucket rate limiter.
// Inactive clients are evicted after 3 minutes to prevent unbounded memory growth.
type RateLimiter struct {
	mu      sync.Mutex
	clients map[string]*client
	rate    int
	burst   int
}

// NewRateLimiter creates a RateLimiter that allows r requests per second with a
// burst capacity of burst. A background goroutine evicts idle clients every minute.
func NewRateLimiter(r, burst int) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*client),
		rate:    r,
		burst:   burst,
	}

	go func() {
		for {
			time.Sleep(time.Minute)
			rl.mu.Lock()
			for ip, c := range rl.clients {
				if time.Since(c.lastSeen) > 3*time.Minute {
					delete(rl.clients, ip)
				}
			}
			rl.mu.Unlock()
		}
	}()

	return rl
}

// getClient returns the rate.Limiter for the given IP, creating one if needed.
func (rl *RateLimiter) getClient(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if _, ok := rl.clients[ip]; !ok {
		rl.clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(rl.rate), rl.burst)}
	}
	rl.clients[ip].lastSeen = time.Now()
	return rl.clients[ip].limiter
}

// Limit is a chi-compatible middleware that enforces the per-IP rate limit.
// Returns 429 Too Many Requests when the token bucket is empty.
func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		if !rl.getClient(ip).Allow() {
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
