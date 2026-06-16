package bybit

import (
	"time"
)

const minRequestInterval = 500 * time.Millisecond

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
}
