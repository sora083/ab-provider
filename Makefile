#export GO111MODULE=on

# all: restapp

# restapp: *.go
# 	go build -o ab-provider


NAME := ab-provider
VERSION := 1.0.0
REVISION := $(shell git rev-parse --short HEAD)

SRCS     := $(shell find . -type f -name '*.go')
LDFLAGS  := -ldflags="-s -w -extldflags \"-static\""
NOVENDOR := $(shell go list ./... | grep -v vendor)

DIST_DIRS := find * -type d -exec

.DEFAULT_GOAL := help

# ifndef GOBIN
# GOBIN := $(shell echo "$${GOPATH%%:*}/bin")
# endif

# LINT := $(GOBIN)/golint
# GOX := $(GOBIN)/gox
# ARCHIVER := $(GOBIN)/archiver
# GHR := $(GOBIN)/ghr

# $(LINT): ; @go get github.com/golang/lint/golint
# $(GOX): ; @go get github.com/mitchellh/gox
# $(ARCHIVER): ; @go get github.com/mholt/archiver/cmd/arc
# $(GHR): ; @go get github.com/tcnksm/ghr

.PHONY: help
help: ## Show help see: https://postd.cc/auto-documented-makefile/
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: deps
deps: ## Install dependency libraries
	go get -d -v .

.PHONY: tidy
tidy: ## remove unnecessary deps
	go mod tidy

.PHONY: update-deps
update-deps: ## Update dependency libraries
	go get -u

.PHONY: build
build: deps ## Build app for developers os
	CGO_ENABLED=0 go build $(LDFLAGS) -o dist/$(NAME)

.PHONY: server
server: ## Run api server for test
ifeq ($(shell command -v realize 2> /dev/null),)
	go get -u github.com/oxequa/realize
endif
	realize start

.PHONY: watch
watch: ## Watching modified go files and run unit test
ifeq ($(shell command -v ginkgo 2> /dev/null),)
	go get -u github.com/onsi/ginkgo/ginkgo
endif
	ginkgo watch ./...

# Alias
s: server
w: watch