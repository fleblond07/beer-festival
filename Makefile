.PHONY: help backend frontend dev install build test test-backend test-frontend clean

# Default target
.DEFAULT_GOAL := help

# Colors for output
BLUE := \033[0;34m
GREEN := \033[0;32m
RESET := \033[0m

help: ## Show this help message
	@echo "$(BLUE)Beer Festival - Makefile commands$(RESET)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-15s$(RESET) %s\n", $$1, $$2}'

install: ## Install dependencies for frontend and backend
	@echo "$(BLUE)Installing frontend dependencies...$(RESET)"
	npm install
	@echo "$(GREEN)Frontend dependencies installed!$(RESET)"

backend: ## Start the backend server
	@echo "$(BLUE)Starting backend server...$(RESET)"
	cd backend && go run .

frontend: ## Start the frontend development server
	@echo "$(BLUE)Starting frontend development server...$(RESET)"
	npm run dev

dev: ## Start both backend and frontend concurrently
	@echo "$(BLUE)Starting backend and frontend...$(RESET)"
	@$(MAKE) -j2 backend frontend

build: ## Build the frontend for production
	@echo "$(BLUE)Building frontend...$(RESET)"
	npm run build
	@echo "$(GREEN)Frontend built successfully!$(RESET)"

test: test-backend test-frontend ## Run all tests

test-backend: ## Run backend tests
	@echo "$(BLUE)Running backend tests...$(RESET)"
	cd backend && go test -v

test-frontend: ## Run frontend unit tests
	@echo "$(BLUE)Running frontend unit tests...$(RESET)"
	npm run test

test-e2e: ## Run frontend e2e tests
	@echo "$(BLUE)Running frontend e2e tests...$(RESET)"
	npm run test:e2e

lint: ## Run linter on frontend code
	@echo "$(BLUE)Running linter...$(RESET)"
	npm run lint

format: ## Format frontend code
	@echo "$(BLUE)Formatting code...$(RESET)"
	npm run format

clean: ## Clean build artifacts and dependencies
	@echo "$(BLUE)Cleaning build artifacts...$(RESET)"
	rm -rf dist coverage node_modules
	@echo "$(GREEN)Clean complete!$(RESET)"
