package ratelimiter

import (
	"sync"
	"time"
)

type RateLimiter struct {
	mu       sync.Mutex
	requests map[string]*RateLimitData
	limit    int
	cooldown time.Duration
}

type RateLimitData struct {
	Count      int
	InitAccess time.Time
	LastAccess time.Time
	Cooldown   time.Time
}
