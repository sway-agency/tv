package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
)

var version = "dev"

const usage = `accueil-bridge: secure MCP and HTTP bridge to control an Android TV

Usage:
  accueil-bridge init --listen <ip:port> --host <name-or-ip>... --tv <ip:5555>
  accueil-bridge token add --name <name> [--scopes observe,control]
  accueil-bridge token list
  accueil-bridge token revoke --name <name>
  accueil-bridge cert client --name <name>
  accueil-bridge serve
  accueil-bridge stdio [--scopes observe,control]
  accueil-bridge version

Every command accepts --dir (default ~/.config/accueil-bridge).
`

const defaultMemoryLimit = 64 << 20

func main() {
	if os.Getenv("GOMEMLIMIT") == "" {
		debug.SetMemoryLimit(defaultMemoryLimit)
	}
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	err := run(os.Args[1], os.Args[2:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(command string, args []string) error {
	switch command {
	case "init":
		return runInit(args)
	case "token":
		return runToken(args)
	case "cert":
		return runCert(args)
	case "serve":
		return runServe(args)
	case "stdio":
		return runStdio(args)
	case "version":
		fmt.Println(version)
		return nil
	}
	fmt.Fprint(os.Stderr, usage)
	return fmt.Errorf("unknown command %q", command)
}

func defaultDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".accueil-bridge"
	}
	return filepath.Join(home, ".config", "accueil-bridge")
}

type listFlag []string

func (l *listFlag) String() string {
	return fmt.Sprint(*l)
}

func (l *listFlag) Set(value string) error {
	*l = append(*l, value)
	return nil
}
