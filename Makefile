.PHONY: install test build binary-image push

BUILDER ?= forecastle-builder
DOCKER_IMAGE ?= stakater/forecastle

# Default value "dev"
DOCKER_TAG ?= dev
REPOSITORY = ${DOCKER_IMAGE}:${DOCKER_TAG}

GOCMD = go

install:
	go mod download

test:
	"$(GOCMD)" test -v ./...

build:

builder-image:
	@docker build --network host -t "${BUILDER}" -f build/package/Dockerfile.build .

binary-image: builder-image
	@docker run --network host --rm "${BUILDER}" | docker build --network host -t "${REPOSITORY}" -f Dockerfile.run -

push:
	docker push $(REPOSITORY)

# Bump Chart
bump-chart:
	sed -i "s/^version:.*/version:  $(VERSION)/" deployments/kubernetes/chart/forecastle/Chart.yaml
	sed -i "s/^appVersion:.*/appVersion:  $(VERSION)/" deployments/kubernetes/chart/forecastle/Chart.yaml
	sed -i "s/tag:.*/tag:  v$(VERSION)/" deployments/kubernetes/chart/forecastle/values.yaml