.PHONY: all
all: test install ## build APIKit and run tests

ALL_PACKAGES=$(shell go list ./... | grep -v "vendor")
GOPATH = $(shell printenv GOPATH)
ifeq ($(GOPATH), )
	GOPATH = ~/go
endif

.PHONY: framework
framework: ## build framework components
	scripts/version.sh
	go install -v ./cmd/fpacker
	$(GOPATH)/bin/fpacker -src ./internal/framework/ -dest ./framework/framework_code.go
	$(GOPATH)/bin/fpacker -src ./internal/framework/ -dest ./framework/framework_code_client.go -exclude=xserver,validation,middleware,unmarshal -kind=client
	$(GOPATH)/bin/fpacker -src ./internal/framework/ -dest ./framework/framework_code_server.go -exclude=xclient,roundtripper,hooks -kind=server

testgenerator: framework
	go build ./cmd/apikit
	mv ./apikit $(GOPATH)/bin/test_apikit
	git checkout @ -- ./internal/framework/version/version.go
	git checkout @ -- ./framework/framework_code.go
	git checkout @ -- ./framework/framework_code_client.go
	git checkout @ -- ./framework/framework_code_server.go

.PHONY: test
test: testgenerator ## run tests
	$(GOPATH)/bin/test_apikit --debug generate ./tests/data/swagger.yaml  ./tests/api/ api
	$(GOPATH)/bin/test_apikit --debug generate --mocked ./tests/data/swagger.yaml  ./tests/mock/ api
	for package in $(ALL_PACKAGES); do ENVIRONMENT=test go test $$package; if [ $$? -ne "0" ]; then echo "Test failed!"; exit 1; fi; done

.PHONY: lint
lint:				## runs linters for all packages except 'gen' (requires golangci-lint)
	golangci-lint run --skip-dirs '(example|tests)'

.PHONY: install
install: framework   ## builds and installs the binaries of the APIKit in $GOPATH/bin
	scripts/version.sh
	cd ./cmd/apikit && \
	go build -o apikit
	mv ./cmd/apikit/apikit $(GOPATH)/bin/apikit
	git checkout @ -- ./internal/framework/version/version.go

.PHONY: help
help:	            ## prints help for most of the make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
