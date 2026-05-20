#!/usr/bin/env bash
# doctor-aws.sh — Pre-flight check for AWS authentication.
# Does NOT call any cloud APIs. Reports what the SDK chain would find.

set -uo pipefail

pass() { printf "✅ %s\n" "$*"; }
warn() { printf "⚠️  %s\n" "$*"; }
fail() { printf "❌ %s\n" "$*"; }

if [[ -n "${AWS_PROFILE:-}" ]]; then
  pass "AWS_PROFILE set: $AWS_PROFILE"
else
  warn "AWS_PROFILE not set; default profile or instance role will be used"
fi

if [[ -f "$HOME/.aws/credentials" ]]; then
  pass "~/.aws/credentials present"
else
  warn "~/.aws/credentials missing; SDK may fall back to instance/container role"
fi

if [[ -n "${AWS_REGION:-}" || -n "${AWS_DEFAULT_REGION:-}" ]]; then
  pass "AWS_REGION set: ${AWS_REGION:-$AWS_DEFAULT_REGION}"
else
  warn "AWS_REGION not set; per-CLI commands will require --region"
fi

if command -v aws >/dev/null 2>&1; then
  pass "aws CLI installed ($(aws --version 2>&1))"
else
  warn "aws CLI not installed; recommended for IAM/VPC fallbacks until v1.1+"
fi

if command -v aws >/dev/null 2>&1; then
  if identity=$(aws sts get-caller-identity --output text 2>&1); then
    pass "STS GetCallerIdentity: $identity"
  else
    fail "STS GetCallerIdentity failed; auth not working"
    echo "    $identity"
  fi
fi
