package auth

import (
	"fmt"
	"slices"
	"strings"
)

type Scope string

const (
	ScopeObserve Scope = "observe"
	ScopeControl Scope = "control"
	ScopeText    Scope = "text"
	ScopeScreen  Scope = "screen"
)

var AllScopes = []Scope{ScopeObserve, ScopeControl, ScopeText, ScopeScreen}

func ParseScopes(list string) ([]Scope, error) {
	scopes := []Scope{}
	for raw := range strings.SplitSeq(list, ",") {
		scope := Scope(strings.TrimSpace(raw))
		if !slices.Contains(AllScopes, scope) {
			return nil, fmt.Errorf("unknown scope %q (allowed: observe, control, text, screen)", scope)
		}
		if !slices.Contains(scopes, scope) {
			scopes = append(scopes, scope)
		}
	}
	return scopes, nil
}

func HasScope(scopes []Scope, wanted Scope) bool {
	return slices.Contains(scopes, wanted)
}
