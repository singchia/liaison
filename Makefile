include ./Makefile.defs

# ============================================================================
# Default targets
# ============================================================================
.PHONY: all
all: build-local test-deps

.PHONY: linux
linux: build-linux

# ============================================================================
# Docker images for code generation
# ============================================================================
.PHONY: image-gen-swagger
image-gen-swagger:
	docker buildx build -t liaison-gen-swagger:${VERSION} -f images/Dockerfile.liaison-swagger .

.PHONY: image-gen-api
image-gen-api:
	docker buildx build -t image-gen-api:${VERSION} -f images/Dockerfile.liaison-api .

# ============================================================================
# Local development builds (for macOS)
# ============================================================================
.PHONY: build-local
build-local: build-liaison build-edge

.PHONY: build-liaison
build-liaison:
	CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -o ./bin/liaison cmd/manager/main.go

.PHONY: build-edge
build-edge:
	go build -trimpath -ldflags "-s -w" -o ./bin/liaison-edge cmd/edge/main.go

# Legacy aliases
.PHONY: liaison liaison-edge
liaison: build-liaison
liaison-edge: build-edge

# ============================================================================
# Linux builds using Docker (required for CGO cross-compilation from macOS)
# ============================================================================
.PHONY: build-linux
build-linux: build-liaison-linux build-edge-linux

.PHONY: build-liaison-linux
build-liaison-linux:
	@echo "Building liaison for Linux (amd64) using Docker..."
	@mkdir -p ./bin
	@docker run --rm \
		--platform linux/amd64 \
		-v "$(shell pwd):/build" \
		-w /build \
		-e CGO_ENABLED=1 \
		-e GOOS=linux \
		-e GOARCH=amd64 \
		-e GOTOOLCHAIN=auto \
		golang:latest sh -c "\
			apt-get update -qq && \
			DEBIAN_FRONTEND=noninteractive apt-get install -y -qq gcc libc6-dev libsqlite3-dev >/dev/null 2>&1 && \
			go env -w GOTOOLCHAIN=auto && \
			go mod download && \
			CC=gcc CGO_ENABLED=1 go build -trimpath -ldflags '-s -w' -o ./bin/liaison cmd/manager/main.go"
	@chmod +x ./bin/liaison
	@echo "✅ Built: ./bin/liaison"

.PHONY: build-edge-linux
build-edge-linux:
	@echo "Building liaison-edge for Linux (amd64) using Docker..."
	@mkdir -p ./bin
	@docker run --rm \
		--platform linux/amd64 \
		-v "$(shell pwd):/build" \
		-w /build \
		-e GOOS=linux \
		-e GOARCH=amd64 \
		-e GOTOOLCHAIN=auto \
		golang:latest sh -c "\
			go env -w GOTOOLCHAIN=auto && \
			go mod download && \
			CGO_ENABLED=0 go build -trimpath -ldflags '-s -w' -o ./bin/liaison-edge cmd/edge/main.go"
	@chmod +x ./bin/liaison-edge
	@echo "✅ Built: ./bin/liaison-edge"

# Legacy aliases
.PHONY: liaison-linux liaison-edge-linux
liaison-linux: build-liaison-linux
liaison-edge-linux: build-edge-linux

# ============================================================================
# Code generation
# ============================================================================
.PHONY: gen-api
gen-api:
	docker run --rm -v ${PWD}/api/v1:/api/v1 image-gen-api:${VERSION}

.PHONY: gen-swagger
gen-swagger:
	docker run --rm -v ${PWD}:/liaison liaison-gen-swagger:${VERSION}

.PHONY: gen
gen: gen-api gen-swagger

# ============================================================================
# Testing
# ============================================================================
.PHONY: test
test:
	go test -v ./...

.PHONY: test-deps
test-deps:
	go get github.com/stretchr/testify/assert
	go get github.com/stretchr/testify/require

.PHONY: test-e2e
test-e2e:
	@echo "Running E2E tests..."
	@./test/e2e/run_simple_test.sh

.PHONY: test-e2e-full
test-e2e-full:
	@echo "Running full E2E tests..."
	@./test/e2e/run_e2e_test.sh

.PHONY: test-all
test-all: test test-e2e

