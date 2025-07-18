include ./Makefile.defs

.PHONY: image-gen-swagger
image-gen-swagger:
	docker buildx build -t liaison-gen-swagger:${VERSION} -f images/Dockerfile.liaison-swagger .

.PHONY: image-gen-api
image-gen-api:
	docker buildx build -t image-gen-api:${VERSION} -f images/Dockerfile.liaison-api .

# api
.PHONY: gen-api
gen-api:
	docker run --rm -v ${PWD}/api/v1:/api/v1 image-gen-api:${VERSION}

.PHONY: gen-swagger
gen-swagger:
	docker run --rm -v ${PWD}:/liaison liaison-gen-swagger:${VERSION}