#!/bin/bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SOURCES="$ROOT/build/fireplace-sources"
VIDEO_URL="https://commons.wikimedia.org/wiki/Special:FilePath/Fire_%D8%A2%D8%AA%D8%B4_02.ogv"
SOUND_URL="https://commons.wikimedia.org/wiki/Special:FilePath/Bones_breaking_wood_fire_ice_crackling.ogg"
USER_AGENT="AccueilTV/1.0 (https://github.com)"

mkdir -p "$SOURCES"
curl -fsSL -A "$USER_AGENT" -o "$SOURCES/fire.ogv" "$VIDEO_URL"
curl -fsSL -A "$USER_AGENT" -o "$SOURCES/crackle.ogg" "$SOUND_URL"
"$ROOT/scripts/prepare_fireplace.sh" "$SOURCES/fire.ogv" "$SOURCES/crackle.ogg"
