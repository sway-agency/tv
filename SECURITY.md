# Security policy

## Reporting a vulnerability

Please report vulnerabilities privately through GitHub: **Security**, **Report a vulnerability** on [sway-agency/tv](https://github.com/sway-agency/tv/security/advisories/new). Do not open a public issue.

We acknowledge reports within 72 hours and aim to release a fix within 14 days for critical issues.

## Supported versions

Only the latest commit on `main` is supported. Rebuild and redeploy the launcher and the bridge to receive fixes.

## Scope

- The Android launcher (`src/`), which requests no Android permission.
- The AI bridge (`bridge/`), see [`docs/bridge.md`](docs/bridge.md) for its threat model.
- The setup scripts (`scripts/`, `Makefile`).

Out of scope: vulnerabilities of the TV firmware, of `adb` itself, or of third-party apps installed on the TV.

## Design commitments

- The launcher never requests network, camera, microphone or location permissions.
- The bridge has no third-party dependency, never listens on a public address, and exposes no tool to run commands, install apps, open URLs, display external content or access the camera or microphone.
- Changes to the TV are reversible and logged locally.
