# Contributing

Thanks for helping. Humans and AI agents follow the same rules.

## Before you start

- Read `NORME.md`. Code that fails the norminette is not merged.
- AI agents also read `CLAUDE.md` (or `AGENTS.md`, the same file).
- The Sway logos and where they are displayed are not open to changes (N-SWAY-01, N-SWAY-02).

## Adding a TV brand

The most useful contribution is a package list for a brand we do not know yet.

1. Install Accueil on the TV (see the quick install in `CLAUDE.md`).
2. Disable the manufacturer preloads you do not need with `make debloat PACKAGES="..."`.
3. Test the TV: home button, every input, the remote, Wi-Fi, Bluetooth, settings, launching apps, a reboot.
4. Add `config/packages/<brand>.txt`, where `<brand>` is the output of `make brand`. One `package.name | Description` per line, only packages you disabled and tested.
5. Run `python3 scripts/norminette.py`.
6. Commit with `docs(packages): add <Brand> package list` and open a pull request from a `brand/<brand>` branch with the template filled in.

Packages found on every Android TV belong in `common.txt`; open a separate pull request for those.

## Code changes

1. Branch from `main`: `feature/<topic>` or `fix/<topic>`.
2. Keep the change small and focused.
3. Run `python3 scripts/norminette.py`, `./build.sh`, and `make bridge-test` if you touched `bridge/`.
4. If you have a TV, install it with `make quick TV=<ip>` and attach a screenshot to the pull request.
5. Commit following the rules below.

## Commits

Commits follow [Conventional Commits](https://www.conventionalcommits.org) on **one line only** (`git log --oneline` must tell the whole story):

```text
type(scope): imperative description in English
```

- Types: `feat`, `fix`, `refactor`, `perf`, `docs`, `style`, `test`, `build`, `ci`, `chore`.
- No body, no trailing full stop.
- No `Co-Authored-By` trailer and no mention of a tool or an AI, in commits or pull requests.

```text
feat(dock): set dock apps over adb without rebuilding
docs(packages): add Philips package list
fix(fireplace): release players when the surface is destroyed
```

## Pull requests

`main` is protected: every change goes through a pull request, CI must pass, and force-pushes are refused. Pull requests are squash-merged, so the PR title becomes the commit and must follow the commit rules above (CI checks it).

Fill in `.github/pull_request_template.md`. CI runs on every pull request:

- **norminette**: the Sway norm (`NORME.md`);
- **java**: compilation with every warning as an error, then the APK build;
- **go**: `gofmt`, `go vet`, tests with the race detector, build, `govulncheck`;
- **CodeQL**: security analysis of the Java, Go and workflow code.

Dependabot opens update pull requests for GitHub Actions and the Go toolchain every week.
