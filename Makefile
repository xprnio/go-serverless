SOURCE := cmd/main.go

BUILD_DIR := bin
BUILD := $(BUILD_DIR)/go-serverless

GOCMD := go
GOBUILD := $(GOCMD) build

.PHONY: run build

build:
	mkdir -p "$(BUILD_DIR)"
	$(GOBUILD) -o "$(BUILD)" $(SOURCE)

run: build
	$(BUILD)

