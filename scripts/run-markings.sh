#!/usr/bin/env bash
set -e

# Determine version from the pre-commit hook's git tag
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
VERSION=$(git -C "$SCRIPT_DIR" describe --tags --exact-match 2>/dev/null || echo "")

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
    
    echo "Downloading markings ${VERSION} for ${OS}-${ARCH}..." >&2
    mkdir -p "$BIN_DIR"
    
    if [ "$EXT" = "zip" ]; then
        curl -sL -o "$BIN_DIR/archive.zip" "$DOWNLOAD_URL"
        unzip -q -o "$BIN_DIR/archive.zip" "markings.exe" -d "$BIN_DIR"
        rm "$BIN_DIR/archive.zip"
    else
        curl -sL "$DOWNLOAD_URL" | tar -xz -C "$BIN_DIR" markings
    fi
    
    chmod +x "$MARKINGS_BIN"
fi

exec "$MARKINGS_BIN" "$@"
