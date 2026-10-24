package tv

import "context"

type App struct {
	Package string `json:"package"`
}

type Status struct {
	Awake             bool   `json:"awake"`
	ForegroundPackage string `json:"foreground_package"`
}

type Controller interface {
	Status(ctx context.Context) (Status, error)
	Apps(ctx context.Context) ([]App, error)
	PressKey(ctx context.Context, key Key, repeat int) error
	LaunchApp(ctx context.Context, packageName string) error
	SetPower(ctx context.Context, on bool) error
	StartFireplace(ctx context.Context) error
	TypeText(ctx context.Context, text string) error
	Screenshot(ctx context.Context) ([]byte, error)
}
