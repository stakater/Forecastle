.PHONY: default build builder-image binary-image test stop clean-images clean push apply deploy release release-all manifest push clean-image

OS ?= linux
ARCH ?= ???
ALL_ARCH ?= arm64 arm amd64

BUILDER ?= forecastle-builder
BINARY ?= Forecastle
DOCKER_IMAGE ?= stakater/forecastle

# Default value "dev"
DOCKER_TAG ?= dev
REPOSITORY_GENERIC = ${DOCKER_IMAGE}:${DOCKER_TAG}
REPOSITORY_ARCH = ${DOCKER_IMAGE}:${TAG}-${ARCH}

VERSION ?= 0.0.1
BUILD=

GOCMD = go
GOFLAGS ?= $(GOFLAGS:)
LDFLAGS =

default: build test

install:
	"$(GOCMD)" mod download

build:
	"$(GOCMD)" build ${GOFLAGS} ${LDFLAGS} -o "${BINARY}"

builder-image:
	docker buildx build --platform ${OS}/${ARCH} --build-arg GOARCH=$(ARCH) -t "${BUILDER}" --load -f build/package/Dockerfile.build .

forecastle-${ARCH}.tar:
	docker buildx build --platform ${OS}/${ARCH} --build-arg GOARCH=$(ARCH) -t "${BUILDER}" --load -f build/package/Dockerfile.build .
	docker run --platform ${OS}/${ARCH} --rm "${BUILDER}" > forecastle-${ARCH}.tar

binary-image: builder-image
	cat forecastle-${ARCH}.tar | docker buildx build --platform ${OS}/${ARCH} -t "${REPOSITORY_ARCH}"  --load -f Dockerfile.run -

push:
	docker push ${REPOSITORY_ARCH}

release:  binary-image push manifest
release-all:
	-rm -rf ~/.docker/manifests/*
	# Make arch-specific release
	@for arch in $(ALL_ARCH) ; do \
		echo Make release: $$arch ; \
		make release ARCH=$$arch ; \
	done

	set -e
	docker manifest push --purge $(REPOSITORY_GENERIC)

manifest:
	set -e
	docker manifest create -a $(REPOSITORY_GENERIC) $(REPOSITORY_ARCH)
	docker manifest annotate --arch $(ARCH) $(REPOSITORY_GENERIC)  $(REPOSITORY_ARCH)

test:
	"$(GOCMD)" test -timeout 1800s -v ./...

stop:
	@docker stop "${BINARY}"

clean-images: stop
	-docker rmi "${BINARY}"
	@for arch in $(ALL_ARCH) ; do \
		echo Clean image: $$arch ; \
		make clean-image ARCH=$$arch ; \
	done
	-docker rmi "${REPOSITORY_GENERIC}"

clean-image:
	-docker rmi "${BUILDER}"
	-docker rmi "${REPOSITORY_ARCH}"
	-rm -rf ~/.docker/manifests/*

clean:
	"$(GOCMD)" clean -i
	-rm -rf forecastle-*.tar

apply:
	kubectl apply -f deployments/manifests/ -n temp-forecastle

deploy: binary-image push apply

# Bump Chart
bump-chart:
	sed -i "s/^version:.*/version: v$(VERSION)/" deployments/kubernetes/chart/forecastle/Chart.yaml
	sed -i "s/^appVersion:.*/appVersion: v$(VERSION)/" deployments/kubernetes/chart/forecastle/Chart.yaml
	sed -i "s/tag:.*/tag: v$(VERSION)/" deployments/kubernetes/chart/forecastle/values.yaml