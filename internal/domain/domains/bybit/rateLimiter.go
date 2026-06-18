package bybit

import (
	"sync"
	"time"
)

const minRequestInterval = 500 * time.Millisecond

var (
	throttleMu      sync.Mutex
	lastRequestTime time.Time
)

func throttle() {
	throttleMu.Lock()
	defer throttleMu.Unlock()

	if !lastRequestTime.IsZero() {
		elapsed := time.Since(lastRequestTime)
		if elapsed < minRequestInterval {
			time.Sleep(minRequestInterval - elapsed)
		}
	}
	lastRequestTime = time.Now()
}
