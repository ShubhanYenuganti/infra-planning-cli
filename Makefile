# infra-press verification harness
# See docs/superpowers/specs/2026-05-21-v1-verification-design.md

.PHONY: verify verify-fast verify-lint verify-build verify-golden verify-smoke verify-recipe verify-spec-sync verify-upstream verify-live

verify: verify-fast verify-upstream verify-live

verify-fast: verify-lint verify-build verify-golden verify-smoke verify-recipe verify-spec-sync
	@echo "==> verify-fast complete (L1-L6)"

verify-lint:
	@echo "==> L1 lint"
	@command -v yq >/dev/null 2>&1 || { echo "ERROR: yq is required for verify-lint (brew install yq | snap install yq)"; exit 1; }
	@for path in $$(yq '.clis[].path' catalog.yaml); do \
		echo "  - $$path"; \
		(cd tools/lint-conventions && go run . "../../$$path/") || exit 1; \
	done

verify-build:
	@echo "==> L2 build matrix"
	cd tests/smoke && go test -v -run TestBuildMatrix . -timeout 120s

verify-golden:
	@echo "==> L3 golden"
	cd tests/golden && go test -v ./...

verify-smoke:
	@echo "==> L4 behavioral smoke"
	cd tests/smoke && go test -v -run TestCLISmoke . -timeout 300s

verify-recipe:
	@echo "==> L5 recipe smoke"
	cd tests/smoke && go test -v -run TestRecipeSmoke . -timeout 300s

verify-spec-sync:
	@echo "==> L6 spec sync"
	cd tests/spec_check && go test -v -run TestSpecSync .

verify-upstream:
	@echo "==> L7 upstream advisory (advisory only, never fails PR)"
	cd tests/spec_check && go test -v -run TestUpstreamAdvisory .

verify-live:
	@echo "==> L8 live cloud (skips clouds without OIDC creds)"
	cd tests/live && go test -tags=live -v .
