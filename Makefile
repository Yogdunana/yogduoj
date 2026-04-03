# ============================================================
#  YogduOJ Makefile
#  Copyright (c) Yogdunana-悠渡
# ============================================================

.PHONY: dev build backend frontend judge test lint clean deploy help

# ---------- Default ----------
.DEFAULT_GOAL := help

# ---------- Variables ----------
GO        := go
GOFLAGS   := -v
NPM       := npm
PNPM      := pnpm
DOCKER    := docker
COMPOSE   := docker-compose

BACKEND_DIR  := backend
FRONTEND_DIR := frontend
JUDGE_DIR    := judge
DEPLOY_DIR   := deploy

GREEN  := \033[0;32m
YELLOW := \033[0;33m
CYAN   := \033[0;36m
RESET  := \033[0m

# ---------- Help ----------
help: ## Show this help message
	@echo ""
	@echo "$(CYAN)  YogduOJ - Online Judge System$(RESET)"
	@echo "$(CYAN)  Copyright Yogdunana-悠渡$(RESET)"
	@echo ""
	@echo "$(GREEN)Usage:$(RESET)"
	@echo "  make [target]"
	@echo ""
	@echo "$(GREEN)Targets:$(RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-12s$(RESET) %s\n", $$1, $$2}'
	@echo ""

# ---------- Development ----------
dev: ## Start development environment with Docker Compose
	@echo "$(GREEN)[Dev] Starting development environment...$(RESET)"
	$(COMPOSE) -f $(DEPLOY_DIR)/docker-compose.dev.yml up --build -d
	@echo "$(GREEN)[Dev] Development environment is running.$(RESET)"

dev-down: ## Stop development environment
	@echo "$(YELLOW)[Dev] Stopping development environment...$(RESET)"
	$(COMPOSE) -f $(DEPLOY_DIR)/docker-compose.dev.yml down
	@echo "$(GREEN)[Dev] Development environment stopped.$(RESET)"

dev-logs: ## Tail development logs
	$(COMPOSE) -f $(DEPLOY_DIR)/docker-compose.dev.yml logs -f

# ---------- Build ----------
build: build-backend build-frontend build-judge ## Build all services
	@echo "$(GREEN)[Build] All services built successfully.$(RESET)"

build-backend: ## Build backend binary
	@echo "$(GREEN)[Build] Building backend...$(RESET)"
	cd $(BACKEND_DIR) && $(GO) build $(GOFLAGS) -o ../bin/yogduoj-server ./cmd/server
	@echo "$(GREEN)[Build] Backend binary: bin/yogduoj-server$(RESET)"

build-frontend: ## Build frontend assets
	@echo "$(GREEN)[Build] Building frontend...$(RESET)"
	cd $(FRONTEND_DIR) && $(PNPM) install && $(PNPM) run build
	@echo "$(GREEN)[Build] Frontend built: $(FRONTEND_DIR)/dist$(RESET)"

build-judge: ## Build judge service binary
	@echo "$(GREEN)[Build] Building judge service...$(RESET)"
	cd $(JUDGE_DIR) && $(GO) build $(GOFLAGS) -o ../bin/yogduoj-judge ./cmd/judge
	@echo "$(GREEN)[Build] Judge binary: bin/yogduoj-judge$(RESET)"

# ---------- Run (local development) ----------
backend: ## Build and run backend locally
	@echo "$(GREEN)[Run] Starting backend server...$(RESET)"
	cd $(BACKEND_DIR) && $(GO) run ./cmd/server

frontend: ## Build and run frontend dev server
	@echo "$(GREEN)[Run] Starting frontend dev server...$(RESET)"
	cd $(FRONTEND_DIR) && $(PNPM) install && $(PNPM) run dev

judge: ## Build and run judge service locally
	@echo "$(GREEN)[Run] Starting judge service...$(RESET)"
	cd $(JUDGE_DIR) && $(GO) run ./cmd/judge

# ---------- Testing ----------
test: test-backend test-frontend test-judge ## Run all tests
	@echo "$(GREEN)[Test] All tests passed.$(RESET)"

test-backend: ## Run backend tests
	@echo "$(GREEN)[Test] Running backend tests...$(RESET)"
	cd $(BACKEND_DIR) && $(GO) test ./... -race -coverprofile=coverage.out
	@echo "$(GREEN)[Test] Backend tests passed.$(RESET)"

test-frontend: ## Run frontend tests
	@echo "$(GREEN)[Test] Running frontend tests...$(RESET)"
	cd $(FRONTEND_DIR) && $(PNPM) run test:unit
	@echo "$(GREEN)[Test] Frontend tests passed.$(RESET)"

test-judge: ## Run judge service tests
	@echo "$(GREEN)[Test] Running judge service tests...$(RESET)"
	cd $(JUDGE_DIR) && $(GO) test ./... -race -coverprofile=coverage.out
	@echo "$(GREEN)[Test] Judge tests passed.$(RESET)"

# ---------- Linting ----------
lint: lint-backend lint-frontend lint-judge ## Run all linters
	@echo "$(GREEN)[Lint] All lint checks passed.$(RESET)"

lint-backend: ## Lint backend Go code
	@echo "$(GREEN)[Lint] Linting backend...$(RESET)"
	cd $(BACKEND_DIR) && golangci-lint run ./...
	@echo "$(GREEN)[Lint] Backend lint passed.$(RESET)"

lint-frontend: ## Lint frontend code
	@echo "$(GREEN)[Lint] Linting frontend...$(RESET)"
	cd $(FRONTEND_DIR) && $(PNPM) run lint
	@echo "$(GREEN)[Lint] Frontend lint passed.$(RESET)"

lint-judge: ## Lint judge Go code
	@echo "$(GREEN)[Lint] Linting judge service...$(RESET)"
	cd $(JUDGE_DIR) && golangci-lint run ./...
	@echo "$(GREEN)[Lint] Judge lint passed.$(RESET)"

# ---------- Docker ----------
docker-build: ## Build all Docker images
	@echo "$(GREEN)[Docker] Building all images...$(RESET)"
	$(DOCKER) compose -f $(DEPLOY_DIR)/docker-compose.yml build
	@echo "$(GREEN)[Docker] All images built.$(RESET)"

docker-up: ## Start all services with Docker Compose
	@echo "$(GREEN)[Docker] Starting services...$(RESET)"
	$(DOCKER) compose -f $(DEPLOY_DIR)/docker-compose.yml up -d
	@echo "$(GREEN)[Docker] All services running.$(RESET)"

docker-down: ## Stop all Docker services
	@echo "$(YELLOW)[Docker] Stopping services...$(RESET)"
	$(DOCKER) compose -f $(DEPLOY_DIR)/docker-compose.yml down
	@echo "$(GREEN)[Docker] All services stopped.$(RESET)"

# ---------- Database ----------
migrate: ## Run database migrations
	@echo "$(GREEN)[DB] Running migrations...$(RESET)"
	cd $(BACKEND_DIR) && $(GO) run ./cmd/migrate
	@echo "$(GREEN)[DB] Migrations completed.$(RESET)"

migrate-down: ## Rollback last database migration
	@echo "$(YELLOW)[DB] Rolling back last migration...$(RESET)"
	cd $(BACKEND_DIR) && $(GO) run ./cmd/migrate --down
	@echo "$(GREEN)[DB] Rollback completed.$(RESET)"

# ---------- Clean ----------
clean: ## Clean all build artifacts and caches
	@echo "$(YELLOW)[Clean] Removing build artifacts...$(RESET)"
	rm -rf bin/
	rm -rf $(FRONTEND_DIR)/dist/
	rm -rf $(FRONTEND_DIR)/node_modules/.cache/
	rm -rf $(BACKEND_DIR)/tmp/
	rm -rf $(JUDGE_DIR)/tmp/
	rm -rf $(JUDGE_DIR)/workspace/
	$(GO) clean -cache
	@echo "$(GREEN)[Clean] All artifacts cleaned.$(RESET)"

# ---------- Deploy ----------
deploy: ## Deploy to production
	@echo "$(GREEN)[Deploy] Building and deploying to production...$(RESET)"
	$(DOCKER) compose -f $(DEPLOY_DIR)/docker-compose.yml up -d --build --force-recreate
	@echo "$(GREEN)[Deploy] Production deployment complete.$(RESET)"

deploy-stop: ## Stop production deployment
	@echo "$(YELLOW)[Deploy] Stopping production...$(RESET)"
	$(DOCKER) compose -f $(DEPLOY_DIR)/docker-compose.yml down
	@echo "$(GREEN)[Deploy] Production stopped.$(RESET)"

# ---------- Code Generation ----------
generate: ## Run code generation (swagger, protobuf, etc.)
	@echo "$(GREEN)[Generate] Running code generation...$(RESET)"
	cd $(BACKEND_DIR) && swag init -g cmd/server/main.go -o docs
	@echo "$(GREEN)[Generate] Code generation complete.$(RESET)"
