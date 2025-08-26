.PHONY: help dev build clean test lint openapi generate generate-api generate-db db-up db-down migrate-up migrate-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development
dev: ## Start development environment
	docker-compose -f compose.base.yml -f compose.dev.overrides.yml up -d

dev-down: ## Stop development environment
	docker-compose -f compose.base.yml -f compose.dev.overrides.yml down

# Build
build: ## Build the application
	go build -o bin/api ./cmd/api

clean: ## Clean build artifacts
	rm -rf bin/ internal/api/gen.go

# Code generation
openapi: ## Combine OpenAPI files and generate combined spec
	@echo "Combining OpenAPI files..."
	swagger-cli bundle --dereference docs/base.yaml -o docs/openapi.json

generate: generate-api generate-db ## Generate all code (OpenAPI + SQLC)

generate-api: openapi ## Generate OpenAPI server code
	@echo "Generating OpenAPI server..."
	mkdir -p internal/api
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config codegen.yaml docs/openapi.json

generate-db: ## Generate SQLC database code
	@echo "Generating SQLC..."
	sqlc generate

# Database
db-up: ## Start database
	docker-compose -f compose.base.yml -f compose.dev.overrides.yml up -d postgres

db-down: ## Stop database
	docker-compose -f compose.base.yml -f compose.dev.overrides.yml down postgres

migrate-up: ## Run database migrations up
	migrate -path migrations -database "$$DATABASE_URL" up

migrate-down: ## Run database migrations down
	migrate -path migrations -database "$$DATABASE_URL" down

# Tools
install-tools: ## Install development tools
	go mod download
	npm install -g @apidevtools/swagger-cli
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Testing & Quality
test: ## Run tests
	go test -v ./...

lint: ## Run linter (if configured)
	@echo "Add golangci-lint configuration if needed"

