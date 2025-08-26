.PHONY: help dev build clean test lint openapi generate generate-api generate-db db-up db-down migrate-up migrate-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development
dev: ## Start development environment
	docker compose -f compose.base.yml -f compose.dev.overrides.yml up -d

dev-down: ## Stop development environment
	docker compose -f compose.base.yml -f compose.dev.overrides.yml down

# Helper function to add entry to .gitignore if not exists
define add_to_gitignore
	@if ! grep -q "$(1)" .gitignore 2>/dev/null; then \
		echo "$(1)" >> .gitignore; \
		echo "Added $(1) to .gitignore"; \
	else \
		echo "$(1) already exists in .gitignore"; \
	fi
endef

env: ## Create environment files and update .gitignore
	@echo "Creating environment files..."
	@if [ ! -f .env.example ]; then \
		echo "Creating .env.example..."; \
		echo "# Database Configuration" > .env.example; \
		echo "DB_HOST=localhost" >> .env.example; \
		echo "DB_PORT=5432" >> .env.example; \
		echo "DB_NAME=go_api_starter" >> .env.example; \
		echo "DB_USER=postgres" >> .env.example; \
		echo "DB_PASSWORD=password" >> .env.example; \
		echo "" >> .env.example; \
		echo "# Server Configuration" >> .env.example; \
		echo "SERVER_PORT=8080" >> .env.example; \
		echo "SERVER_HOST=localhost" >> .env.example; \
		echo "" >> .env.example; \
		echo "# JWT Configuration" >> .env.example; \
		echo "JWT_SECRET=your-secret-key-here" >> .env.example; \
		echo "JWT_EXPIRY=24h" >> .env.example; \
		echo "" >> .env.example; \
		echo "# Environment" >> .env.example; \
		echo "ENV=development" >> .env.example; \
		echo "" >> .env.example; \
		echo "# Logging" >> .env.example; \
		echo "LOG_LEVEL=debug" >> .env.example; \
	fi
	@echo "Copying .env.example to environment files..."
	@cp .env.example .env
	@cp .env.example .env.development
	@cp .env.example .env.production
	@echo "Updating .gitignore..."
	$(call add_to_gitignore,.env)
	$(call add_to_gitignore,.env.development)
	$(call add_to_gitignore,.env.production)
	@echo "Environment files created successfully!"



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
	docker compose -f compose.base.yml -f compose.dev.overrides.yml up -d postgres

db-down: ## Stop database
	docker compose -f compose.base.yml -f compose.dev.overrides.yml down postgres

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

