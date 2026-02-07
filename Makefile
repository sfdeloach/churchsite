.PHONY: generate dev-up dev-up-detached dev-down dev-logs migrate-up migrate-down migrate-create seed test lint build watch clean help

generate: ## Generate Templ templates
	templ generate

dev-up: ## Start Docker containers w/ logs
	docker compose -f compose.yml -f compose.dev.yml up --build

dev-up-detached: ## Start Docker containers detached from logs
	docker compose -f compose.yml -f compose.dev.yml up -d --build

dev-down: ## Stop containers
	docker compose -f compose.yml -f compose.dev.yml down

dev-logs: ## Follow app logs
	docker compose -f compose.yml -f compose.dev.yml logs -f app

migrate-up: ## Run migrations
	docker compose -f compose.yml -f compose.dev.yml exec app go run ./cmd/server migrate up

migrate-down: ## Rollback last migration
	docker compose -f compose.yml -f compose.dev.yml exec app go run ./cmd/server migrate down

migrate-create: ## Create migration (usage: make migrate-create name=<name>)
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-create name=<name>"; exit 1; fi
	@mkdir -p migrations
	@touch migrations/$$(date +%Y%m%d%H%M%S)_$(name).up.sql
	@touch migrations/$$(date +%Y%m%d%H%M%S)_$(name).down.sql
	@echo "Created migration files for $(name)"

seed: ## Seed development data
	docker compose -f compose.yml -f compose.dev.yml exec app go run ./cmd/seed

test: ## Run tests
	go test -v -race -coverprofile=coverage.out ./...

lint: ## Run golangci-lint
	golangci-lint run ./...

build: ## Build binary
	templ generate
	CGO_ENABLED=0 go build -ldflags="-w -s" -o sachapel ./cmd/server

watch: ## Hot reload (templ + air)
	air

clean: ## Remove build artifacts
	rm -f sachapel
	rm -f coverage.out coverage.txt coverage.html
	rm -rf tmp/
	find . -name '*_templ.go' -delete

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