# ============================================================================
# Tools
# ============================================================================
.PHONY: build-tools
build-tools: build-password-verifier build-password-generator

.PHONY: build-password-verifier
build-password-verifier:
	CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -o ./bin/password-verifier ./tools/password-verifier

.PHONY: build-password-generator
build-password-generator:
	CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -o ./bin/password-generator ./tools/password-generator

.PHONY: build-tools-linux
build-tools-linux: build-password-verifier-linux build-password-generator-linux

.PHONY: build-password-verifier-linux
build-password-verifier-linux:
	@echo "Building password-verifier for Linux..."
	@mkdir -p ./bin
	@docker run --rm \
		--platform linux/amd64 \
		-v "$(shell pwd):/build" \
		-w /build \
		-e CGO_ENABLED=1 \
		-e GOOS=linux \
		-e GOARCH=amd64 \
		-e GOTOOLCHAIN=auto \
		golang:latest sh -c "\
			apt-get update -qq && \
			DEBIAN_FRONTEND=noninteractive apt-get install -y -qq gcc libc6-dev libsqlite3-dev >/dev/null 2>&1 && \
			go env -w GOTOOLCHAIN=auto && \
			go mod download && \
			CC=gcc CGO_ENABLED=1 go build -trimpath -ldflags '-s -w' -o ./bin/password-verifier ./tools/password-verifier"

.PHONY: build-password-generator-linux
build-password-generator-linux:
	@echo "Building password-generator for Linux..."
	@mkdir -p ./bin
	@docker run --rm \
		--platform linux/amd64 \
		-v "$(shell pwd):/build" \
		-w /build \
		-e CGO_ENABLED=1 \
		-e GOOS=linux \
		-e GOARCH=amd64 \
		-e GOTOOLCHAIN=auto \
		golang:latest sh -c "\
			apt-get update -qq && \
			DEBIAN_FRONTEND=noninteractive apt-get install -y -qq gcc libc6-dev libsqlite3-dev >/dev/null 2>&1 && \
			go env -w GOTOOLCHAIN=auto && \
			go mod download && \
			CC=gcc CGO_ENABLED=1 go build -trimpath -ldflags '-s -w' -o ./bin/password-generator ./tools/password-generator"

# Legacy aliases
.PHONY: tools tools-linux password-verifier password-generator password-verifier-linux password-generator-linux
tools: build-tools
tools-linux: build-tools-linux
password-verifier: build-password-verifier
password-generator: build-password-generator
password-verifier-linux: build-password-verifier-linux
password-generator-linux: build-password-generator-linux

# ============================================================================
# Edge binary builds for different architectures
# ============================================================================
.PHONY: build-edge-all
build-edge-all: build-edge-linux-amd64 build-edge-linux-arm64 build-edge-darwin-amd64 build-edge-darwin-arm64 build-edge-windows-amd64
	@echo "✅ All edge binaries built in ./bin/"

.PHONY: build-edge-linux-amd64
build-edge-linux-amd64:
	@echo "Building liaison-edge for linux-amd64..."
	@mkdir -p ./bin
	@docker run --rm \
		--platform linux/amd64 \
		-v "$(shell pwd):/build" \
		-w /build \
		-e GOOS=linux \
		-e GOARCH=amd64 \
		-e GOTOOLCHAIN=auto \
		golang:latest sh -c "\
			go env -w GOTOOLCHAIN=auto && \
			go mod download && \
			CGO_ENABLED=0 go build -trimpath -ldflags '-s -w' -o ./bin/liaison-edge-linux-amd64 cmd/edge/main.go"
	@chmod +x ./bin/liaison-edge-linux-amd64
	@echo "✅ Built: ./bin/liaison-edge-linux-amd64"

.PHONY: build-edge-linux-arm64
build-edge-linux-arm64:
	@echo "Building liaison-edge for linux-arm64..."
	@mkdir -p ./bin
	@docker run --rm \
		--platform linux/arm64 \
		-v "$(shell pwd):/build" \
		-w /build \
		-e GOOS=linux \
		-e GOARCH=arm64 \
		-e GOTOOLCHAIN=auto \
		golang:latest sh -c "\
			go env -w GOTOOLCHAIN=auto && \
			go mod download && \
			CGO_ENABLED=0 go build -trimpath -ldflags '-s -w' -o ./bin/liaison-edge-linux-arm64 cmd/edge/main.go"
	@chmod +x ./bin/liaison-edge-linux-arm64
	@echo "✅ Built: ./bin/liaison-edge-linux-arm64"

