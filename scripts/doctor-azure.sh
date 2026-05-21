#!/usr/bin/env bash
# doctor-azure.sh — Pre-flight check for Azure authentication.

set -uo pipefail

pass() { printf "✅ %s\n" "$*"; }
warn() { printf "⚠️  %s\n" "$*"; }
fail() { printf "❌ %s\n" "$*"; }

if [[ -n "${AZURE_TENANT_ID:-}" && -n "${AZURE_CLIENT_ID:-}" && -n "${AZURE_CLIENT_SECRET:-}" ]]; then
  pass "Service principal env vars set (TENANT/CLIENT_ID/CLIENT_SECRET)"
else
  warn "Service principal env vars not set; will fall back to az CLI or managed identity"
fi

if command -v az >/dev/null 2>&1; then
  pass "az CLI installed ($(az --version 2>&1 | head -1))"
  if account=$(az account show --output json 2>&1); then
    if command -v python3 >/dev/null 2>&1; then
      name=$(echo "$account" | python3 -c "import json,sys; print(json.load(sys.stdin).get('user',{}).get('name',''))" 2>/dev/null || true)
      sub=$(echo "$account" | python3 -c "import json,sys; print(json.load(sys.stdin).get('name',''))" 2>/dev/null || true)
      if [[ -n "$name" ]]; then
        pass "az logged in as: $name"
      else
        warn "az account user could not be parsed"
      fi
      if [[ -n "$sub" ]]; then
        pass "Active subscription: $sub"
      else
        warn "az active subscription could not be parsed"
      fi
    else
      warn "python3 not found; cannot parse az account details"
    fi
  else
    warn "az not logged in; run \`az login\`"
  fi
else
  warn "az CLI not installed; recommended fallback"
fi

if [[ -n "${AZURE_SUBSCRIPTION_ID:-}" ]]; then
  pass "AZURE_SUBSCRIPTION_ID set"
else
  warn "AZURE_SUBSCRIPTION_ID not set; CLIs will need --subscription flag"
fi

echo -n "swagger2openapi (Azure spec converter): "
if command -v swagger2openapi >/dev/null 2>&1; then
  echo "OK ($(swagger2openapi --version 2>/dev/null | head -n1))"
else
  echo "MISSING (only needed to regenerate Azure CLIs from spec; run scripts/install.sh)"
fi
