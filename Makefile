# Specify the GOOS and GOARCH variables
export CGO_ENABLED=0
export GOARCH=amd64
export GOOS=darwin

TARGET=esmsgo

# Remove debug information for release builds
ifeq ($(MAKECMDGOALS), build)
	OUTPUT_DIR=bin/debug
	LDFLAGS=""
else 
	LDFLAGS="-s -w"
	OUTPUT_DIR=bin/release
endif

build:
	@echo "Building for $(GOOS)_$(GOARCH)..." && \
	FILE_EXT=""; \
	if [ "$$GOOS" = "windows" ]; then FILE_EXT=".exe"; fi; \
	BUILD_DIR=$(OUTPUT_DIR)/$(GOOS)/$(GOARCH) && \
	OUTPUT_FILE=$$BUILD_DIR/$(TARGET)$$FILE_EXT && \
	echo "Output file: $$OUTPUT_FILE" && \
	rm -rf $$BUILD_DIR && \
	go mod tidy && \
	git diff --exit-code go.mod go.sum && \
	go mod download && \
	go mod verify && \
	go build -ldflags=$(LDFLAGS) -tags $(TARGET) -o $$OUTPUT_FILE ./cmd/cli

package:
	@PKG_DIR=$(OUTPUT_DIR)/$(GOOS)/$(GOARCH) && \
	echo "Packaging $$PKG_DIR..." && \
	if [ "$(GOOS)" = "linux" ]; then \
		tar -czvf "$(TARGET)_$(VERSION)_$(GOOS)_$(GOARCH).tgz" -C "$$PKG_DIR" .; \
	else \
		zip -rj "$$PKG_DIR/$(TARGET)_$(VERSION)_$(GOOS)_$(GOARCH).zip" "$$PKG_DIR"; \
	fi

build-all:
	@$(MAKE) GOOS=windows build
	@$(MAKE) GOOS=windows GOARCH=386 build
	@$(MAKE) GOOS=linux build
	@$(MAKE) GOOS=darwin build

win:
	@$(MAKE) GOOS=windows build package

win32:
	@$(MAKE) GOOS=windows GOARCH=386 build package

linux:
	@$(MAKE) GOOS=linux build package

mac:
	@$(MAKE) GOOS=darwin build package