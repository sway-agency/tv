package auth

import (
	"crypto/subtle"
	"strings"
)

type Store struct {
	tokens []Token
}

func NewStore(tokens []Token) *Store {
	return &Store{tokens: tokens}
}

func (s *Store) Authenticate(authorization string) (Token, bool) {
	plain, found := strings.CutPrefix(authorization, "Bearer ")
	if !found || !strings.HasPrefix(plain, tokenPrefix) || len(plain) > 128 {
		return Token{}, false
	}
	presented := []byte(Hash(plain))
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
	return len(s.tokens)
}
