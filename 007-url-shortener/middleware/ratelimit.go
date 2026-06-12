package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

// Token bucket per client IP: each client gets `burst` tokens, refilled at
// `rps` per second; each request spends one. Out of tokens → 429.
type bucket struct {
	tokens   float64
	lastSeen time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	buckets map[string]*bucket
	rps     float64
	burst   float64
}

func NewRateLimiter(rps float64, burst int) *RateLimiter {
	rl := &RateLimiter{
		buckets: make(map[string]*bucket),
		rps:     rps,
		burst:   float64(burst),
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	b, ok := rl.buckets[ip]
	if !ok {
		rl.buckets[ip] = &bucket{tokens: rl.burst - 1, lastSeen: now}
		return true
	}

	elapsed := now.Sub(b.lastSeen).Seconds()
	b.tokens = min(rl.burst, b.tokens+elapsed*rl.rps)
	b.lastSeen = now

	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

// cleanup evicts idle buckets so the map doesn't grow forever.
func (rl *RateLimiter) cleanup() {
	for range time.Tick(time.Minute) {
		rl.mu.Lock()
		for ip, b := range rl.buckets {
			if time.Since(b.lastSeen) > 3*time.Minute {
				delete(rl.buckets, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}
		if !rl.allow(ip) {
			w.Header().Set("Retry-After", "1")
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
