#!/bin/bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
APK="$ROOT/build/accueil.apk"
LAUNCHER="fr.sway.accueil"
HOME_ACTIVITY="$LAUNCHER/.HomeActivity"
DOCK_RECEIVER="$LAUNCHER/.DockReceiver"
SCREENSAVER="$LAUNCHER/.FireplaceDream"
STATE_DIR="$ROOT/state"
PACKAGE_LISTS="$ROOT/config/packages"

adb_tv() {
    if [ -n "${SERIAL:-}" ]; then
        adb -s "$SERIAL" "$@"
    else
        adb "$@"
    fi
}

shell() {
    adb_tv shell "$@" | tr -d '\r'
}

brand() {
    shell getprop ro.product.manufacturer | tr '[:upper:]' '[:lower:]' | tr -cd 'a-z0-9'
}

bloat_catalog() {
    local catalog="$STATE_DIR/catalog.txt"
    mkdir -p "$STATE_DIR"
    cat "$PACKAGE_LISTS/common.txt" > "$catalog"
    if [ -f "$PACKAGE_LISTS/$(brand).txt" ]; then
        cat "$PACKAGE_LISTS/$(brand).txt" >> "$catalog"
    fi
    echo "$catalog"
}

cmd_brand() {
    local file="$PACKAGE_LISTS/$(brand).txt"
    if [ -f "$file" ]; then
        echo "known $(brand) config/packages/$(brand).txt"
    else
        echo "unknown $(brand) config/packages/$(brand).txt"
    fi
}

device_id() {
    shell getprop ro.serialno
}

state_file() {
    mkdir -p "$STATE_DIR"
    echo "$STATE_DIR/$(device_id)-disabled.txt"
}

catalog_entries() {
    sed -E '/^[[:space:]]*$/d; s/[[:space:]]*\|[[:space:]]*/|/' "$1"
}

enabled_packages() {
    shell pm list packages -e | sed 's/^package://' | sort
}

describe() {
    local package="$1" catalog="$2" pattern description
    while IFS='|' read -r pattern description; do
        if [[ "$package" == $pattern ]]; then
            echo "$description"
            return 0
        fi
    done < <(catalog_entries "$catalog")
    return 1
}

matching_installed() {
    local package
    for package in $(enabled_packages); do
        if describe "$package" "$1" >/dev/null; then
            echo "$package"
        fi
    done
}

disable() {
    local log
    log="$(state_file)"
    for package in "$@"; do
        if [ "$package" = "$LAUNCHER" ]; then
            continue
        fi
        if shell pm disable-user --user 0 "$package" | grep -q "disabled-user"; then
            echo "disabled  $package"
            grep -qx -- "$package" "$log" 2>/dev/null || echo "$package" >> "$log"
        else
            echo "skipped   $package"
        fi
    done
}

cmd_connect() {
    local address="${1:?usage: connect <tv-ip>}"
    [[ "$address" == *:* ]] || address="$address:5555"
    adb connect "$address"
    adb -s "$address" wait-for-device
    if adb devices | grep -q "^$address[[:space:]]*unauthorized"; then
        echo "Accept the debugging prompt on the TV, then run this again." >&2
        exit 1
    fi
    echo "Connected. Use SERIAL=$address for the other commands."
}

cmd_doctor() {
    echo "Device:    $(shell getprop ro.product.manufacturer) $(shell getprop ro.product.model)"
    echo "Android:   $(shell getprop ro.build.version.release) (API $(shell getprop ro.build.version.sdk))"
    echo "Memory:    $(shell grep MemAvailable /proc/meminfo | awk '{print int($2/1024)" MB available"}')"
    echo "Launcher:  $(shell pm list packages "$LAUNCHER" | grep -q . && echo installed || echo missing)"
    echo "Brand:     $(cmd_brand)"
}

cmd_install() {
    adb_tv install -r "$APK"
}

cmd_home() {
    shell cmd package set-home-activity "$HOME_ACTIVITY" >/dev/null
    disable $(matching_installed "$ROOT/config/launchers.txt")
    shell am start -a android.intent.action.MAIN -c android.intent.category.HOME >/dev/null 2>&1
    echo "Accueil is now the home screen."
}

cmd_screensaver() {
    shell settings put secure screensaver_enabled 1
    shell settings put secure screensaver_components "$SCREENSAVER"
    shell settings put secure screensaver_default_component "$SCREENSAVER"
    echo "Fireplace screensaver enabled."
}

cmd_apps() {
    echo "# Launchable apps"
    shell cmd package query-activities -a android.intent.action.MAIN -c android.intent.category.LEANBACK_LAUNCHER \
        | sed -n 's/^ *packageName=//p' | sort -u | grep -vx "$LAUNCHER" || true
    echo
    echo "# Recommended to disable (installed and enabled)"
    local package
    local catalog package
    catalog="$(bloat_catalog)"
    for package in $(matching_installed "$catalog"); do
        echo "$package | $(describe "$package" "$catalog")"
    done
    echo
    echo "# Brand list: $(cmd_brand)"
}

cmd_debloat() {
    if [ "$#" -gt 0 ]; then
        disable "$@"
    else
        disable $(matching_installed "$(bloat_catalog)")
    fi
    shell am kill-all
    shell pm trim-caches 8G
}

cmd_restore() {
    local log
    log="$(state_file)"
    [ -f "$log" ] || { echo "Nothing to restore."; return; }
    while read -r package; do
        shell pm enable "$package" >/dev/null && echo "enabled   $package"
    done < "$log"
    rm "$log"
}

cmd_dock() {
    [ "$#" -gt 0 ] || { echo "usage: dock <package>..." >&2; exit 1; }
    local joined
    joined="$(IFS=,; echo "$*")"
    shell am broadcast -n "$DOCK_RECEIVER" --esa dock "$joined" >/dev/null
    shell am start -a android.intent.action.MAIN -c android.intent.category.HOME >/dev/null 2>&1
    echo "Dock: $joined"
}

command="${1:-}"
shift || true
case "$command" in
    connect|doctor|brand|install|home|screensaver|apps|debloat|restore|dock) "cmd_$command" "$@" ;;
    *) echo "usage: $0 {connect|doctor|brand|install|home|screensaver|apps|debloat|restore|dock}" >&2; exit 1 ;;
esac
