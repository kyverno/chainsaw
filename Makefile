.DEFAULT_GOAL := build

#############
# VARIABLES #
#############

GIT_SHA                            := $(shell git rev-parse HEAD)
ORG                                ?= kyverno
PACKAGE                            ?= github.com/$(ORG)/chainsaw
CRDS_PATH                          := ${PWD}/.crds
CLI_BIN                            := chainsaw
CGO_ENABLED                        ?= 0
GOOS                               ?= $(shell go env GOOS)
ifdef VERSION
LD_FLAGS                           := "-s -w -X $(PACKAGE)/pkg/version.BuildVersion=$(VERSION)"
else
LD_FLAGS                           := "-s -w"
endif
KO_REGISTRY                        := ko.local
KO_TAGS                            := $(GIT_SHA)
KIND_IMAGE                         ?= kindest/node:v1.34.0

#########
# TOOLS #
#########

TOOLS_DIR                          := $(PWD)/.tools
CONTROLLER_GEN                     := $(TOOLS_DIR)/controller-gen
REGISTER_GEN                       := $(TOOLS_DIR)/register-gen
DEEPCOPY_GEN                       := $(TOOLS_DIR)/deepcopy-gen
CONVERSION_GEN                     := $(TOOLS_DIR)/conversion-gen
CODE_GEN_VERSION                   := v0.34.0
REFERENCE_DOCS                     := $(TOOLS_DIR)/genref
REFERENCE_DOCS_VERSION             := latest
KIND                               := $(TOOLS_DIR)/kind
KIND_VERSION                       := v0.30.0
KO                                 ?= $(TOOLS_DIR)/ko
KO_VERSION                         ?= v0.18.0
TOOLS                              := $(CONTROLLER_GEN) $(REGISTER_GEN) $(DEEPCOPY_GEN) $(CONVERSION_GEN) $(REFERENCE_DOCS) $(KIND) $(KO)
PIP                                ?= "pip"
ifeq ($(GOOS), darwin)
SED                                := gsed
else
SED                                := sed
endif
COMMA                              := ,

$(CONTROLLER_GEN):
	@echo Install controller-gen... >&2
	@cd ./hack/controller-gen && GOBIN=$(TOOLS_DIR) go install -buildvcs=false

$(REGISTER_GEN):
	@echo Install register-gen... >&2
	@GOBIN=$(TOOLS_DIR) go install k8s.io/code-generator/cmd/register-gen@$(CODE_GEN_VERSION)

$(DEEPCOPY_GEN):
	@echo Install deepcopy-gen... >&2
	@GOBIN=$(TOOLS_DIR) go install k8s.io/code-generator/cmd/deepcopy-gen@$(CODE_GEN_VERSION)

$(CONVERSION_GEN):
	@echo Install conversion-gen... >&2
	@GOBIN=$(TOOLS_DIR) go install k8s.io/code-generator/cmd/conversion-gen@$(CODE_GEN_VERSION)

$(REFERENCE_DOCS):
	@echo Install genref... >&2
	@GOBIN=$(TOOLS_DIR) go install github.com/kubernetes-sigs/reference-docs/genref@$(REFERENCE_DOCS_VERSION)

$(KIND):
	@echo Install kind... >&2
	@GOBIN=$(TOOLS_DIR) go install sigs.k8s.io/kind@$(KIND_VERSION)

$(KO):
	@echo Install ko... >&2
	@GOBIN=$(TOOLS_DIR) go install github.com/google/ko@$(KO_VERSION)

.PHONY: install-tools
install-tools: ## Install tools
install-tools: $(TOOLS)

.PHONY: clean-tools
clean-tools: ## Remove installed tools
	@echo Clean tools... >&2
	@rm -rf $(TOOLS_DIR)

###########
# CODEGEN #
###########

.PHONY: codegen-register
codegen-register: ## Generate types registrations
codegen-register: $(REGISTER_GEN)
	@echo Generate registration... >&2
	@$(REGISTER_GEN) \
		--go-header-file=./hack/boilerplate.go.txt \
		--output-file zz_generated.register.go \
		./pkg/apis/...

.PHONY: codegen-deepcopy
codegen-deepcopy: ## Generate deep copy functions
codegen-deepcopy: $(DEEPCOPY_GEN)
	@echo Generate deep copy functions... >&2
	@$(DEEPCOPY_GEN) \
		--go-header-file=./hack/boilerplate.go.txt \
		--output-file=zz_generated.deepcopy.go \
		./pkg/apis/...

.PHONY: codegen-conversion
codegen-conversion: ## Generate conversion functions
codegen-conversion: $(CONVERSION_GEN)
	@echo Generate conversion functions... >&2
	@$(CONVERSION_GEN) \
		--go-header-file=./hack/boilerplate.go.txt \
		--output-file=zz_generated.conversion.go \
		./pkg/apis/...

