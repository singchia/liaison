include ./Makefile.defs

all: liaison liaison-edge test-deps
linux: liaison-linux liaison-edge-linux

.PHONY: image-gen-swagger
image-gen-swagger:
	docker buildx build -t liaison-gen-swagger:${VERSION} -f images/Dockerfile.liaison-swagger .

.PHONY: image-gen-api
image-gen-api:
	docker buildx build -t image-gen-api:${VERSION} -f images/Dockerfile.liaison-api .

# local development builds (for macOS)
.PHONY: liaison
liaison:
	CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -o ./bin/liaison cmd/manager/main.go

.PHONY: liaison-edge
liaison-edge:
	go build -trimpath -ldflags "-s -w" -o ./bin/liaison-edge cmd/edge/main.go

# Linux builds using Docker (required for CGO cross-compilation from macOS)
.PHONY: liaison-linux
liaison-linux:
	@echo "Building Linux binaries using Docker..."
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
			CC=gcc CGO_ENABLED=1 go build -trimpath -ldflags '-s -w' -o ./bin/liaison cmd/manager/main.go && \
			CGO_ENABLED=0 go build -trimpath -ldflags '-s -w' -o ./bin/liaison-edge cmd/edge/main.go"
	@chmod +x ./bin/liaison ./bin/liaison-edge
	@echo "✅ Linux binaries built in ./bin/"

.PHONY: liaison-edge-linux
liaison-edge-linux: liaison-linux

# api
.PHONY: gen-api
gen-api:
	docker run --rm -v ${PWD}/api/v1:/api/v1 image-gen-api:${VERSION}

.PHONY: gen-swagger
gen-swagger:
	docker run --rm -v ${PWD}:/liaison liaison-gen-swagger:${VERSION}

# testing
.PHONY: test
test:
	go test -v ./...

.PHONY: test-e2e
test-e2e:
	@echo "Running E2E tests..."
	@./test/e2e/run_simple_test.sh

.PHONY: test-e2e-full
test-e2e-full:
	@echo "Running full E2E tests..."
	@./test/e2e/run_e2e_test.sh

.PHONY: test-deps
test-deps:
	go get github.com/stretchr/testify/assert
	go get github.com/stretchr/testify/require

.PHONY: test-all
test-all: test test-e2e

# tools
.PHONY: password-verifier
password-verifier:
	CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -o ./bin/password-verifier ./tools/password-verifier

.PHONY: password-generator
password-generator:
	CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -o ./bin/password-generator ./tools/password-generator

.PHONY: password-verifier-linux
password-verifier-linux:
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

.PHONY: password-generator-linux
password-generator-linux:
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

.PHONY: tools
tools: password-verifier password-generator

.PHONY: tools-linux
tools-linux: password-verifier-linux password-generator-linux

# packaging
.PHONY: pack
pack: liaison-linux
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
	echo "Package created: $$PACK_DIR.tar.gz"