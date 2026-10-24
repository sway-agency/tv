package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"regexp"
	"time"
)

const tokenPrefix = "acb_"

var namePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,31}$`)

type Token struct {
	Name    string    `json:"name"`
	SHA256  string    `json:"sha256"`
	Scopes  []Scope   `json:"scopes"`
	Created time.Time `json:"created"`
}

func NewToken(name string, scopes []Scope) (string, Token, error) {
	if !namePattern.MatchString(name) {
		return "", Token{}, fmt.Errorf("token name must be 1 to 32 lowercase letters, digits or hyphens")
	}
	secret := make([]byte, 32)
	_, err := rand.Read(secret)
	if err != nil {
		return "", Token{}, err
	}
	plain := tokenPrefix + base64.RawURLEncoding.EncodeToString(secret)
	return plain, Token{Name: name, SHA256: Hash(plain), Scopes: scopes, Created: time.Now().UTC()}, nil
}

func Hash(plain string) string {
	sum := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(sum[:])
}
