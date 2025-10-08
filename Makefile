SHELL = /bin/bash

# Set up development environment
setup:
	@lefthook install
	@echo "✅ Development environment ready"

# Run Cart Service application
run:
	@echo "🚀 Running Service ..."
	@set -a && source ./.env && go run ./...

# Run Cart Service dependencies
dev:
	@echo "📦 Starting Cart Service dependencies..."
	@docker compose down
	@docker compose up -d

# Stop Cart Service dependencies
down:
	@echo "✋ Shutting down Cart Service dependencies..."
	@docker compose down

# Integration tests
itest:
	@echo "🧪 Testing Cart Service ..."
	@go clean -testcache
	@go test -v -run Integration ./...
	@echo "✅ Tests passed"

# Format Go code using golangci-lint
fmt:
	@echo "🔧 Formatting Go code..."
	@golangci-lint fmt
	@echo "✅ Code formatting complete"

# Run linter checks using gloangci-lint
lint:
	@echo "🔨 Running linter checks..."
	@golangci-lint run
	@echo "✅ Linting complete"

# Fix linting if possible and format the source code
fix: 
	@echo "🛠️ Fix linter issues and formatting the code..."
	@golangci-lint run --fix
	@echo "✅ Fixing complete"

# CI Build discarding artefacts
check-build:
	@echo ⏳ "Building..."
	@go build -o /dev/null ./...
	@echo "✅ Building complete"

.PHONY: setup run dev down itest fmt lint fix check-build