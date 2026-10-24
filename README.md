# Accueil, a calm Android TV launcher by Sway

A minimal home screen for Android TV: a dock with your apps, slowly changing winter wallpapers, spring animations, light snowfall and a crackling fireplace (as an app and as the screensaver). It also cleans up the TV by disabling the preloaded apps you do not need, reversibly.

Plain Java, no Gradle, no dependencies. One `make` command per step.

## Install with an AI agent (easiest)

1. On the TV, enable developer options (Settings, Device preferences, About, press Build 7 times), then enable USB debugging and network debugging.
2. Open Claude Code (or any coding agent) and say:

   > Install https://github.com/sway-agency/tv on my TV

The agent follows `CLAUDE.md`: it connects over adb, installs the launcher, replaces the stock home screen, then asks you which apps to disable and which to pin in the dock.

## Install by hand

```sh
make setup                          # macOS: JDK 17, Android SDK, adb, ffmpeg
make connect TV=192.168.1.20        # accept the prompt on the TV
make quick TV=192.168.1.20          # build, install, set as home, fireplace screensaver
make apps TV=192.168.1.20           # launchable apps and recommended packages to disable
make debloat TV=192.168.1.20        # disable the recommended packages
make dock TV=192.168.1.20 DOCK="com.netflix.ninja com.disney.disneyplus"
```

`make restore TV=<ip>` re-enables everything disabled from this repository. `make help` lists every target.

Outside macOS, set `ANDROID_SDK` (with `build-tools;36.0.0` and `platforms;android-36`) and `ACCUEIL_JDK` (JDK 17).

## Supported TVs

Any Android TV or Google TV on Android 9 or newer. Package lists live in `config/packages/`: `common.txt` for every TV, plus one file per brand. Your brand is missing? The agent can write the list with you and open a pull request, see `CONTRIBUTING.md`.

## Customising

- **Dock**: `make dock`, no rebuild needed
- **Wallpapers**: any JPEG in `assets/wallpapers/`, interval in `WallpaperRotator.java`
- **Font**: drop a variable font at `assets/fonts/satoshi.ttf`; without it the system sans-serif is used
- **Fireplace**: `make fireplace` downloads a CC0 fire loop; replace `assets/fireplace/fire.mp4` with your own if you like

Keep your `accueil.keystore`: Android refuses to update the app if it is signed with another key.

## Contributing

Read `NORME.md` and `CONTRIBUTING.md`. The norminette (`python3 scripts/norminette.py`) and a build run on every pull request.

## Licence

Code under the [MIT licence](LICENSE). The Sway logo is a Sway trademark: it may only be redistributed unmodified, where the launcher displays it. Wallpapers, sounds and icons keep their own licences, listed in `CREDITS.txt`.
