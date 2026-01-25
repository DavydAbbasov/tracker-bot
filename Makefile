.PHONY: help build run test clean deps migrate-up migrate-down docker-up docker-down

deps: ## Install dependencies
	go mod download
	go mod tidy

run: ## Run the application
	go run cmd/tracker-bot/main.go

test: ## Run tests
	go test -v ./...

clean: ## Clean build artifacts
	rm -rf bin/
	go clean

migrate-up: ## Run database migrations up
	goose -dir migrations/pgsql postgres "host=localhost port=5434 user=postgres password=12345 dbname=casino_notifications sslmode=disable" up

migrate-down: ## Run database migrations down
	goose -dir migrations/pgsql postgres "host=localhost port=5434 user=postgres password=12345 dbname=casino_notifications sslmode=disable" down

migrate-status: ## Check migration status
	goose -dir migrations/pgsql postgres "host=localhost port=5434 user=postgres password=12345 dbname=casino_notifications sslmode=disable" status

docker-up: ## Start docker compose services
	docker-compose up -d

docker-down: ## Stop docker compose services
	docker-compose down

dev: docker-up ## Start development environment
	@echo "Waiting for database to be ready..."
	@sleep 5
	make migrate-up
	@echo "Development environment is ready!"
	@echo "- PostgreSQL: localhost:5434"
	
logs: ## Show docker logs
	docker-compose logs -f