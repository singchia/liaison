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

# Linux builds without CGO (for cross-compilation)
.PHONY: liaison-linux
liaison-linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -o ./bin/liaison cmd/manager/main.go

.PHONY: liaison-edge-linux
liaison-edge-linux:
	GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o ./bin/liaison-edge cmd/edge/main.go

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
	CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -o ./bin/password-verifier tools/password_verifier.go

.PHONY: tools
tools: password-verifier