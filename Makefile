SHELL := bash

BIN := src/ordna

.PHONY: lint fmt test ci

lint:
	@echo "==> Linting $(BIN)"
	@if command -v shellcheck >/dev/null; then \
		shellcheck $(BIN); \
	else \
		echo "shellcheck not installed; skipping"; \
	fi
	@if command -v shfmt >/dev/null; then \
		shfmt -d -i 4 -ci -sr $(BIN); \
	else \
		echo "shfmt not installed; skipping"; \
	fi

fmt:
	@if command -v shfmt >/dev/null; then \
		shfmt -w -i 4 -ci -sr $(BIN); \
	else \
		echo "shfmt not installed; skipping"; \
	fi

test:
	@if command -v bats >/dev/null; then \
		bats test; \
	else \
		echo "bats not installed; skipping"; \
	fi

ci: lint test
	@echo "==> Building .deb for CI"
	@dpkg-buildpackage -us -uc -b
