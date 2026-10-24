package tv

import (
	"context"
	"fmt"
)

type Mock struct {
	Calls     []string
	Installed []string
	Fail      error
}

func NewMock(installed ...string) *Mock {
	return &Mock{Installed: installed}
}

func (m *Mock) Status(ctx context.Context) (Status, error) {
	m.Calls = append(m.Calls, "status")
	return Status{Awake: true, ForegroundPackage: "fr.sway.accueil"}, m.Fail
}

func (m *Mock) Apps(ctx context.Context) ([]App, error) {
	m.Calls = append(m.Calls, "apps")
	apps := []App{}
	for _, name := range m.Installed {
		apps = append(apps, App{Package: name})
	}
	return apps, m.Fail
}

func (m *Mock) PressKey(ctx context.Context, key Key, repeat int) error {
	m.Calls = append(m.Calls, fmt.Sprintf("key %s x%d", key.Name, repeat))
	return m.Fail
}

func (m *Mock) record(call string) error {
	m.Calls = append(m.Calls, call)
	return m.Fail
}
