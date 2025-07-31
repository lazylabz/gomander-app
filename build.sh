#!/bin/bash

# Clean up the build directory
rm -rf build/bin

# First, we'll create the windows and darwin builds
cd cmd/gomander
wails build --platform="darwin/amd64" && mv ../../build/bin/gomander.app ../../build/bin/gomander-darwin-amd64.app
wails build --platform="darwin/arm64" && mv ../../build/bin/gomander.app ../../build/bin/gomander-darwin-arm64.app
wails build -nsis --platform="windows/amd64" -o="gomander-windows-amd64-portable.exe" && mv ../../build/bin/gomander-amd64-installer.exe ../../build/bin/gomander-windows-amd64-installer.exe
wails build -nsis --platform="windows/arm64" -o="gomander-windows-arm64-portable.exe" && mv ../../build/bin/gomander-arm64-installer.exe ../../build/bin/gomander-windows-arm64-installer.exe
cd ../../

# Then, we'll build the Mac OS dmg installers

# Create the dmg for macOS amd64
mv build/bin/gomander-darwin-amd64.app build/bin/gomander.app
create-dmg \
  --volname "Gomander Installer (amd64)" \
  --window-pos 200 120 \
  --window-size 500 300 \
  --icon-size 100 \
  --icon "gomander.app" 100 100 \
  --app-drop-link 380 100 \
  --hide-extension "gomander.app" \
  "build/bin/gomander-darwin-amd64.dmg" \
  "build/bin/gomander.app"
mv build/bin/gomander.app build/bin/gomander-darwin-amd64.app

# Create the dmg for macOS arm64
mv build/bin/gomander-darwin-arm64.app build/bin/gomander.app
create-dmg \
  --volname "Gomander Installer (arm64)" \
  --window-pos 200 120 \
  --window-size 500 300 \
  --icon-size 100 \
  --icon "gomander.app" 100 100 \
  --app-drop-link 380 100 \
  --hide-extension "gomander.app" \
  "build/bin/gomander-darwin-arm64.dmg" \
  "build/bin/gomander.app"
mv build/bin/gomander.app build/bin/gomander-darwin-arm64.app

# Zip the .app installers
zip -r build/bin/gomander-darwin-amd64.app.zip build/bin/gomander-darwin-amd64.app
zip -r build/bin/gomander-darwin-arm64.app.zip build/bin/gomander-darwin-arm64.app

# Remove the original app directories
rm -rf build/bin/gomander-darwin-amd64.app
rm -rf build/bin/gomander-darwin-arm64.app

# Finally, we'll build the Linux binaries

# Prepare the docker image for building the Linux binary
docker buildx build --platform linux/amd64,linux/arm64 --output "type=docker" -t wails-linux-builder .

# Run the command to build the Linux binary
docker run --rm --platform linux/amd64 -v $(pwd)/build/bin:/app/output wails-linux-builder
docker run --rm --platform linux/arm64 -v $(pwd)/build/bin:/app/output wails-linux-builder