## Set defaults
export GO111MODULE := on

PROJECT_NAME := "api"

MAIN_PACKAGE := "$(PROJECT_NAME)-server"

APP_CONFIG_FILE_NAME := $(or $(APP_CONFIG_FILE_NAME),"default")

SERVER_PORT = $(or $(APP_PORT),3000)

# Fetch OS info
GOVERSION=$(shell go version)
UNAME_OS=$(shell go env GOOS)
UNAME_ARCH=$(shell go env GOARCH)

DOCS_DIR := "docs"
GIT_HOOKS_DIR := "scripts/git_hooks"

BUILD_SCRIPT := "scripts/build.sh"

# Go binary. Change this to experiment with different versions of Go.
GO = go

VERBOSE = 0
Q 		= $(if $(filter 1,$VERBOSE),,@)
M 		= $(shell printf "\033[34;1m▶\033[0m")

BIN 	 = $(CURDIR)/bin
PKGS     = $(or $(PKG),$(shell $(GO) list ./...))

$(BIN)/%: | $(BIN) ; $(info $(M) building package: $(PACKAGE) ...)
	tmp=$$(mktemp -d); \
       go get $(PACKAGE)
	   env GOBIN=$(BIN) go install $(PACKAGE) \
		|| ret=$$?; \
	   rm -rf $$tmp ; exit $$ret

$(BIN)/goimports: PACKAGE=golang.org/x/tools/cmd/goimports@v0.1.5

GOIMPORTS = $(BIN)/goimports

$(BIN)/golint: PACKAGE=golang.org/x/lint/golint

GOLINT = $(BIN)/golint

GOFILES ?= $(shell find . -type f -name '*.go' -not -path "./models/*" -not -path "restapi/*" -not -path "./vendor/*" -not -path "./statik/*")

.PHONY: deps
deps:
	@go install golang.org/x/lint/golint@latest

.PHONY: build-info
build-info:
	@echo "\nBuild Info:\n"
	@echo "\t\033[33mOS\033[0m: $(UNAME_OS)"
	@echo "\t\033[33mArch\033[0m: $(UNAME_ARCH)"
	@echo "\t\033[33mGo Version\033[0m: $(GOVERSION)\n"

.PHONY: setup-git-hooks
setup-git-hooks:
	@chmod +x $(GIT_HOOKS_DIR)/*
	@git config core.hooksPath $(GIT_HOOKS_DIR)

.PHONY: goimports ## Run goimports and format files
goimports: | $(GOIMPORTS) ; $(info $(M) running goimports ...) @
	$Q $(GOIMPORTS) -w $(GOFILES)

.PHONY: goimports-check ## Check goimports without modifying the files
goimports-check: | $(GOIMPORTS) ; $(info $(M) running goimports -l ...) @
	$(eval FILES=$(shell sh -c '$(GOIMPORTS) -l $(GOFILES)'))
	@$(if $(strip $(FILES)), echo $(FILES); exit 1, echo "goimports check passed")

.PHONY: lint-check ## Run golint check
lint-check: | $(GOLINT) ; $(info $(M) running golint ...) @
	$Q $(GOLINT) -set_exit_status $(PKGS)

.PHONY: tidy
tidy:
	@echo "\n + Tidying Go Mod...\n"
	$Q @$(GO) mod tidy

.PHONY: clean
clean:
	@echo "\n + Cleaning generated files...\n"
	$Q @git clean -Xfq

.PHONY: generate
generate:
	@echo "\n + Running Swagger Codegen...\n"
	$Q @swagger generate server -f swagger/swagger.yml --quiet --main-package=$(MAIN_PACKAGE)

.PHONY: build
build:
	@CURR_DIR=$(CURDIR) $(BUILD_SCRIPT)

.PHONY: server
server: build
	$Q @$(GO) run cmd/$(MAIN_PACKAGE)/main.go --port $(SERVER_PORT)

.PHONY: dev-server
dev-server: build
	$Q CONFIG_FILE_PATH="$(CURDIR)/application/configuration/" CONFIG_FILE="$(APP_CONFIG_FILE_NAME)" ENV="BETA" $(GO) run cmd/$(MAIN_PACKAGE)/main.go --port $(SERVER_PORT)

.PHONY: release
release: clean build goimports-check lint-check
	$Q CGO_ENABLED=1 GOOS=$(UNAME_OS) GOARCH=$(UNAME_ARCH) $(GO) build -o $(BIN)/$(PROJECT_NAME) cmd/$(MAIN_PACKAGE)/main.go

.PHONY: release-quick
release-quick: clean build
	$Q CGO_ENABLED=1 GOOS=$(UNAME_OS) GOARCH=$(UNAME_ARCH) $(GO) build -o $(BIN)/$(PROJECT_NAME) cmd/$(MAIN_PACKAGE)/main.go
