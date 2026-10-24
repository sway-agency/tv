package tv

import (
	"context"
	"regexp"
	"sort"
	"strings"
)

var foregroundPattern = regexp.MustCompile(`mResumedActivity: ActivityRecord\{\S+ u\d+ ([a-zA-Z0-9_.]+)/`)

func (a *ADB) Status(ctx context.Context) (Status, error) {
	power, err := a.shell(ctx, "dumpsys", "power")
	if err != nil {
		return Status{}, err
	}
	activities, err := a.shell(ctx, "dumpsys", "activity", "activities")
	if err != nil {
		return Status{}, err
	}
	return Status{
		Awake:             strings.Contains(power, "mWakefulness=Awake"),
		ForegroundPackage: parseForeground(activities),
	}, nil
}

func (a *ADB) Apps(ctx context.Context) ([]App, error) {
	output, err := a.shell(ctx, "cmd", "package", "query-activities",
		"-a", "android.intent.action.MAIN", "-c", "android.intent.category.LEANBACK_LAUNCHER")
	if err != nil {
		return nil, err
	}
	apps := []App{}
	for _, name := range parsePackages(output) {
		apps = append(apps, App{Package: name})
	}
	return apps, nil
}

func parseForeground(activities string) string {
	match := foregroundPattern.FindStringSubmatch(activities)
	if match == nil || ValidatePackage(match[1]) != nil {
		return ""
	}
	return match[1]
}

func parsePackages(output string) []string {
	seen := map[string]bool{}
	for line := range strings.SplitSeq(output, "\n") {
		name, found := strings.CutPrefix(strings.TrimSpace(line), "packageName=")
		if found && ValidatePackage(name) == nil {
			seen[name] = true
		}
	}
	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
