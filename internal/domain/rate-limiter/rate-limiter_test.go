package ratelimiter

import (
	"testing"
	"time"

	"github.com/ryancarlos88/go-rate-limiter/config"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiter_Allow(t *testing.T) {
	cfg := &config.Config{
		MaxIPRequestsPerSecond:    5,
		MaxTokenRequestsPerSecond: make(map[string]int),
		IPLockTimeoutInSeconds:    5,
		TokenLockTimeoutInSeconds: 1,
	}

	cfg.MaxTokenRequestsPerSecond["token1"] = 10

	rl := NewRateLimiter(cfg)

	t.Run("Allow set token (should overwrite ip settings)", func(t *testing.T) {
		m := make(map[int]bool)

		for i := 1; i < 51; i++ {
			time.Sleep(time.Millisecond * 100)
			m[i] = rl.Allow("token1", "")
		}

		assert.True(t, m[1])
		assert.False(t, m[11])
		assert.True(t, m[21])
		assert.False(t, m[31])
		assert.True(t, m[41])
	})

	t.Run("Allow unset token (should use ip settings)", func(t *testing.T) {
		n := make(map[int]bool)

		for i := 1; i < 71; i++ {
			time.Sleep(time.Millisecond * 100)
			n[i] = rl.Allow("", "")
		}

		assert.True(t, n[1])
		assert.False(t, n[6])
		assert.True(t, n[61])
	})

}

func TestRateLimiter_AllowIP(t *testing.T) {
	rl := NewRateLimiter(&config.Config{
		MaxIPRequestsPerSecond: 10,
		IPLockTimeoutInSeconds: 1,
	})
	m := make(map[int]bool)
	for i := 1; i < 41; i++ {
		time.Sleep(time.Millisecond * 100)
		m[i] = rl.AllowIP("")
	}
	assert.True(t, m[1])
	assert.True(t, m[4])
	assert.False(t, m[15])
	assert.False(t, m[19])
	assert.True(t, m[21])
	assert.False(t, m[32])
	assert.False(t, m[35])
}