.PHONY: codegen-crds
codegen-crds: ## Generate CRDs
codegen-crds: $(CONTROLLER_GEN)
codegen-crds: codegen-deepcopy
codegen-crds: codegen-register
codegen-crds: codegen-conversion
	@echo Generate crds... >&2
	@rm -rf $(CRDS_PATH)
	@$(CONTROLLER_GEN) paths=./pkg/apis/... crd:crdVersions=v1,ignoreUnexportedFields=true,generateEmbeddedObjectMeta=false output:dir=$(CRDS_PATH)
	@echo Copy generated CRDs to embed in the CLI... >&2
	@rm -rf pkg/data/crds && mkdir -p pkg/data/crds
	@cp $(CRDS_PATH)/* pkg/data/crds

.PHONY: codegen-cli-docs
codegen-cli-docs: ## Generate CLI docs
codegen-cli-docs: build
	@echo Generate cli docs... >&2
	@rm -rf website/docs/reference/commands && mkdir -p website/docs/reference/commands
	@./$(CLI_BIN) docs -o website/docs/reference/commands --autogenTag=false

.PHONY: codegen-api-docs
codegen-api-docs: ## Generate markdown API docs
codegen-api-docs: $(REFERENCE_DOCS)
codegen-api-docs: codegen-deepcopy
codegen-api-docs: codegen-register
codegen-api-docs: codegen-conversion
	@echo Generate api docs... >&2
	@rm -rf ./website/docs/reference/apis
	@cd ./website/apis && $(REFERENCE_DOCS) -c config.yaml -f markdown -o ../docs/reference/apis

.PHONY: codegen-jp-docs
codegen-jp-docs: ## Generate JP docs
	@echo Generate jp docs... >&2
	@mkdir -p ./website/docs/reference/jp
	@go run ./website/jp/main.go > ./website/docs/reference/jp/functions.md

.PHONY: codegen-mkdocs
codegen-mkdocs: ## Generate mkdocs website
codegen-mkdocs: codegen-cli-docs
codegen-mkdocs: codegen-api-docs
codegen-mkdocs: codegen-jp-docs
	@echo Generate mkdocs website... >&2
	@$(PIP) install -r requirements.txt
	@mkdocs build -f ./website/mkdocs.yaml

.PHONY: codegen-schemas-openapi
codegen-schemas-openapi: ## Generate openapi schemas (v2 and v3)
codegen-schemas-openapi: CURRENT_CONTEXT = $(shell kubectl config current-context)
codegen-schemas-openapi: codegen-crds
codegen-schemas-openapi: $(KIND)
	@echo Generate openapi schema... >&2
	@rm -rf ./.temp/.schemas
	@mkdir -p ./.temp/.schemas/openapi/v2
	@mkdir -p ./.temp/.schemas/openapi/v3/apis/chainsaw.kyverno.io
	@$(KIND) create cluster --name schema --image $(KIND_IMAGE)
	@kubectl create -f $(CRDS_PATH)
	@sleep 15
	@kubectl get --raw /openapi/v2 > ./.temp/.schemas/openapi/v2/schema.json
	@kubectl get --raw /openapi/v3/apis/chainsaw.kyverno.io/v1alpha1 > ./.temp/.schemas/openapi/v3/apis/chainsaw.kyverno.io/v1alpha1.json
	@kubectl get --raw /openapi/v3/apis/chainsaw.kyverno.io/v1alpha2 > ./.temp/.schemas/openapi/v3/apis/chainsaw.kyverno.io/v1alpha2.json
	@$(KIND) delete cluster --name schema
	@kubectl config use-context $(CURRENT_CONTEXT) || true

.PHONY: codegen-schemas-json
codegen-schemas-json: ## Generate json schemas
codegen-schemas-json: codegen-schemas-openapi
	@echo Generate json schema... >&2
	@$(PIP) install -r requirements.txt
	@rm -rf ./.temp/.schemas/json
	@rm -rf ./.schemas/json
	@openapi2jsonschema .temp/.schemas/openapi/v3/apis/chainsaw.kyverno.io/v1alpha1.json --kubernetes --strict --stand-alone --expanded -o ./.temp/.schemas/json
	@openapi2jsonschema .temp/.schemas/openapi/v3/apis/chainsaw.kyverno.io/v1alpha2.json --kubernetes --strict --stand-alone --expanded -o ./.temp/.schemas/json
	@mkdir -p ./.schemas/json
	@cp ./.temp/.schemas/json/configuration-chainsaw-*.json ./.schemas/json
	@cp ./.temp/.schemas/json/steptemplate-chainsaw-*.json ./.schemas/json
	@cp ./.temp/.schemas/json/test-chainsaw-*.json ./.schemas/json
	@echo Copy generated schemas to embed in the CLI... >&2
	@rm -rf pkg/data/schemas/json && mkdir -p pkg/data/schemas/json
	@cp ./.schemas/json/* pkg/data/schemas/json

.PHONY: codegen-tests-catalog
codegen-tests-catalog: ## Generate tests catalog files
codegen-tests-catalog: $(CLI_BIN)
	@echo Generate tests catalog... >&2
	@./$(CLI_BIN) build docs --test-dir ./testdata/e2e --catalog ./testdata/e2e/examples/CATALOG.md

.PHONY: codegen
codegen: ## Rebuild all generated code and docs
codegen: codegen-api-docs
codegen: codegen-cli-docs
codegen: codegen-crds
codegen: codegen-deepcopy
codegen: codegen-mkdocs
codegen: codegen-register
codegen: codegen-schemas-json
codegen: codegen-tests-catalog
codegen: codegen-conversion

.PHONY: verify-codegen
verify-codegen: ## Verify all generated code and docs are up to date
verify-codegen: codegen
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
	@echo Generate and serve mkdocs website... >&2
	@$(PIP) install -r requirements.txt
	@mkdocs serve -f ./website/mkdocs.yaml

#########
# BUILD #
#########

.PHONY: fmt
fmt: ## Run go fmt
fmt: codegen-register
fmt: codegen-deepcopy
	@echo Go fmt... >&2
	@go fmt ./...

.PHONY: vet
vet: ## Run go vet
	@echo Go vet... >&2
	@go vet ./...

.PHONY: $(CLI_BIN)
$(CLI_BIN): fmt
$(CLI_BIN): vet
$(CLI_BIN): codegen-crds
	@echo Build cli binary... >&2
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) go build -o ./$(CLI_BIN) -ldflags=$(LD_FLAGS) .

build: $(CLI_BIN) ## Build

##############
# BUILD (KO) #
##############

.PHONY: build-ko
build-ko: ## Build Docker image with ko
build-ko: fmt
build-ko: vet
build-ko: $(KO)
	@echo "Build Docker image with ko..." >&2
	@LD_FLAGS=$(LD_FLAGS) KO_DOCKER_REPO=$(KO_REGISTRY) \
		$(KO) build . --preserve-import-paths --tags=$(KO_TAGS)

########
# TEST #
########

SET_FLAGS ?= --set env=poc --set clusterDirectory=my-cluster
SET_STRING_FLAGS ?= --set-string image.tag=01

.PHONY: tests
tests: ## Run tests
tests: $(CLI_BIN)
	@echo Running tests... >&2
	@go test ./... -race -coverprofile=coverage.out -covermode=atomic
	@go tool cover -html=coverage.out

.PHONY: e2e-tests
e2e-tests: ## Run e2e tests
e2e-tests: $(CLI_BIN)
	@echo Running e2e tests... >&2
	@./$(CLI_BIN) test ./testdata/e2e --remarshal --config ./testdata/e2e/config.yaml --values ./testdata/e2e/values.yaml $(SET_FLAGS) $(SET_STRING_FLAGS)

e2e-tests-no-cluster: ## Run e2e tests with --no-cluster
e2e-tests-no-cluster: $(CLI_BIN)
	@echo Running e2e tests with --no-cluster... >&2
	@./$(CLI_BIN) test testdata/e2e/examples/script-env --no-cluster --remarshal --config ./testdata/e2e/config.yaml --values ./testdata/e2e/values.yaml $(SET_FLAGS) $(SET_STRING_FLAGS)
	@./$(CLI_BIN) test testdata/e2e/examples/dynamic-clusters --no-cluster --remarshal --config ./testdata/e2e/config.yaml --values ./testdata/e2e/values.yaml

.PHONY: e2e-tests-ko 
e2e-tests-ko: ## Run e2e tests from a docker container
e2e-tests-ko: build-ko
	@echo Running e2e tests in docker... >&2
	@docker run \
		-v ./testdata/e2e/:/chainsaw/ \
		-v ${HOME}/.kube/:/etc/kubeconfig/ \
		-e KUBECONFIG=/etc/kubeconfig/config \
		--network=host \
		--user $(id -u):$(id -g) \
		--name chainsaw \
		--rm \
		ko.local/github.com/kyverno/chainsaw:$(KO_TAGS) test /chainsaw --remarshal --config /chainsaw/config.yaml --values /chainsaw/values.yaml --selector !no-ko-test

########	
# KIND #
########

.PHONY: kind-cluster
kind-cluster: ## Create kind cluster
kind-cluster: $(KIND)
	@echo Create kind cluster... >&2
	@$(KIND) create cluster --image $(KIND_IMAGE) --wait 1m

########
# HELP #
########

.PHONY: help
help: ## Shows the available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}'
