APP_NAME := pulse
BUILD_DIR := tmp

# Migration configuration
GOOSE_DRIVER ?= postgres
GOOSE_MIGRATION_DIR ?= ./store/migrations

# Build PostgreSQL connection string
# Format: postgres://user:password@host:port/dbname?sslmode=disable
GOOSE_DBSTRING ?= postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

.PHONY: build run run-mock clean

# Build the main application (root package only).
# We use "." instead of "./..." so Go only builds a single package
# and can write one binary to -o.
build:
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) .

run: build
	@./$(BUILD_DIR)/$(APP_NAME)

# Run the mock server as a separate program.
run-mock:
	@go run ./mock-server

clean:
	@rm -rf $(BUILD_DIR)


# Check if required environment variables are set
mig-check-env:
	@if [ -z "$(DB_USER)" ] || [ -z "$(DB_PASSWORD)" ] || [ -z "$(DB_HOST)" ] || [ -z "$(DB_PORT)" ] || [ -z "$(DB_NAME)" ]; then \
		echo "Error: Missing required database environment variables"; \
		echo "Please ensure .env file exists with: DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME"; \
		exit 1; \
	fi
	@if [ ! -f .env ]; then \
		echo "Warning: .env file not found"; \
		echo "Creating migrations may fail without proper database configuration"; \
	fi

# Run all pending migrations
mig-up: mig-check-env
	@echo "Running migrations..."
	@echo "Database: $(DB_HOST):$(DB_PORT)/$(DB_NAME)"
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$(GOOSE_DBSTRING)" goose -dir $(GOOSE_MIGRATION_DIR) up
	@echo "Migrations completed!"

# Rollback the last migration
mig-down: mig-check-env
	@echo "Rolling back last migration..."
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$(GOOSE_DBSTRING)" goose -dir $(GOOSE_MIGRATION_DIR) down
	@echo "Migration rolled back!"

# Show migration status
mig-status: mig-check-env
	@echo "Migration Status:"
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$(GOOSE_DBSTRING)" goose -dir $(GOOSE_MIGRATION_DIR) status

# Create a new migration file
mig-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Error: Migration name is required"; \
		echo "Usage: make mig-create NAME=migration_name"; \
		exit 1; \
	fi
	@echo "Creating new migration: $(NAME)"
	@GOOSE_DRIVER=$(GOOSE_DRIVER) goose -dir $(GOOSE_MIGRATION_DIR) create $(NAME) sql
	@echo "Migration file created!"

# Reset all migrations (rollback all, then run all)
# Use with caution - this will rollback ALL migrations
mig-reset: mig-check-env
	@echo "WARNING: This will rollback ALL migrations!"
	@echo "To proceed, run: make mig-reset-confirm"

# Confirmation target for reset (prevents accidental execution)
mig-reset-confirm: mig-check-env
	@echo "Rolling back all migrations..."
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$(GOOSE_DBSTRING)" goose -dir $(GOOSE_MIGRATION_DIR) reset
	@echo "All migrations rolled back!"

# Fix migration version (useful for fixing migration state)
mig-fix: mig-check-env
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: Version number is required"; \
		echo "Usage: make mig-fix VERSION=version_number"; \
		exit 1; \
	fi
	@echo "Fixing migration version to $(VERSION)..."
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$(GOOSE_DBSTRING)" goose -dir $(GOOSE_MIGRATION_DIR) fix $(VERSION)
	@echo "Migration version fixed!"

# Legacy aliases for backward compatibility
migu: mig-up
migd: mig-down
