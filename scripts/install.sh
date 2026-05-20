#!/usr/bin/env bash
# install.sh — bulk install all `status: stable` CLIs from catalog.yaml.
#
# Idempotent: re-run to upgrade to whatever catalog.yaml on main pins.

set -euo pipefail

INSTALL_DIR="${INFRA_PRESS_HOME:-$HOME/.infra-press}/bin"
CATALOG_URL="${INFRA_PRESS_CATALOG:-https://raw.githubusercontent.com/ShubhanYenuganti/infra-press/main/catalog.yaml}"

# Detect platform
case "$(uname -s)" in
  Darwin) OS=darwin ;;
  Linux)  OS=linux ;;
  *) echo "unsupported OS: $(uname -s)" >&2; exit 1 ;;
esac
case "$(uname -m)" in
  x86_64)  ARCH=amd64 ;;
  arm64|aarch64) ARCH=arm64 ;;
  *) echo "unsupported arch: $(uname -m)" >&2; exit 1 ;;
esac

mkdir -p "$INSTALL_DIR"

echo "Fetching catalog from $CATALOG_URL..."
catalog=$(curl -sSL "$CATALOG_URL")

if ! command -v yq >/dev/null 2>&1; then
  echo "yq is required. Install with: brew install yq  (or your package manager)" >&2
  exit 1
fi

echo "$catalog" | yq '.clis[] | select(.status == "stable") | .name + " " + .press_version' \
| while read -r cli version; do
  if [[ -z "$cli" ]]; then continue; fi
  url="https://github.com/ShubhanYenuganti/infra-press/releases/download/${cli}/${version}/${cli}-${OS}-${ARCH}"
  echo "Installing $cli $version..."
  if curl -fL "$url" -o "$INSTALL_DIR/$cli"; then
    chmod +x "$INSTALL_DIR/$cli"
    echo "  installed $cli"
  else
    echo "  failed to fetch $cli @ $version (release may not exist yet); skipping" >&2
  fi
done

echo
echo "Done. Add to PATH:"
echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
