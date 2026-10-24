package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/sway-agency/tv/bridge/internal/auth"
	"github.com/sway-agency/tv/bridge/internal/config"
	"github.com/sway-agency/tv/bridge/internal/mcp"
	"github.com/sway-agency/tv/bridge/internal/pki"
	"github.com/sway-agency/tv/bridge/internal/server"
	"github.com/sway-agency/tv/bridge/internal/tv"
)

func runServe(args []string) error {
	flags := flag.NewFlagSet("serve", flag.ContinueOnError)
	dir := flags.String("dir", defaultDir(), "configuration directory")
	err := flags.Parse(args)
	if err != nil {
		return err
	}
	path := filepath.Join(*dir, "config.json")
	settings, err := config.Load(path)
	if err != nil {
		return err
	}
	if len(settings.Tokens) == 0 {
		return errors.New("no token configured, run: accueil-bridge token add --name <agent>")
	}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	dispatcher := mcp.NewDispatcher(tv.NewADB(settings.ADBPath, settings.TVSerial), version)
	store := auth.NewStore(settings.Tokens)
	httpServer, err := server.New(settings, dispatcher, store, logger)
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go server.WatchTokens(ctx, path, store, logger)
	go func() {
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = httpServer.Shutdown(shutdown)
	}()
	logger.Info("listening", "address", settings.Listen, "client_certificates", settings.RequireClientCert)
	err = httpServer.ListenAndServeTLS("", "")
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func runStdio(args []string) error {
	flags := flag.NewFlagSet("stdio", flag.ContinueOnError)
	dir := flags.String("dir", defaultDir(), "configuration directory")
	scopeList := flags.String("scopes", "observe,control", "observe, control, text, screen")
	err := flags.Parse(args)
	if err != nil {
		return err
	}
	scopes, err := auth.ParseScopes(*scopeList)
	if err != nil {
		return err
	}
	settings, err := config.Load(filepath.Join(*dir, "config.json"))
	if err != nil {
		return err
	}
	dispatcher := mcp.NewDispatcher(tv.NewADB(settings.ADBPath, settings.TVSerial), version)
	return mcp.ServeStdio(context.Background(), os.Stdin, os.Stdout, dispatcher, mcp.Caller{Name: "stdio", Scopes: scopes})
}

func runCert(args []string) error {
	if len(args) == 0 || args[0] != "client" {
		return errors.New("usage: accueil-bridge cert client --name <name>")
	}
	flags := flag.NewFlagSet("cert", flag.ContinueOnError)
	dir := flags.String("dir", defaultDir(), "configuration directory")
	name := flags.String("name", "", "client name")
	err := flags.Parse(args[1:])
	if err != nil {
		return err
	}
	authority, err := pki.LoadAuthority(*dir)
	if err != nil {
		return err
	}
	err = authority.IssueClient(*dir, *name)
	if err != nil {
		return err
	}
	fmt.Printf("Created client-%s.pem and client-%s-key.pem in %s. Set require_client_cert to true to enforce them.\n", *name, *name, *dir)
	return nil
}
