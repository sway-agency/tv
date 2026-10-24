package server

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/sway-agency/tv/bridge/internal/auth"
	"github.com/sway-agency/tv/bridge/internal/config"
)

func TestReloadTokensFailsClosedAndRecovers(t *testing.T) {
	dir := t.TempDir()
	err := os.Chmod(dir, 0o700)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "config.json")
	plain, token, err := auth.NewToken("agent", []auth.Scope{auth.ScopeObserve})
	if err != nil {
		t.Fatal(err)
	}
	settings := config.Config{
		Listen: "127.0.0.1:8443", Hosts: []string{"127.0.0.1"}, ADBPath: "/usr/bin/adb", TVSerial: "10.0.0.2:5555",
		CertFile: "/c.pem", KeyFile: "/k.pem", RateLimitPerSecond: 1, RateLimitBurst: 1, Tokens: []auth.Token{token},
	}
	err = config.Save(path, settings)
	if err != nil {
		t.Fatal(err)
	}
	store := auth.NewStore(nil)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	reloadTokens(path, store, logger)
	_, ok := store.Authenticate("Bearer " + plain)
	if !ok {
		t.Fatal("token from a valid configuration refused")
	}
	err = os.Chmod(path, 0o644)
	if err != nil {
		t.Fatal(err)
	}
	reloadTokens(path, store, logger)
	_, ok = store.Authenticate("Bearer " + plain)
	if ok {
		t.Fatal("a configuration readable by others must refuse every token")
	}
}

func TestLockoutKeyGroupsIPv6Prefix(t *testing.T) {
	if lockoutKey("2001:db8:1:2:aaaa::1") != lockoutKey("2001:db8:1:2:bbbb::9") {
		t.Fatal("addresses of the same /64 must share a lockout")
	}
	if lockoutKey("192.168.1.20") != "192.168.1.20" {
		t.Fatal("IPv4 addresses are kept as is")
	}
}
