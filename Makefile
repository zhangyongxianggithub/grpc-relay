HOMEDIR	:= $(shell pwd)
OUTDIR	:= $(HOMEDIR)/output

GO 		:= go
GOMOD 	:= $(GO) mod
GOBUILD := $(GO) build
GOTEST	:= $(GO) test -gcflags="-N -l"
GOPKGS  := $$(go list ./...| grep -vE "vendor")

include go.env
export GO111MODULE
export GOPRIVATE
export GOPROXY
export CGO_ENABLED

APP_VERSION	:= 1.0.0
GIT_BRANCH 	:= $(shell git rev-parse --abbrev-ref HEAD)
GIT_COMMIT 	:= $(shell git rev-parse HEAD)
BUILD_DATE	:= $(shell date "+%Y-%m-%d %H:%M:%S")
VERSION_VAR := bestzyx.com/grpc-relay/version
GOMODULE 	:= $(shell $(GO) list -m)

LD_FLAGS 	:= " -extldflags=-static \
	-X '$(VERSION_VAR).AppVersion=${APP_VERSION}' \
    -X '$(VERSION_VAR).GitBranch=${GIT_BRANCH}' \
	-X '$(VERSION_VAR).GitCommit=${GIT_COMMIT}' \
	-X '$(VERSION_VAR).BuildDate=${BUILD_DATE}' \
    -X '$(VERSION_VAR).BuildPipeline=${AGILE_PIPELINE_NAME}' \
	-X '$(VERSION_VAR).BuildNumber=${AGILE_PIPELINE_BUILD_NUMBER}' \
	-X '$(SWAGGER_VAR).DocPrefix=${SWAG_DOC_PREFIX}' \
    "

# make, make all
all: clean prepare compile package

prepare:
	git version     # 低于 2.17.1 可能不能正常工作
	go env          # 打印出 go 环境信息，可用于排查问题
	go mod download || go mod download -x  # 下载 依赖

#make compile
compile: build

build:
	cd gateway && $(GOMOD) tidy
	cd relay && $(GOMOD) tidy
	$(GOBUILD) -v -tags musl -ldflags=$(LD_FLAGS) -o $(HOMEDIR)/gateway/bin/gateway gateway/cmd/main.go
	$(GOBUILD) -v -tags musl -ldflags=$(LD_FLAGS) -o $(HOMEDIR)/relay/bin/relay relay/cmd/main.go


# make package
package:
	rm -rf $(OUTDIR)
	mkdir -p $(OUTDIR)
	mv $(HOMEDIR)/gateway/bin/gateway  $(OUTDIR)/
	mv $(HOMEDIR)/relay/bin/relay $(OUTDIR)/
	cp gateway/config/config.toml $(OUTDIR)/gateway.toml
	cp relay/config/config.toml $(OUTDIR)/relay.toml


# make clean
clean:
	go clean
	rm -rf $(OUTDIR)

# avoid filename conflict and speed up build
.PHONY: all prepare compile package clean build
