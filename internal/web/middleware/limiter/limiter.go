package limiter

import (
	"net/http"
	"strings"

	ratelimiter "github.com/ryancarlos88/go-rate-limiter/internal/domain/rate-limiter"
)

const (
	RateError = "you have reached the maximum number of requests or actions allowed within a certain time frame"
)

func LimitMiddleware(next http.Handler, limiter *ratelimiter.RateLimiter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		t := ""
		// get the bearer token from the request
		token := strings.Split(r.Header.Get("Authorization"), "API_KEY ")
		if len(token) > 1 {
			t = token[1]
		}
		// Check if the request is allowed
		if limiter.Allow(t, ip) {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, RateError, http.StatusTooManyRequests)
		}
	})
}
