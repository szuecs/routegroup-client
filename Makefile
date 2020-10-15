.PHONY: clean test check generate build.local build.linux build.osx

BINARY         ?= rg-example
BINARIES       = $(BINARY)
LOCAL_BINARIES = $(addprefix build/,$(BINARIES))
LINUX_BINARIES = $(addprefix build/linux/,$(BINARIES))
VERSION        ?= $(shell git describe --tags --always --dirty)
SOURCES        = $(shell find . -name '*.go')
GOPKGS         = $(shell go list ./...)
BUILD_FLAGS    ?= -v
LDFLAGS        ?= -X main.version=$(VERSION) -w -s
CRD_TYPE_SOURCE = apis/zalando.org/v1/types.go
GENERATED      = apis/zalando.org/v1/zz_generated.deepcopy.go $(shell find client/ -name '*.go')
GENERATED_CRD  = zalando.org_routegroups.yaml

default: build.local

clean:
	rm -rf build
	rm -rf $(GENERATED)
	rm -rf $(GENERATED_CRD)

test: $(GENERATED)
	go test -v $(GOPKGS)

check: $(GENERATED)
	staticcheck $(GOPKGS)
	go vet -v $(GOPKGS)

$(GENERATED): $(CRD_TYPE_SOURCE)
	bash -x ./hack/update-codegen.sh

$(GENERATED_CRD): go.mod $(GENERATED)
	go run sigs.k8s.io/controller-tools/cmd/controller-gen crd:crdVersions=v1 paths=./apis/... output:crd:dir=.
	# workaround to add pattern to array items. Not supported by controller-gen
	# ref: https://github.com/kubernetes-sigs/controller-tools/issues/342
	perl -i -p0e 's|(\s*)(hosts:.*?items:)|$$1$$2$$1    pattern: "^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*\$$"|sg' $(GENERATED_CRD)

build.local: $(LOCAL_BINARIES) $(GENERATED_CRD)
build.linux: $(LINUX_BINARIES) $(GENERATED_CRD)

build/linux/%: $(SOURCES) $(GENERATED_CRD)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o build/linux/$(notdir $@) -ldflags "$(LDFLAGS)" ./cli/rg-client-test

build/%: $(SOURCES) $(GENERATED_CRD)
	CGO_ENABLED=0 go build -o build/$(notdir $@) $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" ./cli/rg-client-test
