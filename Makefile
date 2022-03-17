# Global variables ============================================================

SHELL := /bin/bash
SED := sed

# Go binary.
GO := go

# Project source code test coverage threshold.
COVERAGE := 30

# Usage components ============================================================

define USAGE_OPTIONS

Options:

   BINS        The binaries to build. Default is all commands in cmd/.
               This option is available for: make build/build.multiarch
               Example: make build BINS="apiserver otherbinary"
   IMAGES      Docker images to build. Default is all commands in cmd/.
               This option is available when using: make image/image.multiarch.
               Example: make image.multiarch IMAGES="apiserver otherbinary"
   PLATFORMS   The multiple platforms to build.
               Default is 'darwin_amd64 darwin_arm64 linux_amd64 linux_arm64'.
               This option is available when using: make build.multiarch.
               Example: make build.multiarch PLATFORMS="linux_amd64"
endef
export USAGE_OPTIONS

# Includes ====================================================================

include scripts/makefiles/share.makefile
include scripts/makefiles/go.makefile
include scripts/makefiles/image.makefile
include scripts/makefiles/generate.makefile
include scripts/makefiles/tools.makefile

# Targets =====================================================================

# Print help information by default.
.DEFAULT_GOAL := help

##  all: Make gen, lint, cover, build
.PHONY: all
all: gen lint cover build

##  run.dev: Run in development mode.
.PHONY: run.dev
run.dev:
	@./scripts/run.sh dev

##  run.test: Run in test mode.
.PHONY: run.test
run.test:
	@./scripts/run.sh test

##  build: Compile packages and dependencies to generate binary file for current platform.
.PHONY: build
build:
	@${MAKE} go.build

##  build.multiarch: Build for multiple platforms. See option PLATFORMS.
.PHONY: build.multiarch
build.multiarch:
	@${MAKE} go.build.multiarch

##  image: Build docker images for host arch.
.PHONY: image
image:
	@${MAKE} image.build

##  push: Build docker images for host arch and push images to registry.
.PHONY: push
push:
	@${MAKE} image.push

##  lint: Check syntax and style of Go source code.
.PHONY: lint
lint:
	@${MAKE} go.lint

##  test: Run unit test.
.PHONY: test
test:
	@${MAKE} go.test

##  cover: Run unit test and get test coverage.
.PHONY: cover
cover:
	@${MAKE} go.test.cover

##  gen: Generate necessary source code files and doc files.
.PHONY: gen
gen:
	@${MAKE} gen.run

##  clean: Remove all files that are created by building.
.PHONY: clean
clean:
	@echo "==========> Cleaning all build output"
	@-rm -vrf ${OUTPUT_DIR}

##  help: Show this help.
.PHONY: help
help: Makefile
	@echo -e "\nUsage: make [TARGETS] [OPTIONS] \n\nTargets:\n"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo "$$USAGE_OPTIONS"

# References:
# https://seisman.github.io/how-to-write-makefile/index.html
