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
	if !found {
		return len(l.byAddress) >= maxTrackedKeys
	}
	return l.now().Before(current.lockedUntil)
}

func (l *Lockout) Fail(address string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	current, found := l.byAddress[address]
	if !found || now.Sub(current.windowStart) > l.window {
		if !found && len(l.byAddress) >= maxTrackedKeys {
			l.pruneExpired(now)
		}
		if !found && len(l.byAddress) >= maxTrackedKeys {
			return
		}
		current = &attempts{windowStart: now, lockedUntil: currentLock(current)}
		l.byAddress[address] = current
	}
	current.count++
	if current.count >= l.maxFailures {
		current.lockedUntil = now.Add(l.duration)
	}
}

func (l *Lockout) pruneExpired(now time.Time) {
	for address, entry := range l.byAddress {
		if now.Sub(entry.windowStart) > l.window && !now.Before(entry.lockedUntil) {
			delete(l.byAddress, address)
		}
	}
}

func currentLock(previous *attempts) time.Time {
	if previous == nil {
		return time.Time{}
	}
	return previous.lockedUntil
}
