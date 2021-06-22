.PHONY: all
all: test install ## build APIKit and run tests

PACKAGE_ROOT = github.com/ExperienceOne/apikit
PACKAGE_VERSION = ${PACKAGE_ROOT}/internal/framework/version

GIT_COMMIT=""
GIT_BRANCH=""
GIT_TAG=""
BUILD_TIME=$(shell date)

ifneq ($(wildcard .git),)
    $(info use git meta data)
	GIT_COMMIT=$(shell git rev-parse HEAD)
	GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)
	GIT_TAG=$(shell git describe --abbrev=0 --tags)
	BUILD_TIME=$(shell date)
else
    $(error no git meta data inside of this dir)
endif

BUILD_INFO_FLAGS = -X '${PACKAGE_VERSION}.GitCommit=${GIT_COMMIT}' -X '${PACKAGE_VERSION}.GitBranch=${GIT_BRANCH}' -X '${PACKAGE_VERSION}.GitTag=${GIT_TAG}' -X '${PACKAGE_VERSION}.BuildTime=${BUILD_TIME}'

ALL_PACKAGES=$(shell go list ./... | grep -v "vendor")
GOPATH = $(shell printenv GOPATH)
ifeq ($(GOPATH), )
	GOPATH = ~/go
endif

.PHONY: framework
framework: ## build framework components
	go install -ldflags "${BUILD_INFO_FLAGS}" -v ./cmd/fpacker
	$(GOPATH)/bin/fpacker -src ./internal/framework/ -dest ./framework/framework_code.go
	$(GOPATH)/bin/fpacker -src ./internal/framework/ -dest ./framework/framework_code_client.go -exclude=xserver,validation,middleware,unmarshal -kind=client
	$(GOPATH)/bin/fpacker -src ./internal/framework/ -dest ./framework/framework_code_server.go -exclude=xclient,roundtripper,hooks -kind=server

testgenerator: framework
	go build -ldflags "${BUILD_INFO_FLAGS}" -o apikit  ./cmd/apikit/main.go
	mv ./apikit $(GOPATH)/bin/test_apikit
	git checkout @ -- ./framework/framework_code.go
	git checkout @ -- ./framework/framework_code_client.go
	git checkout @ -- ./framework/framework_code_server.go

.PHONY: test
test: testgenerator ## run tests
	$(GOPATH)/bin/test_apikit --debug generate ./tests/data/swagger.yaml  ./tests/api/ api
	$(GOPATH)/bin/test_apikit --debug generate --mocked ./tests/data/swagger.yaml  ./tests/mock/ api
	$(GOPATH)/bin/test_apikit --debug generate ./example/api.yaml  ./example todo
	$(GOPATH)/bin/test_apikit --debug generate ./example/api.yaml  ./example todo
	for package in $(ALL_PACKAGES); do ENVIRONMENT=test go test -count=1 -v $$package; if [ $$? -ne "0" ]; then echo "Test failed!"; exit 1; fi; done

.PHONY: lint
lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.40.0 golangci-lint run -v --max-issues-per-linter=0 --max-same-issues=0

.PHONY: install
install: framework   ## builds and installs the binaries of the APIKit in $GOPATH/bin
	go build -ldflags "${BUILD_INFO_FLAGS}" -o apikit  ./cmd/apikit/main.go
	mv apikit $(GOPATH)/bin/apikit

.PHONY: help
help:	            ## prints help for most of the make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
