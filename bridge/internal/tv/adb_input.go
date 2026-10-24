package tv

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

const (
	keycodeWakeup = 224
	keycodeSleep  = 223
)

func (a *ADB) PressKey(ctx context.Context, key Key, repeat int) error {
	if repeat < 1 || repeat > MaxKeyRepeat {
		return fmt.Errorf("repeat must be between 1 and %d", MaxKeyRepeat)
	}
	args := []string{"input", "keyevent"}
	for range repeat {
		args = append(args, strconv.Itoa(key.Code))
	}
	_, err := a.shell(ctx, args...)
	return err
}

func (a *ADB) SetPower(ctx context.Context, on bool) error {
	code := keycodeSleep
	if on {
		code = keycodeWakeup
	}
	_, err := a.shell(ctx, "input", "keyevent", strconv.Itoa(code))
	return err
}

func (a *ADB) TypeText(ctx context.Context, text string) error {
	err := ValidateText(text)
	if err != nil {
		return err
	}
	_, err = a.shell(ctx, "input", "text", strings.ReplaceAll(text, " ", "%s"))
	return err
}
