#!make

# Variables
BINARY_NAME := file-catalog
SRC_DIR := .
BIN_DIR := ./bin
SRC_FILES := $(wildcard $(SRC_DIR)/*.go)

# Default target
build: $(BIN_DIR)/$(BINARY_NAME)

# Create binary in bin directory
$(BIN_DIR)/$(BINARY_NAME): $(SRC_FILES)
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BINARY_NAME) $(SRC_FILES)
	@echo "Built $(BINARY_NAME) in $(BIN_DIR)"

# Clean target
clean:
	@rm -rf $(BIN_DIR)
	@echo "Cleaned $(BIN_DIR)"

.PHONY: build clean

.PHONY: run
run: clean build
	./bin/file-catalog hash
