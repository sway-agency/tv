package auth

import (
	"crypto/subtle"
	"strings"
	"sync"
)

type Store struct {
	mu     sync.RWMutex
	tokens []Token
}

func NewStore(tokens []Token) *Store {
	return &Store{tokens: tokens}
}

func (s *Store) Replace(tokens []Token) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens = tokens
}

func (s *Store) Authenticate(authorization string) (Token, bool) {
	plain, found := strings.CutPrefix(authorization, "Bearer ")
	if !found || !strings.HasPrefix(plain, tokenPrefix) || len(plain) > 128 {
		return Token{}, false
	}
	presented := []byte(Hash(plain))
	s.mu.RLock()
	defer s.mu.RUnlock()
	var matched Token
	ok := false
	for _, token := range s.tokens {
		if subtle.ConstantTimeCompare(presented, []byte(token.SHA256)) == 1 {
			matched, ok = token, true
		}
	}
	return matched, ok
}

func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.tokens)
}
