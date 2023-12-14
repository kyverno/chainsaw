.DEFAULT_GOAL := build

#########
# TOOLS #
#########

TOOLS_DIR                          := $(PWD)/.tools
CONTROLLER_GEN                     := $(TOOLS_DIR)/controller-gen
CONTROLLER_GEN_VERSION             := v0.12.0
REGISTER_GEN                       := $(TOOLS_DIR)/register-gen
DEEPCOPY_GEN                       := $(TOOLS_DIR)/deepcopy-gen
CODE_GEN_VERSION                   := v0.28.0
REFERENCE_DOCS                     := $(TOOLS_DIR)/genref
REFERENCE_DOCS_VERSION             := latest
KIND                               := $(TOOLS_DIR)/kind
KIND_VERSION                       := v0.20.0
TOOLS                              := $(CONTROLLER_GEN) $(REGISTER_GEN) $(DEEPCOPY_GEN) $(REFERENCE_DOCS) $(KIND)
PIP                                ?= "pip"
ifeq ($(GOOS), darwin)
SED                                := gsed
else
SED                                := sed
endif
COMMA                              := ,

$(CONTROLLER_GEN):
	@echo Install controller-gen... >&2
	@GOBIN=$(TOOLS_DIR) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_GEN_VERSION)

$(REGISTER_GEN):
	@echo Install register-gen... >&2
	@GOBIN=$(TOOLS_DIR) go install k8s.io/code-generator/cmd/register-gen@$(CODE_GEN_VERSION)

$(DEEPCOPY_GEN):
	@echo Install deepcopy-gen... >&2
	@GOBIN=$(TOOLS_DIR) go install k8s.io/code-generator/cmd/deepcopy-gen@$(CODE_GEN_VERSION)

$(REFERENCE_DOCS):
	@echo Install genref... >&2
	@GOBIN=$(TOOLS_DIR) go install github.com/kubernetes-sigs/reference-docs/genref@$(REFERENCE_DOCS_VERSION)

$(KIND):
	@echo Install kind... >&2
	@GOBIN=$(TOOLS_DIR) go install sigs.k8s.io/kind@$(KIND_VERSION)

.PHONY: install-tools
install-tools: $(TOOLS) ## Install tools

.PHONY: clean-tools
clean-tools: ## Remove installed tools
	@echo Clean tools... >&2
	@rm -rf $(TOOLS_DIR)

###########
# CODEGEN #
###########

ORG                         ?= kyverno
PACKAGE                     ?= github.com/$(ORG)/chainsaw
GOPATH_SHIM                 := ${PWD}/.gopath
PACKAGE_SHIM                := $(GOPATH_SHIM)/src/$(PACKAGE)
INPUT_DIRS                  := $(PACKAGE)/pkg/apis/v1alpha1
CRDS_PATH                   := ${PWD}/.crds

$(GOPATH_SHIM):
	@echo Create gopath shim... >&2
	@mkdir -p $(GOPATH_SHIM)

.INTERMEDIATE: $(PACKAGE_SHIM)
$(PACKAGE_SHIM): $(GOPATH_SHIM)
	@echo Create package shim... >&2
	@mkdir -p $(GOPATH_SHIM)/src/github.com/$(ORG) && ln -s -f ${PWD} $(PACKAGE_SHIM)

.PHONY: codegen-register
codegen-register: $(PACKAGE_SHIM) $(REGISTER_GEN) ## Generate types registrations
	@echo Generate registration... >&2
	@GOPATH=$(GOPATH_SHIM) $(REGISTER_GEN) \
		--go-header-file=./.hack/boilerplate.go.txt \
		--input-dirs=$(INPUT_DIRS)

.PHONY: codegen-deepcopy
codegen-deepcopy: $(PACKAGE_SHIM) $(DEEPCOPY_GEN) ## Generate deep copy functions
	@echo Generate deep copy functions... >&2
	@GOPATH=$(GOPATH_SHIM) $(DEEPCOPY_GEN) \
		--go-header-file=./.hack/boilerplate.go.txt \
		--input-dirs=$(INPUT_DIRS) \
		--output-file-base=zz_generated.deepcopy

