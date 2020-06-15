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
GENERATED      = client apis/zalando.org/v1/zz_generated.deepcopy.go informers listers

default: build.local

clean:
	rm -rf build
	rm -rf $(GENERATED)

test: $(GENERATED)
	go test -v $(GOPKGS)

check: $(GENERATED)
	staticcheck $(GOPKGS)
	go vet -v $(GOPKGS)

$(GENERATED):
	bash -x ./hack/update-codegen.sh

build.local: $(LOCAL_BINARIES)
build.linux: $(LINUX_BINARIES)

build/linux/%: $(SOURCES) $(GENERATED)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o build/linux/$(notdir $@) -ldflags "$(LDFLAGS)" ./cli/rg-client-test

build/%: $(SOURCES) $(GENERATED)
	CGO_ENABLED=0 go build -o build/$(notdir $@) $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" ./cli/rg-client-test
