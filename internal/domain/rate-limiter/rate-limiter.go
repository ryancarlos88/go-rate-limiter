package ratelimiter

import (
	"sync"
	"time"

	"github.com/ryancarlos88/go-rate-limiter/config"
)

type RateLimiter struct {
	ips              map[string]*RateLimit
	tokens           map[string]*RateLimit
	mutex            sync.Mutex
	maxIPRequests    int
	ipLockout        time.Duration
	tokenLockout     time.Duration
	maxTokenRequests map[string]int
}

type RateLimit struct {
	count     int
	lastReset time.Time
}

func NewRateLimiter(cfg *config.Config) *RateLimiter {
	return &RateLimiter{
		ips:              make(map[string]*RateLimit),
		tokens:           make(map[string]*RateLimit),
		maxIPRequests:    cfg.MaxIPRequestsPerSecond,
		ipLockout:        time.Duration(cfg.IPLockTimeoutInSeconds) * time.Second,
		tokenLockout:     time.Duration(cfg.TokenLockTimeoutInSeconds) * time.Second,
		maxTokenRequests: cfg.MaxTokenRequestsPerSecond,
	}
}

func (rl *RateLimiter) Allow(token, ip string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	_, ok := rl.maxTokenRequests[token]
	if !ok {
		return rl.AllowIP(ip)
	}
	return rl.AllowToken(token)

}

func (rl *RateLimiter) AllowIP(ip string) bool {
	return rl.allow(ip, rl.ips, rl.maxIPRequests, rl.ipLockout)
}

func (rl *RateLimiter) AllowToken(token string) bool {
	return rl.allow(token, rl.tokens, rl.maxTokenRequests[token], rl.tokenLockout)
}

func (rl *RateLimiter) allow(key string, keys map[string]*RateLimit, requestLimit int, lockout time.Duration) bool {
	// Check if the label is already in the map
	rateLimit, ok := keys[key]
	if !ok {
		// label is not in the map, create a new rate limit
		rateLimit = &RateLimit{
			count:     0,
			lastReset: time.Now(),
		}
		keys[key] = rateLimit
	}

	// Check if the rate limit has expired
	if time.Since(rateLimit.lastReset) > lockout+time.Second {
		rateLimit.count = 0
		rateLimit.lastReset = time.Now()
	}

	// Check if the rate limit has been exceeded
	if rateLimit.count >= requestLimit {
		return false
	}

	// Increment the count and allow the request
	rateLimit.count++
	return true
}
