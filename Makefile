-include .env
export

DB_URL=$(DATABASE_URL)
TEST_DB_URL=$(TEST_DATABASE_URL)
MIGRATIONS_PATH=db/migrations

.PHONY: migrate migrate-down migrate-force migrate-create migrate-version test test-integration sqlc server mock

migrate:
	@echo "Running migrations..."
	@migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

migrate-down:
	@echo "Running migrations down..."
	@migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down $(or $(steps),1)

migrate-force:
	@echo "Running migrations force..."
	@migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" force $(version)

migrate-create:
	@echo "Creating migration..."
	@migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)

migrate-version:
	@echo "Showing migration version..."
	@migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" version

test:
	@echo "Running unit tests..."
	@go test ./... -v -cover

test-integration:
	@echo "Running integration tests..."
	@go test -tags=integration -v ./... -cover

sqlc:
	@echo "Generating SQLC code..."
	@sqlc generate

server:
	@echo "Starting server..."
	@go run main.go

mock:
	@echo "Generating mock code..."
	@mockgen -destination db/mock/store.go -package mockdb github.com/santiagot714/SimpleBank/db/sqlc Store