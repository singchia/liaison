include ./Makefile.defs

.PHONY: image-gen-api
image-gen-api:
	docker buildx build -t image-gen-api:${VERSION} -f images/Dockerfile.liaison-api .

# api
.PHONY: gen-api
gen-api:
	docker run --rm -v ${PWD}/api/v1:/api/v1 image-gen-api:${VERSION}