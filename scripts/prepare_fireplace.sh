#!/bin/bash
set -euo pipefail

SOURCE_VIDEO="$1"
SOURCE_SOUND="$2"
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
OUT="$ROOT/assets/fireplace"
CROSSFADE_SECONDS=2

if command -v clamscan >/dev/null; then
    clamscan --no-summary "$SOURCE_VIDEO" "$SOURCE_SOUND"
fi
mkdir -p "$OUT"

duration=$(ffprobe -v error -show_entries format=duration -of csv=p=0 "$SOURCE_VIDEO")
loop_offset=$(echo "$duration - 2 * $CROSSFADE_SECONDS" | bc -l)

ffmpeg -v error -y -i "$SOURCE_VIDEO" -filter_complex "
    [0:v]scale=1920:1080:force_original_aspect_ratio=increase,crop=1920:1080,fps=30,split[body][head];
    [head]trim=0:${CROSSFADE_SECONDS},setpts=PTS-STARTPTS[intro];
    [body]trim=${CROSSFADE_SECONDS},setpts=PTS-STARTPTS[main];
    [main][intro]xfade=transition=fade:duration=${CROSSFADE_SECONDS}:offset=${loop_offset}[loop]" \
    -map "[loop]" -an -c:v libx264 -profile:v high -level 4.1 -pix_fmt yuv420p -crf 22 -movflags +faststart \
    "$OUT/fire.mp4"

ffmpeg -v error -y -i "$SOURCE_SOUND" -af "afade=t=in:d=0.5,areverse,afade=t=in:d=0.5,areverse,loudnorm=I=-24" \
    -ac 2 -ar 48000 -c:a aac -b:a 160k "$OUT/crackle.m4a"

ls -la "$OUT"
