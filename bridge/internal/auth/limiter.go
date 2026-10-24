package auth

import (
	"sync"
	"time"
)

const maxTrackedKeys = 4096

type bucket struct {
	tokens float64
	last   time.Time
}

type Limiter struct {
	mu      sync.Mutex
	rate    float64
	burst   float64
	buckets map[string]*bucket
	now     func() time.Time
}

func NewLimiter(perSecond float64, burst int) *Limiter {
	return &Limiter{rate: perSecond, burst: float64(burst), buckets: map[string]*bucket{}, now: time.Now}
}

func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	current, found := l.buckets[key]
	if !found {
		if len(l.buckets) >= maxTrackedKeys {
			l.buckets = map[string]*bucket{}
		}
		current = &bucket{tokens: l.burst, last: now}
		l.buckets[key] = current
	}
	current.tokens = min(l.burst, current.tokens+now.Sub(current.last).Seconds()*l.rate)
	current.last = now
	if current.tokens < 1 {
		return false
	}
	current.tokens--
	return true
}
