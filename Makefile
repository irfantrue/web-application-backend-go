APP_NAME := app
AIR_DIR := air
AIR_PATH := $(AIR_DIR)/$(APP_NAME)
CMD_PATH := ./cmd/api/

.PHONY: all build run watch tidy lint test test-coverage migrate-up migrate-down clean

all: tidy lint test build

tidy:
	go mod tidy
	go mod verify

watch:
	air -c .air.toml

build-air:
	@mkdir -p $(AIR_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(AIR_PATH) $(CMD_PATH)
	@echo "Built: $(AIR_PATH) ($(VERSION))"

build:
	@mkdir -p bin
	go build -o bin/$(APP_NAME) $(CMD_PATH)
	@echo "Successfully built bin/$(APP_NAME)"

lint:
	@echo "Running linters..."
	golangci-lint run

test:
	go test -v -race ./...

test-coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean:
	@rm -rf bin/
	@rm -rf $(AIR_DIR)/
	@rm -f coverage.out coverage.html
	@echo "Cleaned up successfully."
