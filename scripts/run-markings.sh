#!/usr/bin/env bash
set -e

# Determine version from the pre-commit hook's git tag
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
git -C "$SCRIPT_DIR" fetch --tags --quiet 2>/dev/null || true
VERSION=$(git -C "$SCRIPT_DIR" describe --tags --match "v*.*" --exact-match 2>/dev/null || echo "")

if [ -z "$VERSION" ]; then
    echo "markings: Error: Could not determine release version. Ensure you are using a tagged release in your .pre-commit-config.yaml." >&2
    exit 1
fi

OS=$(uname -s)
ARCH=$(uname -m)
EXT="tar.gz"

case "$OS" in
    Linux)
        OS="Linux"
        ;;
    Darwin)
        OS="Darwin"
        ;;
    MINGW* | MSYS* | CYGWIN*)
        OS="Windows"
        EXT="zip"
        ;;
    *)
        echo "markings: Unsupported OS: $OS" >&2
        exit 1
        ;;
esac

case "$ARCH" in
    x86_64 | amd64)
        ARCH="x86_64"
        ;;
    aarch64 | arm64)
        ARCH="arm64"
        ;;
    *)
        echo "markings: Unsupported architecture: $ARCH" >&2
        exit 1
        ;;
esac

CACHE_DIR="${XDG_CACHE_HOME:-$HOME/.cache}/markings"
if [ "$OS" = "Windows" ]; then
    CACHE_DIR="${LOCALAPPDATA:-$HOME/AppData/Local}/markings/cache"
fi

BIN_DIR="$CACHE_DIR/$VERSION"
MARKINGS_BIN="$BIN_DIR/markings"

if [ "$OS" = "Windows" ]; then
    MARKINGS_BIN="${MARKINGS_BIN}.exe"
fi

if [ ! -x "$MARKINGS_BIN" ]; then
    DOWNLOAD_URL="https://github.com/mharrisb1/markings/releases/download/${VERSION}/markings_${OS}_${ARCH}.${EXT}"
    
    # Use a lock-free atomic download approach to prevent "Text file busy"
    # when pre-commit runs multiple hooks in parallel.
    TMP_DIR="$BIN_DIR/tmp.$$"
    mkdir -p "$TMP_DIR"
    
    echo "Downloading markings ${VERSION} for ${OS}-${ARCH}..." >&2
    
    if [ "$EXT" = "zip" ]; then
        if ! curl -sfL -o "$TMP_DIR/archive.zip" "$DOWNLOAD_URL"; then
            echo "markings: Error: Failed to download from $DOWNLOAD_URL" >&2
            rm -rf "$TMP_DIR"
            exit 1
        fi
        unzip -q -o "$TMP_DIR/archive.zip" "markings.exe" -d "$TMP_DIR"
        chmod +x "$TMP_DIR/markings.exe"
        mv -f "$TMP_DIR/markings.exe" "$MARKINGS_BIN" 2>/dev/null || true
    else
        if ! curl -sfL "$DOWNLOAD_URL" | tar -xz -C "$TMP_DIR" markings; then
            echo "markings: Error: Failed to download or extract from $DOWNLOAD_URL" >&2
            rm -rf "$TMP_DIR"
            exit 1
        fi
        chmod +x "$TMP_DIR/markings"
        mv -f "$TMP_DIR/markings" "$MARKINGS_BIN" 2>/dev/null || true
    fi
    
    rm -rf "$TMP_DIR"
fi

exec "$MARKINGS_BIN" "$@"
