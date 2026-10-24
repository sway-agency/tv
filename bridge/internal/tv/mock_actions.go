package tv

import (
	"context"
	"fmt"
	"slices"
)

func (m *Mock) LaunchApp(ctx context.Context, packageName string) error {
	if !slices.Contains(m.Installed, packageName) {
		return ErrAppNotInstalled
	}
	return m.record("launch " + packageName)
}

func (m *Mock) SetPower(ctx context.Context, on bool) error {
	return m.record(fmt.Sprintf("power %t", on))
}

func (m *Mock) StartFireplace(ctx context.Context) error {
	return m.record("fireplace")
}

func (m *Mock) TypeText(ctx context.Context, text string) error {
	err := ValidateText(text)
	if err != nil {
		return err
	}
	return m.record("text")
}

func (m *Mock) Screenshot(ctx context.Context) ([]byte, error) {
	return []byte{0xff, 0xd8, 0xff, 0xd9}, m.record("screenshot")
}
