MIGRATION_DIR=./cmd/migrate/migrations

.PHONY: migration-create
migration:
	@goose -dir $(MIGRATION_DIR) create $(word 2,$(MAKECMDGOALS)) sql

.PHONY: migrate-up
migrate-up:
	@goose -dir $(MIGRATION_DIR) postgres "$$DB_URL" up

.PHONY: migrate-down
migrate-down:
	@goose -dir $(MIGRATION_DIR) postgres "$$DB_URL" down

.PHONY: seed
seed:
	@go run cmd/migrate/seed/main.go

.PHONY: gen-docs
gen-docs:
	@swag init -g ./api/main.go -d cmd,internal
