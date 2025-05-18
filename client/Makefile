BIN_DIR := ./bin
CMD_DIR := ./
LOGS_DIR := ./logs
MAIN_FILE := ./main.go
BINARY_NAME := ./cookieclient
VERSION := 1.1.0

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
	@echo -e "  $(CYAN)make test$(RESET)      - Run tests"
	@echo -e "  $(CYAN)make install$(RESET)   - Install the binary locally"
	@echo -e "  $(CYAN)make version$(RESET)   - Show the version"
	@echo -e "  $(CYAN)make clean$(RESET)     - Remove binaries and logs"
	@echo -e "  $(CYAN)make build-prod$(RESET) - Compile the Go project for production"
	@echo -e "  $(CYAN)make lint$(RESET)      - Lint the Go files"
	@echo -e "  $(CYAN)make fmt$(RESET)       - Format the Go files"
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
		-race   \
		-trimpath \
		-ldflags="-s -w" \
		-o $(BIN_DIR)/$(BINARY_NAME) $(CMD_DIR)/$(MAIN_FILE)
	@echo -e "$(GREEN)[+] Production build complete!$(RESET)"

test:
	@go test ./...

install: build
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
	@echo -e "$(YELLOW)[*] Cleaning binaries and db...$(RESET)"
	@if [ -d $(BIN_DIR) ]; then rm -rf $(BIN_DIR)/*; fi
	@echo -e "$(GREEN)[+] Clean complete!$(RESET)"
	@echo -e "$(YELLOW)[*] Cleaning logs...$(RESET)"
	@if [ -d $(LOGS_DIR) ]; then rm -rf $(LOGS_DIR)/*; fi
	@echo -e "$(GREEN)[+] Clean complete!$(RESET)"

lint:
	@echo -e "$(CYAN)[*] Linting $(BOLD)...$(RESET)"
	@if ! golangci-lint run; then \
		echo -e "$(RED)[!] Linting failed!$(RESET)"; \
		exit 1; \
	fi
	@echo -e "$(GREEN)[+] Linting complete!$(RESET)"

fmt: ## run gofumpt (if present) or gofmt
	@if command -v gofumpt > /dev/null; then \
		if ! gofumpt -w -d .; then \
			echo -e "$(RED)[!] Formatting with gofumpt failed!$(RESET)"; \
			exit 1; \
		fi; \
	else \
		if ! go list -f {{.Dir}} ./... | xargs gofmt -w -s -d; then \
			echo -e "$(RED)[!] Formatting with gofmt failed!$(RESET)"; \
			exit 1; \
		fi; \
	fi
	@echo -e "$(GREEN)[+] Formatting complete!$(RESET)"
