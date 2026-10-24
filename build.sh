#!/bin/bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")" && pwd)"
SDK="${ANDROID_SDK:-/opt/homebrew/share/android-commandlinetools}"
TOOLS="$SDK/build-tools/36.0.0"
PLATFORM="$SDK/platforms/android-36/android.jar"
JDK="${ACCUEIL_JDK:-/opt/homebrew/opt/openjdk@17}"
export JAVA_HOME="$JDK" PATH="$JDK/bin:$PATH"
OUT="$ROOT/build"
KEYSTORE="${ACCUEIL_KEYSTORE:-$ROOT/accueil.keystore}"
KEYSTORE_PASS="${ACCUEIL_KEYSTORE_PASS:-accueil}"

rm -rf "$OUT"
mkdir -p "$OUT/classes" "$OUT/dex"

"$JDK/bin/javac" --release 8 -Xlint:-options -classpath "$PLATFORM" \
    -d "$OUT/classes" $(find "$ROOT/src" -name '*.java')
"$TOOLS/d8" --release --min-api 28 --lib "$PLATFORM" --output "$OUT/dex" \
    $(find "$OUT/classes" -name '*.class')
"$TOOLS/aapt2" link -I "$PLATFORM" --manifest "$ROOT/AndroidManifest.xml" \
    -A "$ROOT/assets" -0 jpg -0 mp4 -0 m4a -o "$OUT/unsigned.apk"
(cd "$OUT/dex" && zip -q -j "$OUT/unsigned.apk" classes.dex)
"$TOOLS/zipalign" -f -p 4 "$OUT/unsigned.apk" "$OUT/aligned.apk"

if [ ! -f "$KEYSTORE" ]; then
    "$JDK/bin/keytool" -genkeypair -keystore "$KEYSTORE" -storepass "$KEYSTORE_PASS" -keypass "$KEYSTORE_PASS" \
        -alias accueil -keyalg RSA -keysize 2048 -validity 10000 -dname "CN=Accueil TV" >/dev/null
fi
"$TOOLS/apksigner" sign --ks "$KEYSTORE" --ks-pass "pass:$KEYSTORE_PASS" \
    --out "$OUT/accueil.apk" "$OUT/aligned.apk"
echo "$OUT/accueil.apk"