.PHONY: build-edge-darwin-amd64
build-edge-darwin-amd64:
	@echo "Building liaison-edge for darwin-amd64..."
	@mkdir -p ./bin
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o ./bin/liaison-edge-darwin-amd64 cmd/edge/main.go
	@chmod +x ./bin/liaison-edge-darwin-amd64
	@echo "✅ Built: ./bin/liaison-edge-darwin-amd64"

.PHONY: build-edge-darwin-arm64
build-edge-darwin-arm64:
	@echo "Building liaison-edge for darwin-arm64..."
	@mkdir -p ./bin
	@CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags "-s -w" -o ./bin/liaison-edge-darwin-arm64 cmd/edge/main.go
	@chmod +x ./bin/liaison-edge-darwin-arm64
	@echo "✅ Built: ./bin/liaison-edge-darwin-arm64"

.PHONY: build-edge-windows-amd64
build-edge-windows-amd64:
	@echo "Building liaison-edge for windows-amd64..."
	@mkdir -p ./bin
	@docker run --rm \
		--platform linux/amd64 \
		-v "$(shell pwd):/build" \
		-w /build \
		-e GOOS=windows \
		-e GOARCH=amd64 \
		-e GOTOOLCHAIN=auto \
		golang:latest sh -c "\
			go env -w GOTOOLCHAIN=auto && \
			go mod download && \
			CGO_ENABLED=0 go build -trimpath -ldflags '-s -w' -o ./bin/liaison-edge-windows-amd64.exe cmd/edge/main.go"
	@echo "✅ Built: ./bin/liaison-edge-windows-amd64.exe"

# ============================================================================
# Edge package creation (tar.gz with binary and config template)
# ============================================================================
.PHONY: package-edge-all
package-edge-all: package-edge-linux-amd64 package-edge-linux-arm64 package-edge-darwin-amd64 package-edge-darwin-arm64 package-edge-windows-amd64
	@echo "✅ All edge packages created in ./packages/edge/"

.PHONY: package-edge-linux-amd64
package-edge-linux-amd64: build-edge-linux-amd64
	@echo "Packaging liaison-edge-linux-amd64..."
	@mkdir -p ./packages/edge
	@TMP_DIR=$$(mktemp -d) && \
		cp ./bin/liaison-edge-linux-amd64 $$TMP_DIR/liaison-edge && \
		cp ./dist/edge/liaison-edge.yaml.template $$TMP_DIR/liaison-edge.yaml.template && \
		cd $$TMP_DIR && \
		tar -czf ../../packages/edge/liaison-edge-linux-amd64.tar.gz liaison-edge liaison-edge.yaml.template && \
		rm -rf $$TMP_DIR && \
		echo "✅ Package created: ./packages/edge/liaison-edge-linux-amd64.tar.gz"

.PHONY: package-edge-linux-arm64
package-edge-linux-arm64: build-edge-linux-arm64
	@echo "Packaging liaison-edge-linux-arm64..."
	@mkdir -p ./packages/edge
	@TMP_DIR=$$(mktemp -d) && \
		cp ./bin/liaison-edge-linux-arm64 $$TMP_DIR/liaison-edge && \
		cp ./dist/edge/liaison-edge.yaml.template $$TMP_DIR/liaison-edge.yaml.template && \
		cd $$TMP_DIR && \
		tar -czf ../../packages/edge/liaison-edge-linux-arm64.tar.gz liaison-edge liaison-edge.yaml.template && \
		rm -rf $$TMP_DIR && \
		echo "✅ Package created: ./packages/edge/liaison-edge-linux-arm64.tar.gz"

