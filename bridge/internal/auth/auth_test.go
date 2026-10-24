package auth

import (
	"fmt"
	"testing"
	"time"
)

func TestTokenRoundTrip(t *testing.T) {
	plain, token, err := NewToken("claude", []Scope{ScopeObserve})
	if err != nil {
		t.Fatal(err)
	}
	store := NewStore([]Token{token})
	found, ok := store.Authenticate("Bearer " + plain)
	if !ok || found.Name != "claude" {
		t.Fatalf("valid token rejected")
	}
	for _, header := range []string{"", plain, "Bearer acb_wrong", "Basic " + plain, "Bearer " + plain + "x"} {
		_, ok := store.Authenticate(header)
		if ok {
			t.Fatalf("header %q must be rejected", header)
		}
	}
}

func TestTokenNameIsValidated(t *testing.T) {
	_, _, err := NewToken("../evil name", nil)
	if err == nil {
		t.Fatal("invalid name accepted")
	}
}

func TestParseScopes(t *testing.T) {
	scopes, err := ParseScopes("observe, control,observe")
	if err != nil || len(scopes) != 2 {
		t.Fatalf("unexpected scopes %v %v", scopes, err)
	}
	_, err = ParseScopes("observe,admin")
	if err == nil {
		t.Fatal("unknown scope accepted")
	}
}

func TestLimiter(t *testing.T) {
	now := time.Unix(0, 0)
	limiter := NewLimiter(1, 2)
	limiter.now = func() time.Time { return now }
	if !limiter.Allow("a") || !limiter.Allow("a") || limiter.Allow("a") {
		t.Fatal("burst not enforced")
	}
	now = now.Add(time.Second)
	if !limiter.Allow("a") {
		t.Fatal("bucket did not refill")
	}
}

func TestLockout(t *testing.T) {
	now := time.Unix(0, 0)
	lockout := NewLockout(3, time.Minute, time.Hour)
	lockout.now = func() time.Time { return now }
	for range 3 {
		lockout.Fail("10.0.0.2")
	}
	if !lockout.Locked("10.0.0.2") || lockout.Locked("10.0.0.3") {
		t.Fatal("lockout not applied per address")
	}
	now = now.Add(2 * time.Hour)
	if lockout.Locked("10.0.0.2") {
		t.Fatal("lockout never expires")
	}
}

func TestLockoutKeepsActiveLocksWhenFull(t *testing.T) {
	now := time.Unix(0, 0)
	lockout := NewLockout(1, time.Minute, time.Hour)
	lockout.now = func() time.Time { return now }
	lockout.Fail("attacker")
	for i := range maxTrackedKeys {
		lockout.Fail(fmt.Sprintf("10.0.%d.%d", i/256, i%256))
	}
	if !lockout.Locked("attacker") {
		t.Fatal("flooding the table must not lift an active lock")
	}
	if !lockout.Locked("192.168.1.50") {
		t.Fatal("a saturated table must refuse unknown addresses")
	}
}

func TestStoreReplaceRevokesTokens(t *testing.T) {
	plain, token, err := NewToken("phone", []Scope{ScopeControl})
	if err != nil {
		t.Fatal(err)
	}
	store := NewStore([]Token{token})
	store.Replace(nil)
	_, ok := store.Authenticate("Bearer " + plain)
	if ok {
		t.Fatal("revoked token still accepted")
	}
}

func TestLockoutRecoversOnceEntriesExpire(t *testing.T) {
	now := time.Unix(0, 0)
	lockout := NewLockout(1, time.Minute, time.Minute)
	lockout.now = func() time.Time { return now }
	for i := range maxTrackedKeys {
		lockout.Fail(fmt.Sprintf("10.1.%d.%d", i/256, i%256))
	}
	if !lockout.Locked("192.168.1.50") {
		t.Fatal("a table full of active locks must refuse unknown addresses")
	}
	now = now.Add(time.Hour)
	if lockout.Locked("192.168.1.50") {
		t.Fatal("expired entries must be pruned so unknown addresses recover")
	}
}
