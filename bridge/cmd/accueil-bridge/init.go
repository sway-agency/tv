package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sway-agency/tv/bridge/internal/auth"
	"github.com/sway-agency/tv/bridge/internal/config"
	"github.com/sway-agency/tv/bridge/internal/pki"
)

func runInit(args []string) error {
	flags := flag.NewFlagSet("init", flag.ContinueOnError)
	dir := flags.String("dir", defaultDir(), "configuration directory")
	listen := flags.String("listen", "", "LAN or WireGuard ip:port to listen on, e.g. 192.168.1.10:8443")
	serial := flags.String("tv", "", "TV adb address, e.g. 192.168.1.20:5555")
	adbPath := flags.String("adb", "adb", "adb binary")
	var hosts listFlag
	flags.Var(&hosts, "host", "name or IP clients use to reach the bridge (repeatable)")
	err := flags.Parse(args)
	if err != nil {
		return err
	}
	resolvedADB, err := exec.LookPath(*adbPath)
	if err != nil {
		return fmt.Errorf("adb not found: %w", err)
	}
	absoluteADB, err := filepath.Abs(resolvedADB)
	if err != nil {
		return err
	}
	settings := config.Config{
		Listen: *listen, Hosts: hosts, TVSerial: *serial, ADBPath: absoluteADB,
		CertFile: filepath.Join(*dir, "server.pem"), KeyFile: filepath.Join(*dir, "server-key.pem"),
		ClientCAFile: filepath.Join(*dir, "ca.pem"), RateLimitPerSecond: 5, RateLimitBurst: 20,
		Tokens: []auth.Token{},
	}
	err = settings.Validate()
	if err != nil {
		return err
	}
	return writeInitialFiles(*dir, settings)
}

func writeInitialFiles(dir string, settings config.Config) error {
	path := filepath.Join(dir, "config.json")
	_, err := os.Stat(path)
	if err == nil {
		return errors.New(path + " already exists")
	}
	err = os.MkdirAll(dir, 0o700)
	if err != nil {
		return err
	}
	err = os.Chmod(dir, 0o700)
	if err != nil {
		return err
	}
	authority, err := pki.NewAuthority(dir)
	if err != nil {
		return err
	}
	err = authority.IssueServer(dir, settings.Hosts)
	if err != nil {
		return err
	}
	err = config.Save(path, settings)
	if err != nil {
		return err
	}
	fmt.Printf("Created %s\nTrust %s on each client, then add a token: accueil-bridge token add --name <agent>\n",
		path, filepath.Join(dir, "ca.pem"))
	return nil
}
