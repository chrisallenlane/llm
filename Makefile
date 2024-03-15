# paths
makefile := $(realpath $(lastword $(MAKEFILE_LIST)))
cmd_dir  := ./cmd/llm
dist_dir := ./dist

# executables
CAT    := cat
COLUMN := column
CTAGS  := ctags
DOCKER := docker
GO     := go
GREP   := grep
GZIP   := gzip --best
LINT   := revive
MKDIR  := mkdir -p
RM     := rm
SCC    := scc
SED    := sed
SORT   := sort
ZIP    := zip -m

docker_image := llm-devel:latest

# build flags
BUILD_FLAGS  := -ldflags="-s -w" -mod vendor -trimpath
GOBIN        :=
TMPDIR       := /tmp

# release binaries
releases :=                        \
	$(dist_dir)/llm-darwin-amd64 \
	$(dist_dir)/llm-linux-386    \
	$(dist_dir)/llm-linux-amd64  \
	$(dist_dir)/llm-linux-arm5   \
	$(dist_dir)/llm-linux-arm6   \
	$(dist_dir)/llm-linux-arm64  \
	$(dist_dir)/llm-linux-arm7   \
	$(dist_dir)/llm-netbsd-amd64  \
	$(dist_dir)/llm-openbsd-amd64  \
	$(dist_dir)/llm-plan9-amd64  \
	$(dist_dir)/llm-solaris-amd64  \
	$(dist_dir)/llm-windows-amd64.exe

## build: build an executable for your architecture
.PHONY: build
build: | clean $(dist_dir) fmt lint vet vendor
	$(GO) build $(BUILD_FLAGS) -o $(dist_dir)/llm $(cmd_dir)

## build-release: build release executables
.PHONY: build-release
build-release: $(releases)

# llm-darwin-amd64
$(dist_dir)/llm-darwin-amd64: prepare
	GOARCH=amd64 GOOS=darwin \
	$(GO) build $(BUILD_FLAGS) -o $@ $(cmd_dir) && $(GZIP) $@ && chmod -x $@.gz

# llm-linux-386
$(dist_dir)/llm-linux-386: prepare
	GOARCH=386 GOOS=linux \
	$(GO) build $(BUILD_FLAGS) -o $@ $(cmd_dir) && $(GZIP) $@ && chmod -x $@.gz

# llm-linux-amd64
$(dist_dir)/llm-linux-amd64: prepare
	GOARCH=amd64 GOOS=linux \
	$(GO) build $(BUILD_FLAGS) -o $@ $(cmd_dir) && $(GZIP) $@ && chmod -x $@.gz

# llm-linux-arm5
$(dist_dir)/llm-linux-arm5: prepare
	GOARCH=arm GOOS=linux GOARM=5 \
	$(GO) build $(BUILD_FLAGS) -o $@ $(cmd_dir) && $(GZIP) $@ && chmod -x $@.gz

# llm-linux-arm6
$(dist_dir)/llm-linux-arm6: prepare
	GOARCH=arm GOOS=linux GOARM=6 \
	$(GO) build $(BUILD_FLAGS) -o $@ $(cmd_dir) && $(GZIP) $@ && chmod -x $@.gz

# llm-linux-arm7
$(dist_dir)/llm-linux-arm7: prepare
	GOARCH=arm GOOS=linux GOARM=7 \
	$(GO) build $(BUILD_FLAGS) -o $@ $(cmd_dir) && $(GZIP) $@ && chmod -x $@.gz
	
# llm-linux-arm64
$(dist_dir)/llm-linux-arm64: prepare
	GOARCH=arm64 GOOS=linux \
	$(GO) build $(BUILD_FLAGS) -o $@ $(cmd_dir) && $(GZIP) $@ && chmod -x $@.gz

# llm-netbsd-amd64
$(dist_dir)/llm-netbsd-amd64: prepare
	GOARCH=amd64 GOOS=netbsd \
	$(GO) build $(BUILD_FLAGS) -o $@ $(cmd_dir) && $(GZIP) $@ && chmod -x $@.gz

