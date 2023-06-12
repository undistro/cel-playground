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

.PHONY: update-data
update-data: ## Update the web/assets/data.json file
	yq -ojson '.' examples.yaml > web/assets/data.json
	yq -ojson -i '.versions.cel-go = "$(CEL_GO_VERSION)"' web/assets/data.json

.PHONY: addlicense
addlicense: ## Add copyright license headers in source code files.
	@test -s $(LOCALBIN)/addlicense || GOBIN=$(LOCALBIN) go install github.com/google/addlicense@latest
	$(LOCALBIN)/addlicense -c "Undistro Authors" -l "apache" -ignore ".github/**" -ignore ".idea/**" .

.PHONY: checklicense
checklicense: ## Check copyright license headers in source code files.
	@test -s $(LOCALBIN)/addlicense || GOBIN=$(LOCALBIN) go install github.com/google/addlicense@latest
	$(LOCALBIN)/addlicense -c "Undistro Authors" -l "apache" -ignore ".github/**" -ignore ".idea/**" -check .

##@ Build

.PHONY: build
build: fmt update-data ## Build WASM
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o web/main.wasm cmd/wasm/main.go
	gzip --best -f web/main.wasm

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)
