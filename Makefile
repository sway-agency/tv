TV ?=
SERIAL ?= $(if $(TV),$(if $(findstring :,$(TV)),$(TV),$(TV):5555))
DOCK ?=
PACKAGES ?=
LISTEN ?=
HOSTS ?=
NAME ?=
SCOPES ?= observe,control
BRIDGE := build/accueil-bridge
BRIDGE_VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

ANDROID_SDK ?= /opt/homebrew/share/android-commandlinetools
ACCUEIL_JDK ?= /opt/homebrew/opt/openjdk@17
export ANDROID_SDK ACCUEIL_JDK SERIAL

APK := build/accueil.apk
FIREPLACE := assets/fireplace/fire.mp4
SOURCES := $(shell find src -name '*.java') AndroidManifest.xml $(shell find assets -type f)
TV_SH := scripts/tv.sh

.PHONY: help setup connect doctor brand build install home screensaver apps debloat restore dock fireplace quick clean \
	privacy privacy-lock bridge bridge-test bridge-audit bridge-init bridge-token bridge-serve bridge-macos

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
	@echo "  make privacy TV=<ip>          List apps allowed to use the microphone, a camera or the location"
	@echo "  make privacy-lock TV=<ip> PACKAGES=\"a b\"  Revoke those permissions from these apps"
	@echo
	@echo "  make bridge                   Build the AI bridge (MCP server) into build/accueil-bridge"
	@echo "  make bridge-test              Vet, test and scan the bridge for known vulnerabilities"
	@echo "  make bridge-init TV=<ip> LISTEN=<ip:port> HOSTS=\"name ip\"  Create keys and config"
	@echo "  make bridge-token NAME=<agent> [SCOPES=observe,control]  Create an access token"
	@echo "  make bridge-serve             Run the bridge in the foreground"
	@echo "  make bridge-macos             Install the bridge as a macOS login agent"

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

privacy:
	@$(TV_SH) privacy

privacy-lock:
	@$(TV_SH) privacy-lock $(PACKAGES)

$(BRIDGE): $(shell find bridge -name '*.go') bridge/go.mod
	cd bridge && CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=$(BRIDGE_VERSION)" -o ../$(BRIDGE) ./cmd/accueil-bridge

bridge: $(BRIDGE)

bridge-test:
	cd bridge && test -z "$$(gofmt -l .)" && go vet ./... && go test -race -count=1 ./...
	cd bridge && go run golang.org/x/vuln/cmd/govulncheck@latest ./...

bridge-audit: bridge-test

bridge-init: $(BRIDGE)
	$(BRIDGE) init --listen $(LISTEN) --tv $(SERIAL) $(foreach host,$(HOSTS),--host $(host))

bridge-token: $(BRIDGE)
	$(BRIDGE) token add --name $(NAME) --scopes $(SCOPES)

bridge-serve: $(BRIDGE)
	$(BRIDGE) serve

bridge-macos: $(BRIDGE)
	sudo install -m 0755 $(BRIDGE) /usr/local/bin/accueil-bridge
	sed "s|__HOME__|$$HOME|" bridge/deploy/fr.sway.accueil-bridge.plist > $$HOME/Library/LaunchAgents/fr.sway.accueil-bridge.plist
	launchctl bootout gui/$$(id -u)/fr.sway.accueil-bridge 2>/dev/null || true
	launchctl bootstrap gui/$$(id -u) $$HOME/Library/LaunchAgents/fr.sway.accueil-bridge.plist

clean:
	rm -rf build
