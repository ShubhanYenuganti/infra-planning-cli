# infra-press verification harness
# See docs/superpowers/specs/2026-05-21-v1-verification-design.md

.PHONY: verify verify-fast verify-spec-sync verify-lint verify-build verify-golden verify-smoke verify-recipe verify-upstream verify-live

# Default target — full verify (will grow as later phases land).
verify: verify-spec-sync

verify-fast: verify-spec-sync

verify-spec-sync:
	@echo "==> L6 spec sync"
	cd tests/spec_check && go test -v -run TestSpecSync .
