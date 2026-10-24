package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateListen(t *testing.T) {
	for _, listen := range []string{"0.0.0.0:8443", "[::]:8443", ":8443", "8.8.8.8:8443", "localhost:8443"} {
		if validateListen(listen) == nil {
			t.Fatalf("listen %q must be refused", listen)
		}
	}
	for _, listen := range []string{"127.0.0.1:8443", "192.168.1.10:8443", "10.8.0.1:8443"} {
		if validateListen(listen) != nil {
			t.Fatalf("listen %q must be accepted", listen)
		}
	}
}

func TestCheckPrivate(t *testing.T) {
	dir := t.TempDir()
	err := os.Chmod(dir, 0o700)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "config.json")
	err = os.WriteFile(path, []byte("{}"), 0o644)
	if err != nil {
		t.Fatal(err)
	}
	if CheckPrivate(path) == nil {
		t.Fatal("world readable config accepted")
	}
	err = os.Chmod(path, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	if CheckPrivate(path) != nil {
		t.Fatal("private config refused")
	}
}
