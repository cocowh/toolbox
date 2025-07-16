#!/bin/bash

set -e

APP_NAME="toolbox"
VERSION=${VERSION:-"v1.0.0"}
OUTPUT_DIR="./build"

DEFAULT_PLATFORMS=(
  "darwin/amd64"
  "darwin/arm64"
  "linux/amd64"
  "linux/arm64"
  "windows/amd64"
  "windows/arm64"
)

if [ "$#" -gt 0 ]; then
  PLATFORMS=("$@")
else
  PLATFORMS=("${DEFAULT_PLATFORMS[@]}")
fi

echo "🔨 Build application: $APP_NAME ($VERSION)"
rm -rf "$OUTPUT_DIR" || true
mkdir -p "$OUTPUT_DIR"

for PLATFORM in "${PLATFORMS[@]}"; do
  OS="${PLATFORM%%/*}"
  ARCH="${PLATFORM##*/}"
  BIN_NAME="${APP_NAME}-${OS}-${ARCH}"
  FILE_NAME="$BIN_NAME"
  EXT=""

  if [ "$OS" = "windows" ]; then
    EXT=".exe"
  fi

  BIN_NAME+="$EXT"

  echo "🚧 Build $OS/$ARCH -> $BIN_NAME"

  GOOS="$OS" GOARCH="$ARCH" CGO_ENABLED=0 \
    go build -ldflags "-X 'main.Version=$VERSION'" -o "$OUTPUT_DIR/$BIN_NAME" ./cmd/main

  # pack
  PACKAGE_NAME="${APP_NAME}-${VERSION}-${OS}-${ARCH}"
  cd "$OUTPUT_DIR"
  if [ "$OS" = "windows" ]; then
    zip -qr "${PACKAGE_NAME}.zip" "$BIN_NAME"
    echo "📦 Packed ${PACKAGE_NAME}.zip"
  else
    tar -czf "${PACKAGE_NAME}.tar.gz" "$BIN_NAME"
    echo "📦 Packed ${PACKAGE_NAME}.tar.gz"
  fi
  rm "$BIN_NAME"
  cd - > /dev/null
done

if command -v open >/dev/null 2>&1; then
  open "$OUTPUT_DIR"
elif command -v xdg-open >/dev/null 2>&1; then
  xdg-open "$OUTPUT_DIR"
fi

echo "✅ All builds are completed and the files are：$OUTPUT_DIR/"
