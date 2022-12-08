-include .env
export

.PHONY: build-cli fmt dep lint clean help fix-git

dep: ## Download app dependencies
	go mod tidy

build-cli: ## Build the cli binary
	go build -trimpath -o bin/cli cmd/cli/main.go

fmt: ## Reformat the code
	go fmt ./...

clean: ## Remove previous build
	rm -f ./bin/*

qa: lint vet  ## Run quality assurance checks

lint: ## Lint the code
	test -z $$(gofmt -l . | grep -v vendor/) || (echo "Formatting issues found in:" $$(gofmt -l . | grep -v vendor/) && exit 1)

vet: ## Vet the code
	go vet ./...
	golangci-lint run

help: ## Show available commands
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' ${MAKEFILE_LIST} | awk 'BEGIN {FS = ":.*?## "}; \
	{printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

fix-git: ## Fixes gitlab project URL
	git config --global url."git@bucket.digitalarsenal.net:".insteadOf "https://bucket.digitalarsenal.net"