.PHONY: package-edge-darwin-amd64
package-edge-darwin-amd64: build-edge-darwin-amd64
	@echo "Packaging liaison-edge-darwin-amd64..."
	@mkdir -p ./packages/edge
	@TMP_DIR=$$(mktemp -d) && \
		cp ./bin/liaison-edge-darwin-amd64 $$TMP_DIR/liaison-edge && \
		cp ./dist/edge/liaison-edge.yaml.template $$TMP_DIR/liaison-edge.yaml.template && \
		cd $$TMP_DIR && \
		tar -czf ../../packages/edge/liaison-edge-darwin-amd64.tar.gz liaison-edge liaison-edge.yaml.template && \
		rm -rf $$TMP_DIR && \
		echo "✅ Package created: ./packages/edge/liaison-edge-darwin-amd64.tar.gz"

.PHONY: package-edge-darwin-arm64
package-edge-darwin-arm64: build-edge-darwin-arm64
	@echo "Packaging liaison-edge-darwin-arm64..."
	@mkdir -p ./packages/edge
	@TMP_DIR=$$(mktemp -d) && \
		cp ./bin/liaison-edge-darwin-arm64 $$TMP_DIR/liaison-edge && \
		cp ./dist/edge/liaison-edge.yaml.template $$TMP_DIR/liaison-edge.yaml.template && \
		cd $$TMP_DIR && \
		tar -czf ../../packages/edge/liaison-edge-darwin-arm64.tar.gz liaison-edge liaison-edge.yaml.template && \
		rm -rf $$TMP_DIR && \
		echo "✅ Package created: ./packages/edge/liaison-edge-darwin-arm64.tar.gz"

.PHONY: package-edge-windows-amd64
package-edge-windows-amd64: build-edge-windows-amd64
	@echo "Packaging liaison-edge-windows-amd64..."
	@mkdir -p ./packages/edge
	@TMP_DIR=$$(mktemp -d) && \
		cp ./bin/liaison-edge-windows-amd64.exe $$TMP_DIR/liaison-edge.exe && \
		cp ./dist/edge/liaison-edge.yaml.template $$TMP_DIR/liaison-edge.yaml.template && \
		cd $$TMP_DIR && \
		tar -czf ../../packages/edge/liaison-edge-windows-amd64.tar.gz liaison-edge.exe liaison-edge.yaml.template && \
		rm -rf $$TMP_DIR && \
		echo "✅ Package created: ./packages/edge/liaison-edge-windows-amd64.tar.gz"

# Legacy aliases
.PHONY: edge-packages edge-package-linux-amd64 edge-package-linux-arm64 edge-package-darwin-amd64 edge-package-darwin-arm64 edge-package-windows-amd64
edge-packages: package-edge-all
edge-package-linux-amd64: package-edge-linux-amd64
edge-package-linux-arm64: package-edge-linux-arm64
edge-package-darwin-amd64: package-edge-darwin-amd64
edge-package-darwin-arm64: package-edge-darwin-arm64
edge-package-windows-amd64: package-edge-windows-amd64

# ============================================================================
# Full package (liaison + edge + systemd files)
# ============================================================================
.PHONY: package
package: build-linux
	@echo "Packaging Liaison for Linux installation..."
	@PACK_DIR=liaison-$(VERSION)-linux-amd64 && \
	rm -rf $$PACK_DIR && \
	mkdir -p $$PACK_DIR/bin $$PACK_DIR/etc $$PACK_DIR/systemd && \
	cp bin/liaison $$PACK_DIR/bin/ && \
	cp bin/liaison-edge $$PACK_DIR/bin/ && \
	cp etc/liaison.yaml $$PACK_DIR/etc/ && \
	cp etc/liaison-edge.yaml $$PACK_DIR/etc/ && \
	cp dist/systemd/liaison.service $$PACK_DIR/systemd/ && \
	cp dist/systemd/install.sh $$PACK_DIR/ && \
	cp dist/systemd/uninstall.sh $$PACK_DIR/ && \
	cp dist/systemd/README.md $$PACK_DIR/systemd/ 2>/dev/null || true && \
	cp VERSION $$PACK_DIR/ && \
	chmod +x $$PACK_DIR/bin/* && \
	chmod +x $$PACK_DIR/install.sh && \
	chmod +x $$PACK_DIR/uninstall.sh && \
	tar -czf $$PACK_DIR.tar.gz $$PACK_DIR && \
	rm -rf $$PACK_DIR && \
	echo "✅ Package created: $$PACK_DIR.tar.gz"

# Legacy alias
.PHONY: pack
pack: package
