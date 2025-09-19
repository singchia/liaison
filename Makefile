include ./Makefile.defs

all: liaison liaison-edge

.PHONY: image-gen-swagger
image-gen-swagger:
	docker buildx build -t liaison-gen-swagger:${VERSION} -f images/Dockerfile.liaison-swagger .

.PHONY: image-gen-api
image-gen-api:
	docker buildx build -t image-gen-api:${VERSION} -f images/Dockerfile.liaison-api .

# binary
.PHONY: liaison
liaison:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -o ./bin/liaison cmd/manager/main.go

.PHONY: liaison-edge
liaison-edge:
	GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o ./bin/liaison-edge cmd/edge/main.go

# api
.PHONY: gen-api
gen-api:
	docker run --rm -v ${PWD}/api/v1:/api/v1 image-gen-api:${VERSION}

.PHONY: gen-swagger
gen-swagger:
	docker run --rm -v ${PWD}:/liaison liaison-gen-swagger:${VERSION}