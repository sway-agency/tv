package auth

import (
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
