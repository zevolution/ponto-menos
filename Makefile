BINARY_NAME := ponto-menos
CMD_DIR := ./cmd
OUTPUT_DIR := ./bin
LAMBDA_LIST := punchclockschedule

.PHONY: clean

build: build-cli ## Build all binaries
	@for target in $(LAMBDA_LIST); do    \
	  make name=$${target} build-lambda; \
	done;

build-cli: ## Build CLI binary (e.g. 'make os=darwin arch=arm64 build-cli' or just 'make build-cli' for linux)
ifneq (${os},)
ifneq (${arch},)
	GOOS=${os} GOARCH=${arch} go build -o bin/$(BINARY_NAME) $(CMD_DIR)/cli/main.go
else
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME) $(CMD_DIR)/cli/main.go
endif
else
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME) $(CMD_DIR)/cli/main.go
endif

build-lambda: ## Build lambda binaries using param 'name' (e.g. make name=lambda-name build-lambda)
ifeq (${name}, $(filter ${name},$(LAMBDA_LIST)))
	GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bin/${name}/bootstrap $(CMD_DIR)/punchclockschedule/main.go
else
	@echo 'Lambda ${name} was not found'
endif

zip-lambda: ## Zip lambda binaries
ifeq (${name}, $(filter ${name},$(LAMBDA_LIST)))
	@zip -j bin/${name}/bootstrap.zip bin/${name}/bootstrap
else
	@echo 'Lambda ${name} was not found'
endif

clean: ## Remove previous build
	@-rm -vrf ${OUTPUT_DIR}

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'