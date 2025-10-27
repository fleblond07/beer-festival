.PHONY: help backend frontend dev install build test test-backend test-frontend clean

.DEFAULT_GOAL := help

BLUE := \033[0;34m
GREEN := \033[0;32m
RESET := \033[0m

help:
	@echo "$(BLUE)Beer Festival - Makefile commands$(RESET)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-15s$(RESET) %s\n", $$1, $$2}'

install: ## Install dependencies for frontend
	@echo "$(BLUE)Installing frontend dependencies...$(RESET)"
	cd frontend && npm install
	@echo "$(GREEN)Frontend dependencies installed!$(RESET)"

backend: ## Start the backend server
	@echo "$(BLUE)Starting backend server...$(RESET)"
	cd backend && go run .

frontend: ## Start the frontend development server
	@echo "$(BLUE)Starting frontend development server...$(RESET)"
	cd frontend && npm run dev

dev: ## Start both backend and frontend concurrently
	@echo "$(BLUE)Starting backend and frontend...$(RESET)"
	@$(MAKE) -j2 backend frontend

build: ## Build the frontend for production
	@echo "$(BLUE)Building frontend...$(RESET)"
	cd frontend && npm run build
	@echo "$(GREEN)Frontend built successfully!$(RESET)"

test: test-backend test-frontend ## Run all tests

test-backend: ## Run backend tests
	@echo "$(BLUE)Running backend tests...$(RESET)"
	cd backend && go test -v

test-frontend: ## Run frontend unit tests
	@echo "$(BLUE)Running frontend unit tests...$(RESET)"
	cd frontend && npm run test

test-e2e: ## Run frontend e2e tests
	@echo "$(BLUE)Running frontend e2e tests...$(RESET)"
	cd frontend && npm run test:e2e

lint: ## Run linter on frontend code
	@echo "$(BLUE)Running linter...$(RESET)"
	cd frontend && npm run lint

format: ## Format frontend code
	@echo "$(BLUE)Formatting code...$(RESET)"
	cd frontend && npm run format

clean: ## Clean build artifacts and dependencies
	@echo "$(BLUE)Cleaning build artifacts...$(RESET)"
	rm -rf frontend/dist frontend/coverage frontend/node_modules
	@echo "$(GREEN)Clean complete!$(RESET)"
