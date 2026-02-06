package passport

import (
	"crypto/rand"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/leeprovoost/go-rest-api-template/pkg/status"
	"golang.org/x/time/rate"
)

// requestID reads or generates a unique request ID and sets it on the response.
func requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = generateID()
		}
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r)
	})
}

func generateID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// cors returns middleware that adds CORS headers and handles preflight requests.
func cors(allowedOrigins string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// rateLimiter implements per-IP rate limiting using a token bucket algorithm.
// Note: This is suitable for single-instance deployments. For distributed
// systems, use an external store like Redis.
type rateLimiter struct {
	clients map[string]*rate.Limiter
	mu      sync.Mutex
	r       rate.Limit
	burst   int
}

func newRateLimiter(rps float64, burst int) *rateLimiter {
	return &rateLimiter{
		clients: make(map[string]*rate.Limiter),
		r:       rate.Limit(rps),
		burst:   burst,
	}
}

func (rl *rateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	limiter, exists := rl.clients[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.r, rl.burst)
		rl.clients[ip] = limiter
	}
	return limiter
}

func (rl *rateLimiter) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := clientIP(r)
		if !rl.getLimiter(ip).Allow() {
			respond(w, http.StatusTooManyRequests, status.Response{
				Status:  strconv.Itoa(http.StatusTooManyRequests),
				Message: "rate limit exceeded",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func clientIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return forwarded
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
