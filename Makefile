# Root Makefile for Drifter Project

# Backend (Golang)
BACKEND_DIR := backend
BACKEND_MAIN := $(BACKEND_DIR)/main.go
BACKEND_BINARY := $(BACKEND_DIR)/bin/drifter
BACKEND_PID := backend.pid
FRONTEND_PID := frontend.pid

# Frontend (Next.js)
FRONTEND_DIR := frontend

.PHONY: all dev-backend run-backend build-backend clean-backend lint-backend \
        dev-frontend build-frontend lint-frontend format-frontend \
        start stop logs all lint deps test-backend test-frontend

## ---------- GLOBAL ----------
## Start everything (backend and frontend in the background)
all: deps start

## Install dependencies
deps:
	@echo "Installing Go modules..."
	cd $(BACKEND_DIR) && go mod download

	@echo "\nInstalling npm packages..."
	cd $(FRONTEND_DIR) && npm install

	@echo "\n✅ Dependencies installed successfully"

## Log files
LOG_DIR := logs
BACKEND_LOG := $(LOG_DIR)/backend.log
FRONTEND_LOG := $(LOG_DIR)/frontend.log

## Create logs directory if it doesn't exist
$(shell mkdir -p $(LOG_DIR))

## Start all services in the background
start: stop deps test-backend test-frontend
	@echo "Starting backend..."
	cd $(BACKEND_DIR) && go run . > ../$(BACKEND_LOG) 2>&1 & echo $$! > ../$(BACKEND_PID)

	@echo "Starting frontend..."
	cd $(FRONTEND_DIR) && npm run dev > ../$(FRONTEND_LOG) 2>&1 & echo $$! > ../$(FRONTEND_PID)

	@echo "\n✅ Services started in the background"
	@echo "Backend log: $(BACKEND_LOG)"
	@echo "Frontend log: $(FRONTEND_LOG)"
	@echo "Run 'make stop' to stop all services"
	@echo "Run 'make logs' to view logs"

## View logs
logs:
	@echo "=== Backend Log (Ctrl+C to exit) ==="
	@tail -f $(BACKEND_LOG) || (echo "No backend log found"; exit 0)

## Stop all running services
stop:
	@if [ -f $(BACKEND_PID) ]; then \
		echo "Stopping backend (PID: $$(cat $(BACKEND_PID)))"; \
		kill $$(cat $(BACKEND_PID)) 2>/dev/null || true; \
		rm -f $(BACKEND_PID); \
	fi
	@if [ -f $(FRONTEND_PID) ]; then \
		echo "Stopping frontend (PID: $$(cat $(FRONTEND_PID)))"; \
		pkill -P $$(cat $(FRONTEND_PID)) 2>/dev/null || true; \
		rm -f $(FRONTEND_PID); \
	fi
	@# Clean up any remaining processes
	@pkill -f "$(BACKEND_BINARY)" 2>/dev/null || true
	@pkill -f "next" 2>/dev/null || true
	@echo "✅ All services stopped"

## ---------- BACKEND ----------
## Run backend in development mode with hot-reload (blocks)
dev-backend: test-backend
	cd $(BACKEND_DIR) && air

## Run backend in development mode (blocks)
run-backend: test-backend
	cd $(BACKEND_DIR) && go run $(BACKEND_MAIN)

## Build backend binary
build-backend:
	cd $(BACKEND_DIR) && go build -o bin/drifter $(BACKEND_MAIN)

clean-backend:
	cd $(BACKEND_DIR) && rm -rf bin/
lint-backend:
	cd $(BACKEND_DIR) && golangci-lint run ./...

test-backend:
	cd $(BACKEND_DIR) && go test ./... -coverprofile=coverage.out
	cd $(BACKEND_DIR) && grep -v internal/world coverage.out > coverage_filtered.out
	@coverage=$$(cd $(BACKEND_DIR) && go tool cover -func=coverage_filtered.out | tail -1 | awk '{print substr($$3,1,length($$3)-1)}'); \
	echo "Backend coverage: $$coverage%"; \
	awk -v cov=$$coverage 'BEGIN { if (cov < 80) { print "Coverage below 80%"; exit 1 } }'

## ---------- FRONTEND ----------
## Run frontend in development mode (blocks)
dev-frontend: test-frontend
	cd $(FRONTEND_DIR) && npm run dev

build-frontend:
	cd $(FRONTEND_DIR) && npm run build

lint-frontend:
	cd $(FRONTEND_DIR) && npm run lint

test-frontend:
	cd $(FRONTEND_DIR) && npm test --silent

format-frontend:
	cd $(FRONTEND_DIR) && npm run format || npm run prettier

## ---------- COMBO ----------
lint: lint-backend lint-frontend

# ---------- SCHEMA GENERATION ----------
SCHEMA_INPUT := schemas/environment.schema.json
SCHEMA_OUTPUT := $(BACKEND_DIR)/internal/world/schema_gen.go

generate-go-schema:
	go-jsonschema -p world -o $(SCHEMA_OUTPUT) $(SCHEMA_INPUT)

# Full regen + build
regen:
	$(MAKE) generate-go-schema
	$(MAKE) build-backend