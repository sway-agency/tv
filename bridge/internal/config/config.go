package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sway-agency/tv/bridge/internal/auth"
)

type Config struct {
	Listen             string       `json:"listen"`
	Hosts              []string     `json:"hosts"`
	AllowedOrigins     []string     `json:"allowed_origins"`
	CertFile           string       `json:"cert_file"`
	KeyFile            string       `json:"key_file"`
	ClientCAFile       string       `json:"client_ca_file"`
	RequireClientCert  bool         `json:"require_client_cert"`
	ADBPath            string       `json:"adb_path"`
	TVSerial           string       `json:"tv_serial"`
	RateLimitPerSecond float64      `json:"rate_limit_per_second"`
	RateLimitBurst     int          `json:"rate_limit_burst"`
	Tokens             []auth.Token `json:"tokens"`
}

func Load(path string) (Config, error) {
	err := CheckPrivate(path)
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var loaded Config
	err = decoder.Decode(&loaded)
	if err != nil {
		return Config{}, fmt.Errorf("config %s: %w", path, err)
	}
	return loaded, loaded.Validate()
}

func Save(path string, value Config) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	temporary, err := os.CreateTemp(filepath.Dir(path), ".config-*.json")
	if err != nil {
		return err
	}
	defer os.Remove(temporary.Name())
	err = temporary.Chmod(0o600)
	if err != nil {
		return err
	}
	_, err = temporary.Write(append(data, '\n'))
	if err != nil {
		return err
	}
	err = temporary.Close()
	if err != nil {
		return err
	}
	return os.Rename(temporary.Name(), path)
}
