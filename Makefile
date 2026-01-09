.PHONY: default build build-frontend test test-e2e test-e2e-frontend stop clean-images clean push apply deploy release release-all manifest push clean-image docker-build

OS ?= linux
ARCH ?= amd64
ALL_ARCH ?= arm64 arm amd64

BINARY ?= forecastle
DOCKER_IMAGE ?= stakater/forecastle

# Default value "dev"
DOCKER_TAG ?= dev
REPOSITORY_GENERIC = ${DOCKER_IMAGE}:${DOCKER_TAG}
REPOSITORY_ARCH = ${DOCKER_IMAGE}:${DOCKER_TAG}-${ARCH}

VERSION ?= 0.0.1
BUILD=

GOCMD = go
GOFLAGS ?= $(GOFLAGS:)
LDFLAGS = -s -w

default: build test

install:
	"$(GOCMD)" mod download

build-frontend:
	cd frontend && yarn install && yarn build && cd ..
	rm -rf internal/web/build
	mkdir -p internal/web/build
	cp -r frontend/build/* internal/web/build/

build: build-frontend
	"$(GOCMD)" build ${GOFLAGS} -ldflags="${LDFLAGS}" -o "${BINARY}" ./cmd/forecastle

build-go:
	"$(GOCMD)" build ${GOFLAGS} -ldflags="${LDFLAGS}" -o "${BINARY}" ./cmd/forecastle

docker-build:
	docker buildx build --platform ${OS}/${ARCH} -t "${REPOSITORY_ARCH}" --load -f Dockerfile .

push:
	docker push ${REPOSITORY_ARCH}

release: docker-build push manifest

release-all:
	-rm -rf ~/.docker/manifests/*
	@for arch in $(ALL_ARCH) ; do \
		echo Make release: $$arch ; \
		make release ARCH=$$arch ; \
	done
	set -e
	docker manifest push --purge $(REPOSITORY_GENERIC)

manifest:
	set -e
	docker manifest create -a $(REPOSITORY_GENERIC) $(REPOSITORY_ARCH)
	docker manifest annotate --arch $(ARCH) $(REPOSITORY_GENERIC) $(REPOSITORY_ARCH)

test:
	"$(GOCMD)" test -timeout 1800s -v ./cmd/... ./internal/... ./pkg/...

# Run e2e tests against a running Forecastle instance
# Requires: kubectl access to a cluster, FORECASTLE_URL env var (default: http://localhost:3000)
test-e2e:
	"$(GOCMD)" test -timeout 600s -v ./e2e/...

# Run frontend e2e tests with Playwright
# Requires: cd frontend && yarn install first
test-e2e-frontend:
	cd frontend && yarn test:e2e

# Run frontend e2e tests in headed mode (for debugging)
test-e2e-frontend-headed:
	cd frontend && yarn test:e2e:headed

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
# Chart version has no 'v' prefix (helm semver standard)
# appVersion and image tag keep 'v' prefix (matches Docker tag)
bump-chart:
	sed -i "s/^version:.*/version: $(VERSION)/" deployments/kubernetes/chart/forecastle/Chart.yaml
	sed -i "s/^appVersion:.*/appVersion: \"v$(VERSION)\"/" deployments/kubernetes/chart/forecastle/Chart.yaml
	sed -i "s/^    version:.*/    version: $(VERSION)/" deployments/kubernetes/chart/forecastle/values.yaml
	sed -i "s/^    tag:.*/    tag: v$(VERSION)/" deployments/kubernetes/chart/forecastle/values.yaml