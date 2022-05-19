UNAME := $(shell uname)
GO_PATH:=$(shell go env GOPATH)
BUILD_PATH=$(CURDIR)/.aws-sam/build
FUNCTION_NAME=HelloWorldFunction
LINTER_ARGS = run -c .golangci.yml --timeout 5m

.PHONY: help
help:	## Show a list of available commands
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

.PHONY: install
install:	## Download dependencies
	go mod download

.PHONY: build
build:	## Build from sam template
	sam build

.PHONY: validate-template
validate-template: ## Validate sam template
	sam validate

.PHONY: zip-deployment
zip-deployment:	## zip sam build deployment
	cd $(BUILD_PATH)/$(FUNCTION_NAME) && zip -r $(FUNCTION_NAME).zip . && mv $(FUNCTION_NAME).zip ..

.PHONY: clean-zip-deployment
clean-zip-deployment:	## remove zip sam build deployment
	cd $(BUILD_PATH)/ && rm -f $(FUNCTION_NAME).zip

.PHONY: test
test:	## Run tests
	go test -p 1 -cover -v ./tests/acceptance -timeout 5m

.PHONY: lint
lint:	## Run static linting of source files. See .golangci.yml for options
	golangci-lint $(LINTER_ARGS)

.PHONY: download-tools
download-tools:	## Download all required tools to generate documentation, code analysis...
	@echo "Installing tools on $(GO_PATH)/bin"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.0
	go install golang.org/x/tools/cmd/goimports@v0.1.9
	@echo "Tools installed"	