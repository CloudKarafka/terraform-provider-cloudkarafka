GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=cloudkarafka
PROVIDER_VERSION = 0.10.0

default: build

tools:
	GO111MODULE=on go install github.com/client9/misspell/cmd/misspell
	GO111MODULE=on go install github.com/golangci/golangci-lint/cmd/golangci-lint

build: fmtcheck
	go install -ldflags "-X 'main.version=$(PROVIDER_VERSION)'"

local-clean:
	rm -rf ~/.terraform.d/plugins/localhost/cloudkarafka/cloudkarafka/$(PROVIDER_VERSION)/darwin_amd64/*

local-build: local-clean
	@echo $(GOOS);
	@echo $(GOARCH);
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-X 'main.version=$(PROVIDER_VERSION)'" -o terraform-provider-cloudkarafka_v$(PROVIDER_VERSION)

local-install: local-build
	mkdir -p ~/.terraform.d/plugins/localhost/cloudkarafka/cloudkarafka/$(PROVIDER_VERSION)/darwin_amd64
	cp $(CURDIR)/terraform-provider-cloudkarafka_v$(PROVIDER_VERSION) ~/.terraform.d/plugins/localhost/cloudkarafka/cloudkarafka/$(PROVIDER_VERSION)/darwin_amd64

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

lint:
	@echo "==> Checking source code against linters..."
	golangci-lint run ./...

.PHONY: fmtcheck lint tools
