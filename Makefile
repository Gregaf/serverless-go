
FUNCTION_DIRECTORY := ./functions

FUNCTION_SUBDIRECTORIES := $(wildcard $(FUNCTION_DIRECTORY)/*)
BASENAME = $(shell basename $@)

build: clean tidy $(FUNCTION_SUBDIRECTORIES) ## Build all Serverless function binaries

tidy: ## Resolve "go.sum" file
	go mod tidy

$(FUNCTION_SUBDIRECTORIES):
	env CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/$(BASENAME) $@/main.go

clean: ## Remove generated go binaries
	rm -rf ./bin ./vendor go.sum

deploy: clean build ## Deploy the service
	sls deploy --verbose

help: ## Show this help menu
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: help build tidy clean deploy $(FUNCTION_SUBDIRECTORIES)
