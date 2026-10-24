package tv

import (
	"slices"
	"testing"
)

func TestParsePackages(t *testing.T) {
	output := "Activity #0:\n  packageName=com.netflix.ninja\n  packageName=com.netflix.ninja\n  packageName=bad;name\n  packageName=net.oqee.androidtv.store\n"
	names := parsePackages(output)
	if !slices.Equal(names, []string{"com.netflix.ninja", "net.oqee.androidtv.store"}) {
		t.Fatalf("unexpected packages %v", names)
	}
}

func TestParseForeground(t *testing.T) {
	line := "    mResumedActivity: ActivityRecord{921eb1d u0 fr.sway.accueil/.HomeActivity t2761}"
	if parseForeground(line) != "fr.sway.accueil" || parseForeground("nothing") != "" {
		t.Fatal("foreground parsing failed")
	}
}
