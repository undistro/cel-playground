CEL_GO_VERSION=$(shell go list -f "{{.Version}}" -m github.com/google/cel-go)

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: test
test: fmt ## Run tests.
	go test ./... -coverprofile cover.out

.PHONY: serve
serve: ## Serve static files
	python3 -m http.server -d web/ 8080
	#Uncomment the command below to serve with Go
	#go run cmd/server/main.go --dir web/

.PHONY: dotenv
dotenv: ## Update web/.env file with cel-go dependency version
	@if grep -q '^CEL_GO_VERSION=' web/.env; then \
  		sed -i 's/^\(CEL_GO_VERSION=\).*/\1$(CEL_GO_VERSION)/' web/.env; \
  		else echo "CEL_GO_VERSION=$(CEL_GO_VERSION)" >> web/.env; \
  	fi
	@cat web/.env

##@ Build

.PHONY: build
build: fmt dotenv ## Build WASM
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o web/main.wasm cmd/wasm/main.go
