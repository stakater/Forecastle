.PHONY: install test build binary-image push

DOCKER_IMAGE ?= stakater/forecastle

# Default value "dev"
DOCKER_TAG ?= dev
REPOSITORY = ${DOCKER_IMAGE}:${DOCKER_TAG}

install:
	cd src && npm install

test:

build:

binary-image:
	docker build --network host -t ${REPOSITORY} .

push:
	docker push $(REPOSITORY)