package auth

import (
	"sync"
	"time"
)

type attempts struct {
	count       int
	windowStart time.Time
	lockedUntil time.Time
}

type Lockout struct {
	mu          sync.Mutex
	maxFailures int
	window      time.Duration
	duration    time.Duration
	byAddress   map[string]*attempts
	now         func() time.Time
}

func NewLockout(maxFailures int, window time.Duration, duration time.Duration) *Lockout {
	return &Lockout{maxFailures: maxFailures, window: window, duration: duration,
		byAddress: map[string]*attempts{}, now: time.Now}
}

func (l *Lockout) Locked(address string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	current, found := l.byAddress[address]
	return found && l.now().Before(current.lockedUntil)
}

func (l *Lockout) Fail(address string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	current, found := l.byAddress[address]
	if !found || now.Sub(current.windowStart) > l.window {
		if len(l.byAddress) >= maxTrackedKeys {
			l.byAddress = map[string]*attempts{}
		}
		current = &attempts{windowStart: now}
		l.byAddress[address] = current
	}
	current.count++
	if current.count >= l.maxFailures {
		current.lockedUntil = now.Add(l.duration)
	}
}
