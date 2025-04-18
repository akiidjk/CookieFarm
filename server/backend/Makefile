BIN_DIR := ./bin
CMD_DIR := ./cmd/api
LOGS_DIR := ./logs
MAIN_FILE := main.go
BINARY_NAME := clientbackend
VERSION := 0.1.0

RESET := \033[0m
BOLD := \033[1m
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
CYAN := \033[36m

help:
	@echo -e "$(BOLD)Available commands:$(RESET)"
	@echo -e "  $(CYAN)make build$(RESET)     - Compile the Go project"
	@echo -e "  $(CYAN)make run$(RESET)       - Compile and run the binary"
	@echo -e "  $(CYAN)make run-prod$(RESET)  - Compile and run the binary for production"
	@echo -e "  $(CYAN)make test$(RESET)      - Run tests"
	@echo -e "  $(CYAN)make install$(RESET)   - Install the binary locally"
	@echo -e "  $(CYAN)make version$(RESET)   - Show the version"
	@echo -e "  $(CYAN)make clean$(RESET)     - Remove binaries and logs"
	@echo -e "  $(CYAN)make build-prod$(RESET) - Compile the Go project for production"
	@echo -e "  $(CYAN)make watch$(RESET)     - Live reload"
	@echo -e "  $(CYAN)make help$(RESET)      - Show this message"

build:
	@echo -e "$(CYAN)[*] Building $(BOLD)...$(RESET)"
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(BINARY_NAME) $(CMD_DIR)/$(MAIN_FILE)
	@echo -e "$(GREEN)[+] Build complete!$(RESET)"

build-prod:
	@echo -e "$(CYAN)[*] Building $(BOLD)for production...$(RESET)"
	@mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=amd64 \
		go build \
		-trimpath \
		-ldflags="-s -w" \
		-o $(BIN_DIR)/$(BINARY_NAME) $(CMD_DIR)/$(MAIN_FILE)
	@echo -e "$(GREEN)[+] Production build complete!$(RESET)"

test:
	@go test ./...

install: build-prod
	@echo -e "$(CYAN)[*] Installing $(BOLD)...$(RESET)"
	@go install .
	@echo -e "$(GREEN)[+] Install complete!$(RESET)"

run: build
	@echo -e "$(CYAN)[*] Running $(BOLD)...$(RESET)"
	@$(BIN_DIR)/$(BINARY_NAME) $(ARGS)

run-prod: build-prod
	@echo -e "$(CYAN)[*] Running $(BOLD)...$(RESET)"
	@$(BIN_DIR)/$(BINARY_NAME) $(ARGS)

version:
	@echo -e "$(GREEN)[+] Version $(BINARY_NAME) v$(VERSION) $(RESET)"

clean:
	@echo -e "$(YELLOW)[*] Cleaning binaries...$(RESET)"
	@rm -rf $(BIN_DIR)/*
	@echo -e "$(GREEN)[+] Clean complete!$(RESET)"
	@echo -e "$(YELLOW)[*] Cleaning logs...$(RESET)"
	@rm -rf $(LOGS_DIR)/*
	@echo -e "$(GREEN)[+] Clean complete!$(RESET)"

watch:
	@if command -v air > /dev/null; then \
		air; \
		echo "Watching...";\
	else \
		read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/air-verse/air@latest; \
			air; \
			echo "Watching...";\
		else \
			echo "You chose not to install air. Exiting..."; \
			exit 1; \
		fi; \
	fi
