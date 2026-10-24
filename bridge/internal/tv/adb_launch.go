package tv

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

const FireplaceComponent = "fr.sway.accueil/.FireplaceActivity"

var ErrAppNotInstalled = errors.New("app is not installed or cannot be launched on the TV")

func (a *ADB) LaunchApp(ctx context.Context, packageName string) error {
	err := ValidatePackage(packageName)
	if err != nil {
		return err
	}
	installed, err := a.isLaunchable(ctx, packageName)
	if err != nil {
		return err
	}
	if !installed {
		return ErrAppNotInstalled
	}
	component, err := a.resolveLauncher(ctx, packageName)
	if err != nil {
		return err
	}
	_, err = a.shell(ctx, "am", "start", "-n", component)
	return err
}

func (a *ADB) StartFireplace(ctx context.Context) error {
	_, err := a.shell(ctx, "am", "start", "-n", FireplaceComponent)
	return err
}

func (a *ADB) isLaunchable(ctx context.Context, packageName string) (bool, error) {
	apps, err := a.Apps(ctx)
	if err != nil {
		return false, err
	}
	for _, app := range apps {
		if app.Package == packageName {
			return true, nil
		}
	}
	return false, nil
}

func (a *ADB) resolveLauncher(ctx context.Context, packageName string) (string, error) {
	output, err := a.shell(ctx, "cmd", "package", "resolve-activity", "--brief",
		"-a", "android.intent.action.MAIN", "-c", "android.intent.category.LEANBACK_LAUNCHER", packageName)
	if err != nil {
		return "", err
	}
	lines := strings.Split(strings.TrimSpace(output), "\n")
	component := strings.TrimSpace(lines[len(lines)-1])
	if !strings.HasPrefix(component, packageName+"/") || validateComponent(component) != nil {
		return "", fmt.Errorf("%w: no TV launcher activity", ErrAppNotInstalled)
	}
	return component, nil
}
