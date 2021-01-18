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
	

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.3.0 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif



bump-chart-operator:
	sed -i "s/^version:.*/version:  $(VERSION)/" deployments/kubernetes/chart/forecastle/Chart.yaml
	sed -i "s/^appVersion:.*/appVersion:  $(VERSION)/" deployments/kubernetes/chart/forecastle/Chart.yaml
	sed -i "s/tag:.*/tag:  v$(VERSION)/" deployments/kubernetes/chart/forecastle/values.yaml

# Bump Chart
bump-chart: bump-chart-operator

generate-crds: controller-gen
	$(CONTROLLER_GEN) crd paths="./..." output:crd:artifacts:config=deployments/kubernetes/chart/forecastle/crds