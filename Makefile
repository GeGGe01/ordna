SHELL := bash

BIN := src/ordna

.PHONY: lint fmt test ci

lint:
	@echo "==> Linting $(BIN)"
	@command -v shellcheck >/dev/null \
		&& shellcheck $(BIN) \
		|| echo "shellcheck not installed; skipping"
	@command -v shfmt >/dev/null \
		&& shfmt -d -i 4 -ci -sr $(BIN) \
		|| echo "shfmt not installed; skipping"

fmt:
	@command -v shfmt >/dev/null \
		&& shfmt -w -i 4 -ci -sr $(BIN) \
		|| echo "shfmt not installed; skipping"

test:
	@command -v bats >/dev/null \
		&& bats test \
		|| echo "bats not installed; skipping"

ci: lint test
	@echo "==> Building .deb for CI"
	@dpkg-buildpackage -us -uc -b
