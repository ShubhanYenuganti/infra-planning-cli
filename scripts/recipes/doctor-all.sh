#!/usr/bin/env bash
# doctor-all.sh - run `doctor --json` against every v1 CLI in sequence.
#
# Emits one JSON object per CLI on stdout. Exits 0 if every binary was
# runnable; exits 1 if any binary is missing from PATH.
#
# This is a v1.0 smoke test recipe. A richer cross-CLI audit (requires
# live cloud creds + mocked fixtures) is planned for v1.1+.
set -uo pipefail

CLIS=(
  cloud-run-admin-pp-cli
  cloud-functions-pp-cli
  lambda-pp-cli
  apprunner-pp-cli
  functions-pp-cli
  container-apps-pp-cli
)

missing=0
echo "{"
echo '  "clis": ['
first=true
for cli in "${CLIS[@]}"; do
  if ! command -v "$cli" >/dev/null 2>&1; then
    missing=1
    report='{}'
  else
    report=$("$cli" doctor --json 2>/dev/null || echo "{}")
  fi

  [ "$first" = false ] && echo ","
  first=false
  echo "    {\"cli\": \"$cli\", \"report\":"
  echo "$report"
  echo "    }"
done
echo ""
echo "  ]"
echo "}"
exit "$missing"
