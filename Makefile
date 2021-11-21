.PHONY: all
all: test install ## build APIKit and run tests

ifeq ($(OS), Windows_NT)
SHELL := powershell.exe
.SHELLFLAGS := -NoProfile -Command
endif

PACKAGE_ROOT = github.com/ExperienceOne/apikit
PACKAGE_VERSION = ${PACKAGE_ROOT}/internal/framework/version

ifeq ($(OS), Windows_NT)
BUILD_TIME=$(shell Get-Date -Format "yyyy-MM-ddTHH:mm:ss")
else
BUILD_TIME=$(shell date)
endif

GIT_COMMIT=""
GIT_BRANCH=""
GIT_TAG=""

ifneq ($(wildcard .git),)
    $(info use git meta data)
	GIT_COMMIT=$(shell git rev-parse HEAD)
	GIT_TAG=$(shell git describe --abbrev=0 --tags)
	ifeq ($(OS), Windows_NT)
		GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
	else
		GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)
	endif
else
    $(error no git meta data inside of this dir)
endif

BUILD_INFO_FLAGS = -X '${PACKAGE_VERSION}.GitCommit=${GIT_COMMIT}' -X '${PACKAGE_VERSION}.GitBranch=${GIT_BRANCH}' -X '${PACKAGE_VERSION}.GitTag=${GIT_TAG}' -X '${PACKAGE_VERSION}.BuildTime=${BUILD_TIME}'

ALL_PACKAGES=$(shell go list ./...)
GOPATH=$(shell go env GOPATH)

.PHONY: framework
framework: ## build framework components
ifeq ($(OS), Windows_NT)
	go build -ldflags "${BUILD_INFO_FLAGS}" -o fpacker.exe .\cmd\fpacker\main.go
	move -force .\fpacker.exe $(GOPATH)\bin\fpacker.exe
	$(GOPATH)\bin\fpacker -src '.\internal\framework\' -dest '.\framework\framework_code.go'
	$(GOPATH)\bin\fpacker -src '.\internal\framework\' -dest '.\framework\framework_code_client.go' -exclude='xserver,validation,middleware,unmarshal' -kind=client
	$(GOPATH)\bin\fpacker -src '.\internal\framework\' -dest '.\framework\framework_code_server.go' -exclude='xclient,roundtripper,hooks' -kind=server
else
	go build -ldflags "${BUILD_INFO_FLAGS}" -o fpacker ./cmd/fpacker
	mv ./fpacker $(GOPATH)/bin/fpacker
	$(GOPATH)/bin/fpacker -src ./internal/framework/ -dest ./framework/framework_code.go
	$(GOPATH)/bin/fpacker -src ./internal/framework/ -dest ./framework/framework_code_client.go -exclude=xserver,validation,middleware,unmarshal -kind=client
	$(GOPATH)/bin/fpacker -src ./internal/framework/ -dest ./framework/framework_code_server.go -exclude=xclient,roundtripper,hooks -kind=server
endif

testgenerator: framework
ifeq ($(OS), Windows_NT)
	go build -ldflags "${BUILD_INFO_FLAGS}" -o apikit.exe  .\cmd\apikit\main.go
	move -force .\apikit.exe $(GOPATH)\bin\test_apikit.exe
	git checkout -- .\framework\framework_code.go
	git checkout -- .\framework\framework_code_client.go
	git checkout -- .\framework\framework_code_server.go
else
	go build -ldflags "${BUILD_INFO_FLAGS}" -o apikit  ./cmd/apikit/main.go
	mv ./apikit $(GOPATH)/bin/test_apikit
	git checkout @ -- ./framework/framework_code.go
	git checkout @ -- ./framework/framework_code_client.go
	git checkout @ -- ./framework/framework_code_server.go
endif

.PHONY: test
test: testgenerator ## run tests
ifeq ($(OS), Windows_NT)
	$(GOPATH)\bin\test_apikit --debug  generate  .\tests\data\swagger.yaml  .\tests\api\ api --mocked
	$(GOPATH)\bin\test_apikit --debug  generate .\example\api.yaml  .\example todo --mocked
else
	$(GOPATH)/bin/test_apikit --debug  generate  ./tests/data/swagger.yaml  ./tests/api/ api --mocked
	$(GOPATH)/bin/test_apikit --debug  generate ./example/api.yaml  ./example todo --mocked
endif
	go test -v -failfast ./...

.PHONY: lint
lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.40.0 golangci-lint run -v --max-issues-per-linter=0 --max-same-issues=0

.PHONY: install
install: framework   ## builds and installs the binaries of the APIKit in $GOPATH/bin
ifeq ($(OS), Windows_NT)
	go build -ldflags "${BUILD_INFO_FLAGS}" -o apikit.exe  .\cmd\apikit\main.go
	move -force apikit.exe $(GOPATH)\bin\apikit.exe
else
	go build -ldflags "${BUILD_INFO_FLAGS}" -o apikit  ./cmd/apikit/main.go
	mv apikit $(GOPATH)/bin/apikit
endif

.PHONY: help
help:	            ## prints help for most of the make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
