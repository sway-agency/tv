TV ?=
SERIAL ?= $(if $(TV),$(if $(findstring :,$(TV)),$(TV),$(TV):5555))
DOCK ?=
PACKAGES ?=

ANDROID_SDK ?= /opt/homebrew/share/android-commandlinetools
ACCUEIL_JDK ?= /opt/homebrew/opt/openjdk@17
export ANDROID_SDK ACCUEIL_JDK SERIAL

APK := build/accueil.apk
FIREPLACE := assets/fireplace/fire.mp4
SOURCES := $(shell find src -name '*.java') AndroidManifest.xml $(shell find assets -type f)
TV_SH := scripts/tv.sh

.PHONY: help setup connect doctor brand build install home screensaver apps debloat restore dock fireplace quick clean

help:
	@echo "Accueil, a calm Android TV launcher"
	@echo
	@echo "  make setup                    Install build tools (macOS, Homebrew)"
	@echo "  make connect TV=<ip>          Connect to the TV over adb"
	@echo "  make quick TV=<ip>            Build, install, set as home, enable the fireplace screensaver"
	@echo "  make apps TV=<ip>             List apps and recommended packages to disable"
	@echo "  make debloat TV=<ip>          Disable every recommended package found on the TV"
	@echo "  make debloat TV=<ip> PACKAGES=\"a b\"  Disable only these packages"
	@echo "  make restore TV=<ip>          Re-enable everything disabled from this repo"
	@echo "  make dock TV=<ip> DOCK=\"a b c\"  Set dock apps, in order"
	@echo "  make doctor TV=<ip>           Show TV model, Android version and memory"
	@echo "  make brand TV=<ip>            Tell whether this TV brand has a package list"

setup:
	brew install openjdk@17 ffmpeg
	brew install --cask android-commandlinetools android-platform-tools
	yes | JAVA_HOME=$(ACCUEIL_JDK) sdkmanager --sdk_root=$(ANDROID_SDK) "build-tools;36.0.0" "platforms;android-36" >/dev/null

connect:
	@$(TV_SH) connect $(TV)

doctor:
	@$(TV_SH) doctor

brand:
	@$(TV_SH) brand

$(FIREPLACE):
	scripts/fetch_fireplace.sh

fireplace: $(FIREPLACE)

$(APK): $(FIREPLACE) $(SOURCES) build.sh
	./build.sh

build: $(APK)

install: $(APK)
	@$(TV_SH) install

home:
	@$(TV_SH) home

screensaver:
	@$(TV_SH) screensaver

quick: install home screensaver

apps:
	@$(TV_SH) apps

debloat:
	@$(TV_SH) debloat $(PACKAGES)

restore:
	@$(TV_SH) restore

dock:
	@$(TV_SH) dock $(DOCK)

clean:
	rm -rf build