# llm-openbsd-amd64
$(dist_dir)/llm-openbsd-amd64: prepare
	GOARCH=amd64 GOOS=openbsd \
	$(GO) build $(BUILD_FLAGS) -o $@ $(cmd_dir) && $(GZIP) $@ && chmod -x $@.gz

# llm-plan9-amd64
$(dist_dir)/llm-plan9-amd64: prepare
	GOARCH=amd64 GOOS=plan9 \
	$(GO) build $(BUILD_FLAGS) -o $@ $(cmd_dir) && $(GZIP) $@ && chmod -x $@.gz

# llm-solaris-amd64
$(dist_dir)/llm-solaris-amd64: prepare
	GOARCH=amd64 GOOS=solaris \
	$(GO) build $(BUILD_FLAGS) -o $@ $(cmd_dir) && $(GZIP) $@ && chmod -x $@.gz

# llm-windows-amd64
$(dist_dir)/llm-windows-amd64.exe: prepare
	GOARCH=amd64 GOOS=windows \
	$(GO) build $(BUILD_FLAGS) -o $@ $(cmd_dir) && $(ZIP) $@.zip $@ -j

# ./dist
$(dist_dir):
	$(MKDIR) $(dist_dir)

## install: build and install llm on your PATH
.PHONY: install
install: build
	$(GO) install $(BUILD_FLAGS) $(GOBIN) $(cmd_dir) 

## clean: remove compiled executables
.PHONY: clean
clean:
	$(RM) -f $(dist_dir)/*

## distclean: remove the tags file
.PHONY: distclean
distclean:
	$(RM) -f tags
	@$(DOCKER) image rm -f $(docker_image)

## setup: install revive (linter) and scc (sloc tool)
.PHONY: setup
setup:
	GO111MODULE=off $(GO) get -u github.com/boyter/scc github.com/mgechev/revive

## sloc: count "semantic lines of code"
.PHONY: sloc
sloc:
	$(SCC) --exclude-dir=vendor

## tags: build a tags file
.PHONY: tags
tags:
	$(CTAGS) -R --exclude=vendor --languages=go

## vendor: download, tidy, and verify dependencies
.PHONY: vendor
vendor:
	$(GO) mod vendor && $(GO) mod tidy && $(GO) mod verify

## vendor-update: update vendored dependencies
vendor-update:
	$(GO) get -t -u ./... && $(GO) mod vendor && $(GO) mod tidy && $(GO) mod verify

## fmt: run go fmt
.PHONY: fmt
fmt:
	$(GO) fmt ./...

## lint: lint go source files
.PHONY: lint
lint: vendor
	$(LINT) -exclude vendor/... ./...

## vet: vet go source files
.PHONY: vet
vet:
	$(GO) vet ./...

## test: run unit-tests
.PHONY: test
test: 
	$(GO) test ./...

## coverage: generate a test coverage report
.PHONY: coverage
coverage:
	$(GO) test ./... -coverprofile=$(TMPDIR)/llm-coverage.out && \
	$(GO) tool cover -html=$(TMPDIR)/llm-coverage.out

## check: format, lint, vet, vendor, and run unit-tests
.PHONY: check
check: | vendor fmt lint vet test

.PHONY: prepare
prepare: | clean $(dist_dir) vendor fmt lint vet test

## docker-setup: create a docker image for use during development
.PHONY: docker-setup
docker-setup:
	$(DOCKER) build  -t $(docker_image) -f Dockerfile .

## docker-sh: shell into the docker development container
.PHONY: docker-sh
docker-sh:
	$(DOCKER) run -v $(shell pwd):/app -ti $(docker_image) /bin/ash

## help: display this help text
.PHONY: help
help:
	@$(CAT) $(makefile) | \
	$(SORT)             | \
	$(GREP) "^##"       | \
	$(SED) 's/## //g'   | \
	$(COLUMN) -t -s ':'
