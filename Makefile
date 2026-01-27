.PHONY: default
default: build lint test

.PHONY: build
build:
	go build ./...

GOLANG_TOOL_PATH_TO_BIN=$(shell go env GOPATH)
GOLANGCI_LINT_CLI_VERSION?=v2.8.0
GOLANGCI_LINT_CLI=$(GOLANG_TOOL_PATH_TO_BIN)/bin/golangci-lint
$(GOLANGCI_LINT_CLI):
	$(MAKE) golangci-lint-cli-install

# binary will be $(go env GOPATH)/bin/golangci-lint
golangci-lint-cli-install:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | \
	sh -s -- -b $(GOLANG_TOOL_PATH_TO_BIN)/bin $(GOLANGCI_LINT_CLI_VERSION)

.PHONY: lint
lint: $(GOLANGCI_LINT_CLI)
	$(GOLANGCI_LINT_CLI) run

GO_TEST_OPTIONS?=

.PHONY: test
test:
	go test $(GO_TEST_OPTIONS) ./...
