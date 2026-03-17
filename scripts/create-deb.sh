#!/bin/bash
# Creates a .deb package from a pre-built Gomander binary.
# Usage: ./scripts/create-deb.sh <binary_path> <arch> [version]
#   arch: amd64 or arm64
#   version: defaults to reading from internal/releases/release.go

set -e

BINARY="$1"
ARCH="$2"
VERSION="${3:-$(grep -oP '(?<=CurrentRelease = "v)[^"]+' internal/releases/release.go 2>/dev/null || echo "0.0.0")}"

if [ -z "$BINARY" ] || [ -z "$ARCH" ]; then
    echo "Usage: $0 <binary_path> <arch> [version]"
    echo "  arch: amd64 or arm64"
    exit 1
fi

if [ ! -f "$BINARY" ]; then
    echo "Error: binary not found at '$BINARY'"
    exit 1
fi

DEB_NAME="gomander_${VERSION}_${ARCH}"
DEB_STAGING="/tmp/${DEB_NAME}"
OUTPUT_DIR="$(dirname "$BINARY")"

echo "Packaging gomander v${VERSION} (${ARCH}) -> ${OUTPUT_DIR}/${DEB_NAME}.deb"

# Cleanup any previous staging directory
rm -rf "${DEB_STAGING}"

# Create .deb directory structure
mkdir -p "${DEB_STAGING}/DEBIAN"
mkdir -p "${DEB_STAGING}/usr/bin"
mkdir -p "${DEB_STAGING}/usr/share/applications"
mkdir -p "${DEB_STAGING}/usr/share/icons/hicolor/256x256/apps"

# Copy binary
cp "${BINARY}" "${DEB_STAGING}/usr/bin/gomander"
chmod 755 "${DEB_STAGING}/usr/bin/gomander"

# Copy application icon (relative to project root)
ICON_SRC="build/appicon.png"
if [ -f "${ICON_SRC}" ]; then
    cp "${ICON_SRC}" "${DEB_STAGING}/usr/share/icons/hicolor/256x256/apps/gomander.png"
fi

# Create .desktop entry
cat > "${DEB_STAGING}/usr/share/applications/gomander.desktop" << 'EOF'
[Desktop Entry]
Name=Gomander
Comment=Launch, monitor, and organize commands
Exec=/usr/bin/gomander
Icon=gomander
Terminal=false
Type=Application
Categories=Development;Utility;
StartupWMClass=gomander
EOF

# Calculate installed size in KB
INSTALLED_SIZE=$(du -sk "${DEB_STAGING}" | cut -f1)

# Create DEBIAN/control
cat > "${DEB_STAGING}/DEBIAN/control" << EOF
Package: gomander
Version: ${VERSION}
Architecture: ${ARCH}
Maintainer: Gomander <noreply@gomander.app>
Installed-Size: ${INSTALLED_SIZE}
Depends: libgtk-3-0, libwebkit2gtk-4.1-0 | libwebkit2gtk-4.0-37
Description: Launch, monitor, and organize commands
 Gomander is a GUI application to help developers manage multiple project
 commands without juggling multiple terminal windows.
EOF

# Build the .deb package
dpkg-deb --build --root-owner-group "${DEB_STAGING}" "${OUTPUT_DIR}/${DEB_NAME}.deb"

# Cleanup staging directory
rm -rf "${DEB_STAGING}"

echo "Done: ${OUTPUT_DIR}/${DEB_NAME}.deb"
