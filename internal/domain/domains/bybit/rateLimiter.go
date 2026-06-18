package bybit

import (
	"github.com/shatylos/trader/tools/logger"
	"time"
)

const minRequestInterval = 1 * time.Second

type rateLimiter struct {
	lastRequestTime time.Time
}

func (r *rateLimiter) throttle() {
	if r.lastRequestTime.IsZero() {
		r.lastRequestTime = time.Now()
		return
	}
	elapsed := time.Since(r.lastRequestTime)
	if elapsed < minRequestInterval {
		time.Sleep(minRequestInterval - elapsed)
	}
	r.lastRequestTime = time.Now()
	logger.Info("throttle")
}
