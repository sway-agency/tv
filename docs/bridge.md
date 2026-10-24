# AI bridge: control the TV from your AI agent

`accueil-bridge` lets an AI agent (Claude Code, Claude Desktop, Hermes Agent, OpenJarvis, a voice assistant, any [MCP](https://modelcontextprotocol.io) client or plain HTTPS client) control the TV: remote keys, launching apps, turning the screen on or off, the fireplace.

It is built to be safe to leave running at home, with security first and features second.

## How it works

```text
 AI agent ── TLS 1.3 + token ──> accueil-bridge ── adb (LAN) ──> Android TV
 (laptop, phone,                 (Raspberry Pi,                   (launcher has no
  home server)                    Mac, NAS)                        network permission)
```

- The bridge runs on an always-on machine at home, never on the TV. The Accueil launcher requests **no Android permission at all**, not even internet access, so it cannot leak anything.
- The bridge speaks MCP over Streamable HTTP (protocol `2026-07-28`, and the `2025-03-26` to `2025-11-25` versions for older agents), MCP over stdio, and a small REST endpoint for voice shortcuts.
- It is written in Go with the standard library only: **zero third-party dependencies**, so there is no supply chain to compromise. Updating Go is enough to receive every security fix.

## What an agent can and cannot do

| Tool | Scope | What it does |
| --- | --- | --- |
| `tv_status` | `observe` | Screen on or off, foreground app |
| `tv_list_apps` | `observe` | Installed TV apps |
| `tv_press_key` | `control` | Remote keys from a fixed list (navigation, playback, volume) |
| `tv_launch_app` | `control` | Opens an installed app, checked against the installed list |
| `tv_power` | `control` | Wakes the screen or puts it in standby |
| `tv_fireplace` | `control` | Starts the fireplace |
| `tv_type_text` | `text` | Types letters, digits and spaces into a search field |
| `tv_screenshot` | `screen` | Small JPEG of the screen |

Each token only sees the tools of its scopes. `text` and `screen` are **off by default** and must be granted explicitly.

By design, there is no tool and no code path to:

- use the TV camera or microphone, or any other app permission;
- run a shell command, install or remove an app, open a URL, or send arbitrary Android intents;
- display images, videos or text pushed from outside.

Every argument is validated against a strict allowlist before it reaches `adb`: key names come from a fixed table, package names must match an installed app, text only accepts letters, digits and spaces. Processes can only be started from one file (`internal/tv/adb.go`), which the norminette enforces.

## Install

On the machine that will run the bridge (macOS or Linux), with the TV already set up with `make quick`:

```sh
make bridge
make bridge-init TV=192.168.1.20 LISTEN=192.168.1.10:8443 HOSTS="tv-bridge.lan 192.168.1.10"
make bridge-token NAME=claude
make bridge-serve
```

- `LISTEN` is the LAN (or WireGuard) address of the bridge machine. The bridge refuses `0.0.0.0`, `::` and any public address.
- `HOSTS` lists the names and addresses clients will use. Requests with another `Host` header are rejected, which blocks DNS rebinding.
- `bridge-init` creates a private certificate authority, a TLS certificate and `~/.config/accueil-bridge/config.json`, all readable by your user only. The bridge refuses to start if they are not.
- `bridge-token` prints a token **once**. Only its SHA-256 hash is stored. Add `SCOPES=observe,control,text` to allow typing text.

To run it permanently:

- **Linux (Raspberry Pi, NAS, home server)**: copy the binary to `/usr/local/bin`, create an `accueil-bridge` system user, put the configuration in `/var/lib/accueil-bridge/config`, and install [`bridge/deploy/accueil-bridge.service`](../bridge/deploy/accueil-bridge.service). The unit is sandboxed: no privileges, read-only system, no access to home folders, private network addresses only.
- **macOS**: `make bridge-macos` installs a login agent. macOS may ask once to allow local network access.

## Connect an agent

Every client needs the bridge address, the token, and to trust `ca.pem` from the configuration folder.

**Claude Code**

```sh
claude mcp add --transport http tv https://tv-bridge.lan:8443/mcp --header "Authorization: Bearer acb_..."
export NODE_EXTRA_CA_CERTS=~/.config/accueil-bridge/ca.pem
```

**Agent on the same machine as the bridge** (Claude Desktop, Hermes Agent, OpenJarvis or any MCP client that launches local servers): use stdio, no network involved.

```json
{ "mcpServers": { "tv": { "command": "/usr/local/bin/accueil-bridge", "args": ["stdio", "--scopes", "observe,control"] } } }
```

**Other MCP clients**: Streamable HTTP endpoint `https://<host>:8443/mcp` with the header `Authorization: Bearer <token>`.

**Voice assistants and shortcuts** (Siri Shortcuts, Home Assistant, scripts): `POST https://<host>:8443/v1/tools/<tool>` with the same header and the tool arguments as a JSON body.

```sh
curl --cacert ca.pem -H "Authorization: Bearer acb_..." -H "Content-Type: application/json" \
  -d '{"package":"com.netflix.ninja"}' https://tv-bridge.lan:8443/v1/tools/tv_launch_app
```

## Remote access, safely

**Never forward the bridge port (8443) or the TV adb port (5555) on your internet box.** An exposed adb port gives full control of the TV to anyone.

To use the bridge away from home, use a VPN that ends at home:

1. Turn on the **WireGuard server of your box** (Freebox: Freebox OS, Paramètres de la Freebox, Serveur VPN, WireGuard; most recent boxes and routers have one). Otherwise run WireGuard on the bridge machine.
2. Install the WireGuard profile on your phone or laptop. WireGuard does not answer unauthenticated packets, so the port it opens stays invisible to scanners.
3. Reach the bridge at its LAN address through the VPN. Add the VPN subnet to the firewall below.

Cloud agents that cannot join your VPN (for example a chat app running in a vendor's cloud) are deliberately not supported: that would require exposing the bridge to the internet.

## Firewall

On Linux, [`bridge/deploy/nftables.conf`](../bridge/deploy/nftables.conf) only lets the LAN and the WireGuard subnet reach port 8443. On macOS, keep the application firewall on and the bridge bound to the LAN address.

On the TV side:

- In **Developer options**, revoke USB debugging authorizations, then only accept the bridge machine: the TV then refuses adb from any other device.
- If your box supports it, put the TV on a separate network (guest or IoT) that can reach the internet but not your computers.

## Privacy on the TV

```sh
make privacy TV=192.168.1.20
make privacy-lock TV=192.168.1.20 PACKAGES="com.example.app"
```

`make privacy` lists the camera and microphone hardware, and every app allowed to use the microphone, a camera or the location. `make privacy-lock` revokes those permissions from the apps you name (reversible with `make restore`). Keep the voice assistant only if you use the remote microphone.

## Security checklist

- TLS 1.3 only, private CA, optional client certificates (`accueil-bridge cert client --name phone`, then `"require_client_cert": true`).
- 256-bit tokens, stored hashed, compared in constant time, one per agent, revocable (`accueil-bridge token revoke --name ...`).
- Rate limit per token, lockout of an address after 10 failed attempts.
- `Host` and `Origin` checks, `application/json` only, 64 KB request limit, strict timeouts, security headers.
- Local JSON audit log: client, address, method, tool, status. Never tokens, never typed text.
- No telemetry, no outbound connection except to the TV.
- CI runs the tests, `go vet`, `govulncheck` and CodeQL on every change. Dependabot keeps the toolchain and actions up to date.

Report vulnerabilities privately, see [`SECURITY.md`](../SECURITY.md).
