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
    name=$(echo "$account" | python3 -c "import json,sys; print(json.load(sys.stdin).get('user',{}).get('name','?'))" 2>/dev/null || echo '?')
    pass "az logged in as: $name"
    sub=$(echo "$account" | python3 -c "import json,sys; print(json.load(sys.stdin).get('name','?'))" 2>/dev/null || echo '?')
    pass "Active subscription: $sub"
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