.PHONY: codegen-crds
codegen-crds: $(CONTROLLER_GEN) ## Generate CRDs
	@echo Generate crds... >&2
	@rm -rf $(CRDS_PATH)
	@$(CONTROLLER_GEN) crd paths=./pkg/apis/... crd:crdVersions=v1 output:dir=$(CRDS_PATH)
	@echo Copy generated CRDs to embed in the CLI... >&2
	@rm -rf pkg/data/crds && mkdir -p pkg/data/crds
	@cp $(CRDS_PATH)/* pkg/data/crds

.PHONY: codegen-cli-docs
codegen-cli-docs: build ## Generate CLI docs
	@echo Generate cli docs... >&2
	@rm -rf website/docs/commands && mkdir -p website/docs/commands
	@rm -rf docs/user/commands && mkdir -p docs/user/commands
	@./$(CLI_BIN) docs -o website/docs/commands --autogenTag=false

.PHONY: codegen-api-docs
codegen-api-docs: $(REFERENCE_DOCS) ## Generate markdown API docs
	@echo Generate api docs... >&2
	@rm -rf ./website/docs/apis
	@cd ./website/apis && $(REFERENCE_DOCS) -c config.yaml -f markdown -o ../docs/apis

.PHONY: codegen-jp-docs
codegen-jp-docs: ## Generate JP docs
	@echo Generate jp docs... >&2
	@rm -rf ./website/docs/jp && mkdir -p ./website/docs/jp
	@go run ./website/jp/main.go > ./website/docs/jp/functions.md

.PHONY: codegen-mkdocs
codegen-mkdocs: codegen-cli-docs codegen-api-docs codegen-jp-docs ## Generate mkdocs website
	@echo Generate mkdocs website... >&2
	@$(PIP) install mkdocs
	@$(PIP) install --upgrade pip
	@$(PIP) install -U mkdocs-material mkdocs-redirects mkdocs-minify-plugin mkdocs-include-markdown-plugin lunr mkdocs-rss-plugin mike
	@mkdocs build -f ./website/mkdocs.yaml

.PHONY: codegen-schemas-openapi
codegen-schemas-openapi: CURRENT_CONTEXT = $(shell kubectl config current-context)
codegen-schemas-openapi: codegen-crds $(KIND) ## Generate openapi schemas (v2 and v3)
	@echo Generate openapi schema... >&2
	@rm -rf ./.temp/.schemas
	@mkdir -p ./.temp/.schemas/openapi/v2
	@mkdir -p ./.temp/.schemas/openapi/v3/apis/chainsaw.kyverno.io
	@$(KIND) create cluster --name schema --image $(KIND_IMAGE)
	@kubectl create -f $(CRDS_PATH)
	@sleep 15
	@kubectl get --raw /openapi/v2 > ./.temp/.schemas/openapi/v2/schema.json
	@kubectl get --raw /openapi/v3/apis/chainsaw.kyverno.io/v1alpha1 > ./.temp/.schemas/openapi/v3/apis/chainsaw.kyverno.io/v1alpha1.json
	@$(KIND) delete cluster --name schema
	@kubectl config use-context $(CURRENT_CONTEXT)

.PHONY: codegen-schemas-json
codegen-schemas-json: codegen-schemas-openapi ## Generate json schemas
	@$(PIP) install openapi2jsonschema
	@rm -rf ./.temp/.schemas/json
	@rm -rf ./.schemas/json
	@openapi2jsonschema ./.temp/.schemas/openapi/v2/schema.json --kubernetes --stand-alone --expanded -o ./.temp/.schemas/json
	@mkdir -p ./.schemas/json
	@cp ./.temp/.schemas/json/test-chainsaw-*.json ./.schemas/json
	@cp ./.temp/.schemas/json/configuration-chainsaw-*.json ./.schemas/json
	@echo Copy generated schemas to embed in the CLI... >&2
	@rm -rf pkg/data/schemas/json && mkdir -p pkg/data/schemas/json
	@cp ./.schemas/json/* pkg/data/schemas/json

.PHONY: codegen
codegen: codegen-crds codegen-deepcopy codegen-register codegen-mkdocs codegen-cli-docs codegen-api-docs codegen-schemas-json ## Rebuild all generated code and docs

.PHONY: verify-codegen
verify-codegen: codegen ## Verify all generated code and docs are up to date
	@echo Checking codegen is up to date... >&2
	@git --no-pager diff -- .
	@echo 'If this test fails, it is because the git diff is non-empty after running "make codegen".' >&2
	@echo 'To correct this, locally run "make codegen", commit the changes, and re-run tests.' >&2
	@git diff --quiet --exit-code -- .

##########
# MKDOCS #
##########

.PHONY: mkdocs-serve
mkdocs-serve: ## Generate and serve mkdocs website
	@echo Generate and servemkdocs website... >&2
	@$(PIP) install mkdocs
	@$(PIP) install --upgrade pip
	@$(PIP) install -U mkdocs-material mkdocs-redirects mkdocs-minify-plugin mkdocs-include-markdown-plugin lunr mkdocs-rss-plugin mike
	@mkdocs serve -f ./website/mkdocs.yaml

#########
# BUILD #
#########

CLI_BIN        := chainsaw
CGO_ENABLED    ?= 0
GOOS           ?= $(shell go env GOOS)
ifdef VERSION
LD_FLAGS       := "-s -w -X $(PACKAGE)/pkg/version.BuildVersion=$(VERSION)"
else
LD_FLAGS       := "-s -w"
endif

.PHONY: fmt
fmt: ## Run go fmt
	@echo Go fmt... >&2
	@go fmt ./...

.PHONY: vet
vet: ## Run go vet
	@echo Go vet... >&2
	@go vet ./...

.PHONY: $(CLI_BIN)
$(CLI_BIN): fmt vet
	@echo Build cli binary... >&2
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) go build -o ./$(CLI_BIN) -ldflags=$(LD_FLAGS) .

build: $(CLI_BIN) ## Build

########
# TEST #
########

.PHONY: tests
tests: $(CLI_BIN) ## Run tests
	@echo Running tests... >&2
	@go test ./... -race -coverprofile=coverage.out -covermode=atomic

########
# KIND #
########

KIND_IMAGE     ?= kindest/node:v1.28.0

.PHONY: kind-cluster
kind-cluster: $(KIND) ## Create kind cluster
	@echo Create kind cluster... >&2
	@$(KIND) create cluster --image $(KIND_IMAGE)

########
# HELP #
########

.PHONY: help
help: ## Shows the available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}'
