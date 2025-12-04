APP_NAME := pulse
BUILD_DIR := tmp

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
