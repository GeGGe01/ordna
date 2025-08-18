# Repository Guidelines

## Project Structure & Module Organization
- `src/ordna`: Primary Bash CLI. Keep all runtime logic here; prefer small helpers over subshell pipelines.
- `debian/`: Packaging metadata for building a `.deb` (debhelper 13, native package).
- `.github/workflows/publish.yml`: CI to build and publish an APT repo from tags (`v*`).
- Tests: none yet. If added, place under `test/` (see Testing Guidelines).

## Build, Test, and Development Commands
- Local run (no install): `bash src/ordna --help` or `src/ordna ...`
- Build Debian package: `dpkg-buildpackage -us -uc -b` (from repo root; outputs `../ordna_*_all.deb`).
- Lint shell: `shellcheck src/ordna`
- Format shell: `shfmt -w -i 4 src/ordna`
- Quick dry-run example:
  - `src/ordna ~/Pictures ~/out -c -x --dry-run --from 2024-01-01 --to 2024-12-31`
 - Make targets: `make lint` (shellcheck+shfmt check), `make fmt` (format), `make test` (run bats), `make ci` (lint, test, build .deb)

## Coding Style & Naming Conventions
- Language: Bash (bash >= 4.2). Use `set -euo pipefail` and quote variables.
- Indentation: 4 spaces; wrap long pipelines into readable steps.
- Functions: `lower_snake_case` (e.g., `get_epoch`, `target_dir_for`).
- Globals/flags: UPPER_SNAKE_CASE (e.g., `DRYRUN`, `STRICT_EXT`).
- Prefer arrays and built-ins over external processes; batch `exiftool` calls when possible.

## Testing Guidelines
- Framework: optional; recommended `bats-core` for CLI tests under `test/*.bats`.
- Naming: mirror feature names (e.g., `test/sort_by_ext.bats`).
- Run (if using bats): `bats test`
- Until a suite exists, include reproducible `--dry-run` transcripts in PRs and cover edge cases (empty dirs, duplicate hashes, unknown extensions, date ranges).

## Commit & Pull Request Guidelines
- Commits: imperative mood, concise scope (e.g., "add ext-based sorting"), group related changes.
- Link issues in commit bodies (`Fixes #123`) when relevant.
- PRs must include: summary, rationale, sample commands + output (prefer `--dry-run`), and risk notes.
- CI: tags `v*` trigger publish workflow; avoid tagging until review is complete.

## Security & Configuration Tips
- Dependencies: `exiftool`, `file`, `coreutils`, `findutils`, `grep`, `gawk`, `sed`. Ensure GNU `getopt` and `date -d` are available.
- Do not run as root. Review outputs before moving files; use `--dry-run` to validate.
- Publishing: requires GPG secrets in GitHub Actions; never commit keys.
