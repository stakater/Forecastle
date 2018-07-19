.PHONY: install test build binary-image push

BUILDER ?= forecastle-builder
DOCKER_IMAGE ?= stakater/forecastle

# Default value "dev"
DOCKER_TAG ?= dev
REPOSITORY = ${DOCKER_IMAGE}:${DOCKER_TAG}

GOCMD = go
GLIDECMD = glide

install:
	"$(GLIDECMD)" install

test:
	"$(GOCMD)" test -v ./...

build:

builder-image:
	@docker build --network host -t "${BUILDER}" -f build/package/Dockerfile.build .

binary-image: builder-image
	@docker run --network host --rm "${BUILDER}" | docker build --network host -t "${REPOSITORY}" -f Dockerfile.run -

push:
	docker push $(REPOSITORY)