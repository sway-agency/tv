package tv

import (
	"fmt"
	"regexp"
)

const MaxTextLength = 64

var (
	packagePattern   = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*(\.[a-zA-Z][a-zA-Z0-9_]*)+$`)
	componentPattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_.]*/[a-zA-Z0-9_.$]+$`)
	textPattern      = regexp.MustCompile(`^[a-zA-Z0-9 ]+$`)
)

func ValidatePackage(name string) error {
	if len(name) > 255 || !packagePattern.MatchString(name) {
		return fmt.Errorf("invalid package name %q", name)
	}
	return nil
}

func validateComponent(component string) error {
	if len(component) > 512 || !componentPattern.MatchString(component) {
		return fmt.Errorf("invalid activity component")
	}
	return nil
}

func ValidateText(text string) error {
	if len(text) == 0 || len(text) > MaxTextLength || !textPattern.MatchString(text) {
		return fmt.Errorf("text must be 1 to %d letters, digits or spaces", MaxTextLength)
	}
	return nil
}
