# Gomander Build Makefile
.PHONY: all clean windows darwin linux dmg docker-build help

# Variables
BUILD_DIR = build/bin
CMD_DIR = cmd/gomander
DOCKER_IMAGE = wails-linux-builder

# Default target
all: clean darwin windows linux

# Help target
help:
	@echo "Available targets:"
	@echo "  all          - Build all platforms (darwin, windows, linux)"
	@echo "  clean        - Clean build directory"
	@echo "  darwin       - Build macOS binaries and create DMG installers"
	@echo "  darwin-amd64 - Build macOS AMD64 binary only"
	@echo "  darwin-arm64 - Build macOS ARM64 binary only"
	@echo "  dmg          - Create DMG installers from existing .app files"
	@echo "  dmg-amd64    - Create DMG installer for AMD64"
	@echo "  dmg-arm64    - Create DMG installer for ARM64"
	@echo "  windows      - Build Windows binaries (both architectures)"
	@echo "  windows-amd64 - Build Windows AMD64 binary"
	@echo "  windows-arm64 - Build Windows ARM64 binary"
	@echo "  linux        - Build Linux binaries (both architectures)"
	@echo "  linux-amd64  - Build Linux AMD64 binary"
	@echo "  linux-arm64  - Build Linux ARM64 binary"
	@echo "  docker-build - Build Docker image for Linux builds"

# Clean build directory
clean:
	rm -rf $(BUILD_DIR)

# Darwin (macOS) targets
darwin: darwin-amd64 darwin-arm64 dmg zip-darwin cleanup-darwin

darwin-amd64:
	cd $(CMD_DIR) && \
	wails build --platform="darwin/amd64" && \
	mv ../../$(BUILD_DIR)/gomander.app ../../$(BUILD_DIR)/gomander-darwin-amd64.app

darwin-arm64:
	cd $(CMD_DIR) && \
	wails build --platform="darwin/arm64" && \
	mv ../../$(BUILD_DIR)/gomander.app ../../$(BUILD_DIR)/gomander-darwin-arm64.app

# DMG creation targets
dmg: dmg-amd64 dmg-arm64

dmg-amd64:
	mv $(BUILD_DIR)/gomander-darwin-amd64.app $(BUILD_DIR)/gomander.app
	create-dmg \
		--volname "Gomander Installer (amd64)" \
		--window-pos 200 120 \
		--window-size 500 300 \
		--icon-size 100 \
		--icon "gomander.app" 100 100 \
		--app-drop-link 380 100 \
		--hide-extension "gomander.app" \
		"$(BUILD_DIR)/gomander-darwin-amd64.dmg" \
		"$(BUILD_DIR)/gomander.app"
	mv $(BUILD_DIR)/gomander.app $(BUILD_DIR)/gomander-darwin-amd64.app

dmg-arm64:
	mv $(BUILD_DIR)/gomander-darwin-arm64.app $(BUILD_DIR)/gomander.app
	create-dmg \
		--volname "Gomander Installer (arm64)" \
		--window-pos 200 120 \
		--window-size 500 300 \
		--icon-size 100 \
		--icon "gomander.app" 100 100 \
		--app-drop-link 380 100 \
		--hide-extension "gomander.app" \
		"$(BUILD_DIR)/gomander-darwin-arm64.dmg" \
		"$(BUILD_DIR)/gomander.app"
	mv $(BUILD_DIR)/gomander.app $(BUILD_DIR)/gomander-darwin-arm64.app

# Zip Darwin .app files
zip-darwin:
	zip -r $(BUILD_DIR)/gomander-darwin-amd64.app.zip $(BUILD_DIR)/gomander-darwin-amd64.app
	zip -r $(BUILD_DIR)/gomander-darwin-arm64.app.zip $(BUILD_DIR)/gomander-darwin-arm64.app

# Clean up Darwin .app directories
cleanup-darwin:
	rm -rf $(BUILD_DIR)/gomander-darwin-amd64.app
	rm -rf $(BUILD_DIR)/gomander-darwin-arm64.app

# Windows targets
windows: windows-amd64 windows-arm64

windows-amd64:
	cd $(CMD_DIR) && \
	wails build -nsis --platform="windows/amd64" -o="gomander-windows-amd64-portable.exe" && \
	mv ../../$(BUILD_DIR)/gomander-amd64-installer.exe ../../$(BUILD_DIR)/gomander-windows-amd64-installer.exe

windows-arm64:
	cd $(CMD_DIR) && \
	wails build -nsis --platform="windows/arm64" -o="gomander-windows-arm64-portable.exe" && \
	mv ../../$(BUILD_DIR)/gomander-arm64-installer.exe ../../$(BUILD_DIR)/gomander-windows-arm64-installer.exe

# Linux targets
linux: docker-build linux-amd64 linux-arm64

docker-build:
	docker buildx build --platform linux/amd64,linux/arm64 --output "type=docker" -t $(DOCKER_IMAGE) .

linux-amd64: docker-build
	docker run --rm --platform linux/amd64 -v $$(pwd)/$(BUILD_DIR):/app/output $(DOCKER_IMAGE)

linux-arm64: docker-build
	docker run --rm --platform linux/arm64 -v $$(pwd)/$(BUILD_DIR):/app/output $(DOCKER_IMAGE)

# Individual platform builds (without dependencies)
build-darwin-only: clean darwin-amd64 darwin-arm64
build-windows-only: clean windows-amd64 windows-arm64
build-linux-only: clean linux-amd64 linux-arm64