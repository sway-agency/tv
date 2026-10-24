package server

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/sway-agency/tv/bridge/internal/auth"
	"github.com/sway-agency/tv/bridge/internal/config"
)

const reloadInterval = 2 * time.Second

func WatchTokens(ctx context.Context, path string, store *auth.Store, logger *slog.Logger) {
	lastChange := modified(path)
	ticker := time.NewTicker(reloadInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
		change := modified(path)
		if change.Equal(lastChange) {
			continue
		}
		lastChange = change
		settings, err := config.Load(path)
		if err != nil {
			store.Replace(nil)
			logger.Error("configuration rejected, every token is refused until it is fixed", "error", err)
			continue
		}
		store.Replace(settings.Tokens)
		logger.Info("tokens reloaded", "count", len(settings.Tokens))
	}
}

func modified(path string) time.Time {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}
	}
	return info.ModTime()
}
