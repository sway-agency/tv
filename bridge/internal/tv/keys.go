package tv

import (
	"fmt"
	"sort"
)

type Key struct {
	Name string
	Code int
}

const MaxKeyRepeat = 20

var keyCodes = map[string]int{
	"up":           19,
	"down":         20,
	"left":         21,
	"right":        22,
	"select":       23,
	"back":         4,
	"home":         3,
	"menu":         82,
	"play_pause":   85,
	"play":         126,
	"pause":        127,
	"stop":         86,
	"next":         87,
	"previous":     88,
	"rewind":       89,
	"fast_forward": 90,
	"volume_up":    24,
	"volume_down":  25,
	"mute":         164,
	"channel_up":   166,
	"channel_down": 167,
	"info":         165,
	"guide":        172,
}

func KeyNames() []string {
	names := make([]string, 0, len(keyCodes))
	for name := range keyCodes {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func ParseKey(name string) (Key, error) {
	code, ok := keyCodes[name]
	if !ok {
		return Key{}, fmt.Errorf("unknown key %q", name)
	}
	return Key{Name: name, Code: code}, nil
}
