SHELL = /bin/bash

GO = godep go
GOFMT = gofmt -l
GOLINT = golint
GOTEST = $(GO) test --cover --race -v
GOVET = $(GO) vet

SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
DOCKER_HOST = $(shell echo $$DOCKER_HOST)
REGISTRY = $(shell echo $$DEV_REGISTRY)
GIT_SHA = $(shell git rev-parse --short HEAD)
IMAGE_PREFIX := deis/

ifndef BUILD_TAG
  BUILD_TAG = git-$(GIT_SHA)
endif

ifndef S3_BUCKET
  S3_BUCKET = deis-updates
endif

ifndef DEIS_NUM_INSTANCES
  DEIS_NUM_INSTANCES = 3
endif

define echo_cyan
  @echo "\033[0;36m$(subst ",,$(1))\033[0m"
endef

define echo_yellow
  @echo "\033[0;33m$(subst ",,$(1))\033[0m"
endef

check-docker:
	@if [ -z $$(which docker) ]; then \
	  echo "Missing \`docker\` client which is required for development"; \
	  exit 2; \
	fi

check-registry:
	@if [ -z "$$DEV_REGISTRY" ]; then \
	  echo "DEV_REGISTRY is not exported, try:  make dev-registry"; \
	exit 2; \
	fi

check-deisctl:
	@if [ -z $$(which deisctl) ]; then \
	  echo "Missing \`deisctl\` utility, please install from https://github.com/deis/deis"; \
	fi
