package server

import (
	"context"
	"log/slog"
	"slices"
	"time"

	"github.com/sway-agency/tv/bridge/internal/auth"
	"github.com/sway-agency/tv/bridge/internal/config"
)

const reloadInterval = 2 * time.Second

func WatchTokens(ctx context.Context, path string, store *auth.Store, logger *slog.Logger) {
	ticker := time.NewTicker(reloadInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
		reloadTokens(path, store, logger)
	}
}

func reloadTokens(path string, store *auth.Store, logger *slog.Logger) {
	settings, err := config.Load(path)
	if err != nil {
		if store.Len() > 0 {
			logger.Error("configuration rejected, every token is refused until it is fixed", "error", err)
		}
		store.Replace(nil)
		return
	}
	if !slices.Equal(fingerprint(store.Tokens()), fingerprint(settings.Tokens)) {
		logger.Info("tokens reloaded", "count", len(settings.Tokens))
	}
	store.Replace(settings.Tokens)
}

func fingerprint(tokens []auth.Token) []string {
	hashes := make([]string, 0, len(tokens))
	for _, token := range tokens {
		hashes = append(hashes, token.Name+":"+token.SHA256+":"+scopeList(token.Scopes))
	}
	slices.Sort(hashes)
	return hashes
}

func scopeList(scopes []auth.Scope) string {
	joined := ""
	for _, scope := range scopes {
		joined += string(scope) + ","
	}
	return joined
}
