GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: all

all: help

## General:
clean: ## Remove report files
	"$(CURDIR)/scripts/clean_files.sh"

createpkg: ## Remove report files
	"$(CURDIR)/scripts/create_pkg.sh"

## Go Helpers:
update: ## Update dependencies
	"$(CURDIR)/scripts/update_dependencies.sh"

sync: ## Sync project dependencies, equivalent to execute go mod tidy
	"$(CURDIR)/scripts/sync_dependencies.sh"

gofmt: ## Format project using gofmt and gofumpt
	"$(CURDIR)/scripts/gofmt.sh"

## Test:
test: ## Run the tests of the project
	"$(CURDIR)/scripts/unit_test.sh" $(PKG_NAME)

coverage: ## Run the tests of the project and export the coverage
	"$(CURDIR)/scripts/coverage.sh"

## Analyst:
lint: ## Run lint on your project
	"$(CURDIR)/scripts/linter.sh"

static_checks: ## Run static check and go vet on your project
	"$(CURDIR)/scripts/static_check.sh" $(PKG_NAME)

## Security:
vuln: ## Scan vulnerabilities
	"$(CURDIR)/scripts/vuln.sh" $(PKG_NAME)

## Docs:
docs: ## LocalDocumentation
	mkdocs serve

## Help:
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)
