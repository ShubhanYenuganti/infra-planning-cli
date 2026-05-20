#!/usr/bin/env bash
# doctor-gcp.sh — Pre-flight check for GCP authentication.

set -uo pipefail

pass() { printf "✅ %s\n" "$*"; }
warn() { printf "⚠️  %s\n" "$*"; }
fail() { printf "❌ %s\n" "$*"; }

if [[ -n "${GOOGLE_APPLICATION_CREDENTIALS:-}" ]]; then
  if [[ -f "$GOOGLE_APPLICATION_CREDENTIALS" ]]; then
    pass "GOOGLE_APPLICATION_CREDENTIALS: $GOOGLE_APPLICATION_CREDENTIALS"
  else
    fail "GOOGLE_APPLICATION_CREDENTIALS set but file not found"
  fi
else
  warn "GOOGLE_APPLICATION_CREDENTIALS not set; falling back to gcloud ADC"
fi

ADC_PATH="$HOME/.config/gcloud/application_default_credentials.json"
if [[ -f "$ADC_PATH" ]]; then
  pass "gcloud ADC present: $ADC_PATH"
else
  warn "gcloud ADC not present; run \`gcloud auth application-default login\`"
fi

if command -v gcloud >/dev/null 2>&1; then
  pass "gcloud CLI installed ($(gcloud --version 2>&1 | head -1))"
  active=$(gcloud config get-value account 2>/dev/null || echo "(none)")
  pass "Active gcloud account: $active"
  project=$(gcloud config get-value project 2>/dev/null || echo "(unset)")
  if [[ "$project" == "(unset)" ]]; then
    warn "gcloud project not set; per-CLI commands will require --project"
  else
    pass "gcloud project: $project"
  fi
else
  warn "gcloud CLI not installed; recommended for fallback work"
fi
