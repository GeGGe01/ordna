# Releasing Ordna

This guide describes how to cut a release, publish the APT repo, and how users can install it.

## Prerequisites
- GitHub Pages enabled (Source: GitHub Actions).
- Secrets set: `GPG_PRIVATE_KEY` (ASCII‑armored private key) and `GPG_PASSPHRASE`.
- Local tools (optional): `shellcheck`, `shfmt`, `bats`, `devscripts`.

## Release Steps
1. Verify quality locally:
   - `make lint && make test`
   - `make ci` (builds `../ordna_*_all.deb`)
2. Pick a tag (SemVer): `vX.Y.Z` (e.g., `v0.1.0`).
3. Tag and push:
   - `git tag v0.1.0`
   - `git push origin v0.1.0`
4. GitHub Actions (publish workflow) will:
   - Lint and test, build `.deb`.
   - Generate and sign APT metadata.
   - Publish the APT repo to GitHub Pages.
   - Upload the `.deb` as an artifact.
5. Confirm deployment: check the “deploy” job for the Pages URL.

## Test/Dry Run
- Use a pre‑release tag like `v0.0.0-test1`.
- Remove when done: `git push --delete origin v0.0.0-test1 && git tag -d v0.0.0-test1`.

## APT Install Instructions (for users)
Replace `<PAGES_URL>` with your deployed Pages URL (shown in the workflow), which serves the `publish/` directory.

```sh
# 1) Add key
curl -fsSL <PAGES_URL>/ordna-archive-keyring.asc | \
  sudo tee /usr/share/keyrings/ordna-archive-keyring.asc >/dev/null

# 2) Add source (stable channel)
echo "deb [signed-by=/usr/share/keyrings/ordna-archive-keyring.asc arch=all] \
<PAGES_URL> stable main" | sudo tee /etc/apt/sources.list.d/ordna.list

# 3) Install
sudo apt-get update && sudo apt-get install ordna
```

## Troubleshooting
- Secrets missing: publish job fails early—set both secrets.
- Lint/test failures: fix locally and retag.
- Pages cache: changes take ~1–2 minutes to propagate.
- GPG issues: ensure the key has a passphrase and was imported successfully in the workflow logs.

