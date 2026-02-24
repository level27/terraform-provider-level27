default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

# Run unit tests
.PHONY: test
test:
	go test ./... -v $(TESTARGS) -timeout 30m

# Build the provider binary
.PHONY: build
build:
	go build -v ./...

# Install the provider locally for development
.PHONY: install
install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/level27/level27/$(VERSION)/linux_amd64
	cp terraform-provider-level27 ~/.terraform.d/plugins/registry.terraform.io/level27/level27/$(VERSION)/linux_amd64/

# Generate documentation from schema descriptions
.PHONY: docs
docs:
	go generate ./...

# Lint
.PHONY: lint
lint:
	golangci-lint run

# Format
.PHONY: fmt
fmt:
	go fmt ./...
	gofmt -s -w .

# Tidy modules
.PHONY: tidy
tidy:
	go mod tidy
