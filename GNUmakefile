GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=cloudkarafka
PROVIDER_VERSION = 0.10.0

default: build

## Check if a 64 bit kernel is running
UNAME_M := $(shell uname -m)

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
    GOOS += linux
endif
ifeq ($(UNAME_S),Darwin)
    GOOS += darwin
endif

UNAME_P := $(shell uname -p)
ifeq ($(UNAME_P),i386)
	ifeq ($(UNAME_M),x86_64)
		GOARCH += amd64
	else
		GOARCH += i386
	endif
else
    ifeq ($(UNAME_P),AMD64)
        GOARCH += amd64
    endif
endif
PROVIDER_ARCH = $(GOOS)_$(GOARCH)

tools:
	GO111MODULE=on go install github.com/client9/misspell/cmd/misspell
	GO111MODULE=on go install github.com/golangci/golangci-lint/cmd/golangci-lint

build: fmtcheck
	go install -ldflags "-X 'main.version=$(PROVIDER_VERSION)'"

local-clean:
	rm -rf ~/.terraform.d/plugins/localhost/cloudkarafka/cloudkarafka/$(PROVIDER_VERSION)/$(PROVIDER_ARCH)/terraform-provider-cloudkarafka_v$(PROVIDER_VERSION)

local-build: local-clean
	@echo $(GOOS);
	@echo $(GOARCH);
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-X 'main.version=$(PROVIDER_VERSION)'" -o terraform-provider-cloudkarafka_v$(PROVIDER_VERSION)

local-install: local-build
	mkdir -p ~/.terraform.d/plugins/localhost/cloudkarafka/cloudkarafka/$(PROVIDER_VERSION)/$(PROVIDER_ARCH)
	cp $(CURDIR)/terraform-provider-cloudkarafka_v$(PROVIDER_VERSION) ~/.terraform.d/plugins/localhost/cloudkarafka/cloudkarafka/$(PROVIDER_VERSION)/$(PROVIDER_ARCH)

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

lint:
	@echo "==> Checking source code against linters..."
	golangci-lint run ./...

.PHONY: fmtcheck lint tools
