package main

import (
	"errors"
	"flag"
	"fmt"
	"path/filepath"
	"slices"

	"github.com/sway-agency/tv/bridge/internal/auth"
	"github.com/sway-agency/tv/bridge/internal/config"
)

func runToken(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: accueil-bridge token add|list|revoke")
	}
	flags := flag.NewFlagSet("token", flag.ContinueOnError)
	dir := flags.String("dir", defaultDir(), "configuration directory")
	name := flags.String("name", "", "token name, e.g. claude or phone")
	scopes := flags.String("scopes", "observe,control", "observe, control, text, screen")
	err := flags.Parse(args[1:])
	if err != nil {
		return err
	}
	path := filepath.Join(*dir, "config.json")
	settings, err := config.Load(path)
	if err != nil {
		return err
	}
	switch args[0] {
	case "add":
		return addToken(path, settings, *name, *scopes)
	case "list":
		for _, token := range settings.Tokens {
			fmt.Printf("%-20s %-30v created %s\n", token.Name, token.Scopes, token.Created.Format("2006-01-02"))
		}
		return nil
	case "revoke":
		before := len(settings.Tokens)
		settings.Tokens = slices.DeleteFunc(settings.Tokens, func(token auth.Token) bool { return token.Name == *name })
		if len(settings.Tokens) == before {
			return fmt.Errorf("no token named %q", *name)
		}
		err = config.Save(path, settings)
		if err != nil {
			return err
		}
		fmt.Printf("Revoked %s. A running bridge refuses it within 2 seconds.\n", *name)
		return nil
	}
	return fmt.Errorf("unknown token command %q", args[0])
}

func addToken(path string, settings config.Config, name string, scopeList string) error {
	scopes, err := auth.ParseScopes(scopeList)
	if err != nil {
		return err
	}
	if slices.ContainsFunc(settings.Tokens, func(token auth.Token) bool { return token.Name == name }) {
		return fmt.Errorf("a token named %q exists, revoke it first", name)
	}
	plain, token, err := auth.NewToken(name, scopes)
	if err != nil {
		return err
	}
	settings.Tokens = append(settings.Tokens, token)
	err = config.Save(path, settings)
	if err != nil {
		return err
	}
	fmt.Printf("Token for %s (shown once, store it in the agent secrets):\n%s\n", name, plain)
	return nil
}
