.PHONY: clean help test
.DEFAULT: help

gen: ## go generate ./...
	go generate ./...

test: ## Run go vet, and test the whole repo
	go vet ./...
	go test ./...

integration-test-safe: ## Run integration tests. Requires that you fill out tests/key.yaml (see tests/key.template.yaml). These hit prod (there is no sandbox), but every test doesn't involve using money
	INTEGRATION=1 go test ./tests/...

integration-test: ## Run riskier integration tests. Requires that you fill out tests/key.yaml (see tests/key.template.yaml). These hit prod (there is no sandbox), and will create orders that will basically never fill (buying BTC for $0.01)
	INTEGRATION=1 go test ./tests/...

clean: gen ## tidy modules, delete the bin folder, go generate
	go mod tidy

help: ## Print help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@printf "\033[36m%-30s\033[0m %s\n" "(target)" "Build a target binary: $(targets)"
