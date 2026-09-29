APP_NAME := app
AIR_DIR := air
AIR_PATH := $(AIR_DIR)/$(APP_NAME)
CMD_PATH := ./cmd/api/

.PHONY: all build run watch

tidy:
	go mod tidy
	go mod verify

watch:
	air -c .air.toml

build-air:
	@mkdir -p $(AIR_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(AIR_PATH) $(CMD_PATH)
	@echo "Built: $(AIR_PATH) ($(VERSION))"
