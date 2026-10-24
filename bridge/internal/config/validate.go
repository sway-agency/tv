package config

import (
	"errors"
	"fmt"
	"net"
	"net/netip"
	"path/filepath"
	"regexp"
)

var serialPattern = regexp.MustCompile(`^[a-zA-Z0-9.:_-]{1,64}$`)

func (c Config) Validate() error {
	err := validateListen(c.Listen)
	if err != nil {
		return err
	}
	if len(c.Hosts) == 0 {
		return errors.New("hosts must list the names or IP addresses clients use to reach the bridge")
	}
	for _, path := range []string{c.CertFile, c.KeyFile, c.ADBPath} {
		if !filepath.IsAbs(path) {
			return fmt.Errorf("path %q must be absolute", path)
		}
	}
	if c.RequireClientCert && !filepath.IsAbs(c.ClientCAFile) {
		return errors.New("require_client_cert needs an absolute client_ca_file")
	}
	if !serialPattern.MatchString(c.TVSerial) {
		return errors.New("tv_serial must look like 192.168.1.20:5555")
	}
	if c.RateLimitPerSecond <= 0 || c.RateLimitBurst < 1 {
		return errors.New("rate_limit_per_second and rate_limit_burst must be positive")
	}
	return nil
}

func validateListen(listen string) error {
	host, _, err := net.SplitHostPort(listen)
	if err != nil {
		return fmt.Errorf("listen must be ip:port: %w", err)
	}
	address, err := netip.ParseAddr(host)
	if err != nil {
		return errors.New("listen must use an explicit IP address, such as the LAN or WireGuard address")
	}
	if address.IsUnspecified() {
		return errors.New("listen must not be 0.0.0.0 or ::, bind to the LAN or WireGuard address only")
	}
	if !address.IsLoopback() && !address.IsPrivate() && !address.IsLinkLocalUnicast() {
		return errors.New("listen must be a loopback or private address, never a public one")
	}
	return nil
}
