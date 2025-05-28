package ratelimiter

import (
	"time"
)

// NewRateLimiter initializes a new RateLimiter with a specified request limit and cooldown duration.
func NewRateLimiter(limit int, cooldown time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string]*RateLimitData),
		limit:    limit,
		cooldown: cooldown,
	}
}

// Allow checks if a request is allowed for the given key (e.g., IP or username) and returns a boolean and cooldown duration.
func (rl *RateLimiter) Allow(key string) (bool, time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	data, exists := rl.requests[key]

	if !exists || rl.isCooldownExpired(data, now) {
		rl.resetRateLimitData(key, now)
		return true, 0
	}

	if rl.isWithinTimeWindow(data) {
		return rl.handleWithinTimeWindow(data, now)
	}

	rl.resetRateLimitData(key, now)
	return true, 0
}

func (rl *RateLimiter) isCooldownExpired(data *RateLimitData, now time.Time) bool {
	return data == nil || (now.After(data.Cooldown) && !data.Cooldown.IsZero())
}

func (rl *RateLimiter) resetRateLimitData(key string, now time.Time) {
	rl.requests[key] = &RateLimitData{
		Count:      1,
		InitAccess: now,
		LastAccess: now,
		Cooldown:   time.Time{}, // Reset cooldown
	}
}

func (rl *RateLimiter) isWithinTimeWindow(data *RateLimitData) bool {
	return data.LastAccess.Sub(data.InitAccess) <= time.Minute
}

// handleWithinTimeWindow processes requests within the rate limit time window and updates the rate limit data.
func (rl *RateLimiter) handleWithinTimeWindow(data *RateLimitData, now time.Time) (bool, time.Duration) {
	data.Count++
	data.LastAccess = now

	if data.Count > rl.limit {
		if data.Cooldown.IsZero() {
			data.Cooldown = now.Add(rl.cooldown)
		}
		return false, data.Cooldown.Sub(now)
	}

	return true, 0
}
