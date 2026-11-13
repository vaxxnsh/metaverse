GO ?= go
GOOSE ?= goose

REALTIME_DIR := services/realtime
API_DIR := services/api
SHARED_DIR := shared
BIN_DIR := bin
COMPOSE_FILE := docker/docker-compose.yml
ENV_FILE ?= $(API_DIR)/.env
MIGRATIONS_DIR := shared/sql/schema

.PHONY: help fmt test check api-test shared-test realtime-test api-run realtime-run api-build realtime-build migrate-up migrate-status docker-up docker-down docker-logs docker-ps

help:
	@printf '%s\n' 'Targets:' \
		'  make fmt             Format Go source files' \
		'  make test            Run all Go tests' \
		'  make check           Run formatting checks and tests' \
		'  make api-test        Run API-service tests' \
		'  make realtime-test   Run realtime-service tests' \
		'  make shared-test     Run shared-module tests' \
		'  make api-run         Start the API server' \
		'  make realtime-run    Start the realtime WebSocket server' \
		'  make api-build       Build the API server into bin/' \
		'  make realtime-build  Build the realtime server into bin/' \
		'  make migrate-up      Apply pending database migrations' \
		'  make migrate-status  Show database migration status' \
		'  make docker-up       Start Docker Compose services' \
		'  make docker-down     Stop Docker Compose services' \
		'  make docker-logs     Follow Docker Compose logs' \
		'  make docker-ps       Show Docker Compose service status'

fmt:
	@gofmt -w $$(find $(API_DIR) $(REALTIME_DIR) $(SHARED_DIR) -name '*.go' -type f)

test: api-test realtime-test shared-test

check:
	@gofmt -d $$(find $(API_DIR) $(REALTIME_DIR) $(SHARED_DIR) -name '*.go' -type f) | diff -u /dev/null -
	@$(MAKE) --no-print-directory test

api-test:
	@cd $(API_DIR) && $(GO) test ./...

realtime-test:
	@cd $(REALTIME_DIR) && $(GO) test ./...

shared-test:
	@cd $(SHARED_DIR) && $(GO) test ./...

realtime-run:
	@cd $(REALTIME_DIR) && $(GO) run ./cmd/server

api-run:
	@cd $(API_DIR) && $(GO) run ./cmd/server

api-build:
	@mkdir -p $(BIN_DIR)
	@cd $(API_DIR) && $(GO) build -o ../../$(BIN_DIR)/api ./cmd/server

realtime-build:
	@mkdir -p $(BIN_DIR)
	@cd $(REALTIME_DIR) && $(GO) build -o ../../$(BIN_DIR)/realtime ./cmd/server

migrate-up:
	@test -f "$(ENV_FILE)" || { echo "Missing $(ENV_FILE); copy services/api/.env.example first."; exit 1; }
	@set -a; . "$(ENV_FILE)"; set +a; test -n "$$DATABASE_URL" || { echo "DATABASE_URL is required in $(ENV_FILE)."; exit 1; }; $(GOOSE) -dir "$(MIGRATIONS_DIR)" postgres "$$DATABASE_URL" up

migrate-status:
	@test -f "$(ENV_FILE)" || { echo "Missing $(ENV_FILE); copy services/api/.env.example first."; exit 1; }
	@set -a; . "$(ENV_FILE)"; set +a; test -n "$$DATABASE_URL" || { echo "DATABASE_URL is required in $(ENV_FILE)."; exit 1; }; $(GOOSE) -dir "$(MIGRATIONS_DIR)" postgres "$$DATABASE_URL" status

docker-up:
	@docker compose -f $(COMPOSE_FILE) up -d

docker-down:
	@docker compose -f $(COMPOSE_FILE) down

docker-logs:
	@docker compose -f $(COMPOSE_FILE) logs -f

docker-ps:
	@docker compose -f $(COMPOSE_FILE) ps
