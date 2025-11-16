.PHONY: help build run test clean proto fmt lint tidy deps install build-lambda package-lambda upload-lambda update-lambda deploy-lambda deploy-all clean-lambda

# Variables
BINARY_NAME=server
CMD_DIR=./cmd/server
PROTO_DIR=./proto
RPC_DIR=./rpc
GOPATH=$(shell go env GOPATH)
PROTOC=$(shell which protoc)

# Lambda deployment variables
LAMBDA_BINARY=bootstrap
LAMBDA_ZIP=function.zip
S3_BUCKET?=
S3_KEY?=quiz-backend/$(LAMBDA_ZIP)
AWS_REGION?=us-east-1
FUNCTION_NAME?=quiz-backend

# Colors for output
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[0;33m
NC := \033[0m # No Color

help: ## Show this help message
	@echo "$(BLUE)Available targets:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}'

build: ## Build the application
	@echo "$(BLUE)Building $(BINARY_NAME)...$(NC)"
	@go build -o $(BINARY_NAME) $(CMD_DIR)
	@echo "$(GREEN)Build complete!$(NC)"

run: ## Run the application
	@echo "$(BLUE)Running $(BINARY_NAME)...$(NC)"
	@go run $(CMD_DIR)

dev: ## Run the application in development mode (with auto-reload if air is installed)
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "$(YELLOW)Air not found, running with go run$(NC)"; \
		go run $(CMD_DIR); \
	fi

test: ## Run tests
	@echo "$(BLUE)Running tests...$(NC)"
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

clean: ## Clean build artifacts
	@echo "$(BLUE)Cleaning...$(NC)"
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out coverage.html
	@go clean
	@echo "$(GREEN)Clean complete!$(NC)"

proto: ## Generate protocol buffer code
	@echo "$(BLUE)Generating protocol buffer code...$(NC)"
	@if [ -z "$(PROTOC)" ]; then \
		echo "$(YELLOW)protoc not found. Installing...$(NC)"; \
		echo "Please install protoc: brew install protobuf"; \
		exit 1; \
	fi
	@export PATH="$$PATH:$(GOPATH)/bin" && \
	protoc --proto_path=$(PROTO_DIR) \
		--go_out=. --go_opt=module=github.com/rprajapati0067/quiz-game-backend \
		--go-grpc_out=. --go-grpc_opt=module=github.com/rprajapati0067/quiz-game-backend \
		$(PROTO_DIR)/*.proto
	@echo "$(GREEN)Protocol buffer code generated!$(NC)"

proto-install: ## Install protobuf plugins
	@echo "$(BLUE)Installing protocol buffer plugins...$(NC)"
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "$(GREEN)Plugins installed!$(NC)"

fmt: ## Format code
	@echo "$(BLUE)Formatting code...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)Format complete!$(NC)"

lint: ## Run linter
	@echo "$(BLUE)Running linter...$(NC)"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "$(YELLOW)golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(NC)"; \
		echo "Running go vet instead..."; \
		go vet ./...; \
	fi

tidy: ## Run go mod tidy
	@echo "$(BLUE)Running go mod tidy...$(NC)"
	@go mod tidy
	@echo "$(GREEN)Tidy complete!$(NC)"

deps: ## Download dependencies
	@echo "$(BLUE)Downloading dependencies...$(NC)"
	@go mod download
	@echo "$(GREEN)Dependencies downloaded!$(NC)"

install: ## Install the application
	@echo "$(BLUE)Installing $(BINARY_NAME)...$(NC)"
	@go install $(CMD_DIR)
	@echo "$(GREEN)Install complete!$(NC)"

check: fmt lint test ## Run format, lint, and test

build-lambda: ## Build for AWS Lambda (Linux binary)
	@echo "$(BLUE)Building for AWS Lambda...$(NC)"
	@GOOS=linux GOARCH=amd64 go build -o $(LAMBDA_BINARY) $(CMD_DIR)
	@echo "$(GREEN)Lambda build complete! Binary: $(LAMBDA_BINARY)$(NC)"

package-lambda: build-lambda ## Package Lambda function for deployment
	@echo "$(BLUE)Packaging Lambda function...$(NC)"
	@zip -j $(LAMBDA_ZIP) $(LAMBDA_BINARY)
	@echo "$(GREEN)Package created: $(LAMBDA_ZIP)$(NC)"

upload-lambda: package-lambda ## Upload Lambda package to S3
	@if [ -z "$(S3_BUCKET)" ]; then \
		echo "$(YELLOW)Error: S3_BUCKET is not set$(NC)"; \
		echo "$(YELLOW)Usage: make upload-lambda S3_BUCKET=my-bucket-name$(NC)"; \
		echo "$(YELLOW)Optional: S3_KEY=path/to/function.zip AWS_REGION=us-east-1$(NC)"; \
		exit 1; \
	fi
	@if ! command -v aws > /dev/null; then \
		echo "$(YELLOW)Error: AWS CLI not found. Install it first.$(NC)"; \
		exit 1; \
	fi
	@echo "$(BLUE)Uploading $(LAMBDA_ZIP) to s3://$(S3_BUCKET)/$(S3_KEY)...$(NC)"
	@aws s3 cp $(LAMBDA_ZIP) s3://$(S3_BUCKET)/$(S3_KEY) --region $(AWS_REGION)
	@echo "$(GREEN)Upload complete!$(NC)"
	@echo "$(YELLOW)Package location: s3://$(S3_BUCKET)/$(S3_KEY)$(NC)"

update-lambda: upload-lambda ## Upload to S3 and update Lambda function code
	@if [ -z "$(S3_BUCKET)" ] || [ -z "$(FUNCTION_NAME)" ]; then \
		echo "$(YELLOW)Error: S3_BUCKET and FUNCTION_NAME are required$(NC)"; \
		echo "$(YELLOW)Usage: make update-lambda S3_BUCKET=my-bucket FUNCTION_NAME=my-function$(NC)"; \
		exit 1; \
	fi
	@echo "$(BLUE)Updating Lambda function $(FUNCTION_NAME)...$(NC)"
	@aws lambda update-function-code \
		--function-name $(FUNCTION_NAME) \
		--s3-bucket $(S3_BUCKET) \
		--s3-key $(S3_KEY) \
		--region $(AWS_REGION)
	@echo "$(GREEN)Lambda function updated!$(NC)"

deploy-lambda: package-lambda ## Build and package for Lambda deployment
	@echo "$(BLUE)Ready for Lambda deployment$(NC)"
	@echo "$(YELLOW)Package created: $(LAMBDA_ZIP)$(NC)"
	@echo "$(YELLOW)Upload options:$(NC)"
	@echo "  - Manual: Upload $(LAMBDA_ZIP) to AWS Lambda console"
	@echo "  - S3: make upload-lambda S3_BUCKET=my-bucket"
	@echo "  - S3 + Update: make update-lambda S3_BUCKET=my-bucket FUNCTION_NAME=my-function"
	@echo "$(YELLOW)Lambda Configuration:$(NC)"
	@echo "  - Runtime: provided.al2023 or provided.al2"
	@echo "  - Handler: bootstrap"
	@echo "  - Architecture: x86_64"

deploy-all: upload-lambda update-lambda ## Build, package, upload to S3, and update Lambda function
	@echo "$(GREEN)Deployment complete!$(NC)"

clean-lambda: ## Clean Lambda build artifacts
	@echo "$(BLUE)Cleaning Lambda artifacts...$(NC)"
	@rm -f $(LAMBDA_BINARY) $(LAMBDA_ZIP)
	@echo "$(GREEN)Lambda artifacts cleaned!$(NC)"

all: clean proto tidy deps build ## Clean, generate proto, tidy, download deps, and build

