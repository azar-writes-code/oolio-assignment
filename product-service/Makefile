-- Active: 1774208111445@@localhost@5432@invoice_db
down:
	@docker compose -f infra/docker/docker-compose.yml down -v

up:
	@docker compose -f infra/docker/docker-compose.yml up -d

logs:
	@docker compose -f infra/docker/docker-compose.yml logs -f

clean-db:
	@echo "Dropping database..."
	@go run ./cmd/db/main.go drop
	@echo "Waiting 2 seconds for file system to sync..."
	@sleep 2
	@echo "Running sql migrations..."
	@go run ./cmd/db/main.go migrate
	@sleep 2
	@echo "Generating sqlc files..."
	@make generate-sqlc
	

run-migrations:
	@go run ./cmd/db/main.go migrate

start:
	@echo "Stopping all services..."
	@make down 
	@sleep 5
	@echo "Starting all services..."
	@make up 
	@sleep 5
	@echo "Run sql migrations"
	@go run ./cmd/db/main.go migrate
	@echo "Starting air..."
	@air

generate-sqlc:
	@echo "Deleting old sqlc generated files..."
	@rm -rf pkg/server/rest/dtos/db/sqlc
	@echo "Waiting 5 seconds for file system to sync..."
	@sleep 5
	@echo "Generating new sqlc files..."
	@sqlc generate

restart-telemetry:
	@echo "Stopping loki..."
	@docker compose -f infra/docker/docker-compose.yml stop loki
	@echo "Stopping prometheus..."
	@docker compose -f infra/docker/docker-compose.yml stop prometheus
	@echo "Starting loki..."
	@docker compose -f infra/docker/docker-compose.yml start loki
	@echo "Starting prometheus..."
	@docker compose -f infra/docker/docker-compose.yml start prometheus

test:
	@go test ./tests/unit/... -v -coverprofile=coverage.out

end:
	@make down && clear

.PHONY: down up logs start end generate-sqlc restart-telemetry clean-db