package ratelimiter

import (
	"math"
	"time"
)

type RateLimitInfo struct {
	Limit         int
	Remaining     int
	RetryAfterSec int
}

type Limiter struct {
	config Config
	store  Store
}

func NewLimiter(config Config, store Store) *Limiter {
	return &Limiter{
		config: config,
		store:  store,
	}
}

func (l *Limiter) Allow(key string) (bool, RateLimitInfo) {

	total_burst := l.config.Rate + l.config.Burst

	prev_current, curr_current, windowStart, err := l.store.increment(key, l.config.Window)

	if err != nil {
		return false, RateLimitInfo{
			Limit:         l.config.Rate,
			Remaining:     0,
			RetryAfterSec: 0,
		}
	}

	now := time.Now()
	elapsed := now.Sub(windowStart)
	windowDuration := l.config.Window

	var fraction float64
	if windowDuration > 0 {
		fraction = float64(elapsed) / float64(windowDuration)
		if fraction > 1.0 {
			fraction = 1.0
		}
	}

	prevWeight := 1.0 - fraction
	weightedCount := float64(prev_current)*prevWeight + float64(curr_current)

	remaining := total_burst - int(math.Ceil(weightedCount))
	if remaining < 0 {
		remaining = 0
	}

	if weightedCount > float64(total_burst) {
		// Request denied — calculate retry-after
		// Estimate how long until enough of the previous window's weight
		// has decayed to bring the count below the limit.
		retryAfter := int(math.Ceil(windowDuration.Seconds() - elapsed.Seconds()))
		if retryAfter < 1 {
			retryAfter = 1
		}

		return false, RateLimitInfo{
			Limit:         total_burst,
			Remaining:     0,
			RetryAfterSec: retryAfter,
		}
	}

	return true, RateLimitInfo{
		Limit:     total_burst,
		Remaining: remaining,
	}

}

func (l *Limiter) IsWhiteListerd(ip string) bool {

	for _, ip_data := range l.config.IPWhiteList {
		if ip == ip_data {
			return true
		}
	}

	return false

}

func (l *Limiter) GetConfig() Config {
	return l.config
}
