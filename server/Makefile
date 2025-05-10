BIN_DIR := ./bin
CMD_DIR := ./cmd/api
LOGS_DIR := ./logs
MAIN_FILE := main.go
BINARY_NAME := cookieserver
VERSION := 1.0.0
GOOS ?= linux
GOARCH ?= amd64


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
	@echo -e "  $(CYAN)make install$(RESET)   - Install the binary locally"
	@echo -e "  $(CYAN)make version$(RESET)   - Show the version"
	@echo -e "  $(CYAN)make clean$(RESET)     - Remove binaries and logs"
	@echo -e "  $(CYAN)make build-prod$(RESET) - Compile the Go project for production"
	@echo -e "  $(CYAN)make build-plugins$(RESET) - Compile the plugins"
	@echo -e "  $(CYAN)make build-plugins-prod$(RESET) - Compile the plugins for production"
	@echo -e "  $(CYAN)make watch$(RESET)     - Live reload"
	@echo -e "  $(CYAN)make tailwindcss-build$(RESET)     - Build Tailwind CSS minified"
	@echo -e "  $(CYAN)make tailwindcss-watch$(RESET)     - Tailwind CSS watch"
	@echo -e "  $(CYAN)make minify$(RESET)    - Minify the JS files"
	@echo -e "  $(CYAN)make lint$(RESET)      - Lint the Go files"
	@echo -e "  $(CYAN)make fmt$(RESET)       - Format the Go files"
	@echo -e "  $(CYAN)make help$(RESET)      - Show this message"

build:
	@echo -e "$(CYAN)[*] Building $(BOLD)...$(RESET)"
	@mkdir -p $(BIN_DIR)
	@go build -gcflags='github.com/ByteTheCookies/cookieserver/...="-m"' -o $(BIN_DIR)/$(BINARY_NAME) $(CMD_DIR)/$(MAIN_FILE)
	@echo -e "$(GREEN)[+] Build complete!$(RESET)"

build-plugins-prod:
	@echo -e "$(GREEN)[*] Building plugins in production for $(GOOS)/$(GOARCH)$(RESET)"
	@for file in $$(find ./protocols -name '*.go' ! -name 'protocols.go'); do \
		filename=$$(basename $$file); \
		pluginname=$${filename%.go}; \
		GOOS=$(GOOS) GOARCH=$(GOARCH) go build -trimpath -buildmode=plugin -ldflags="-s -w" -o "protocols/$$pluginname.so" "$$file"; \
		echo -e "$(GREEN)[+] Plugin $$pluginname.so built!$(RESET)"; \
	done

build-plugins:
	@for file in $$(find ./protocols -name '*.go' ! -name 'protocols.go'); do \
		filename=$$(basename $$file); \
		pluginname=$${filename%.go}; \
		go build -race -gcflags -m -buildmode=plugin -o "protocols/$$pluginname.so" "$$file"; \
		echo -e "$(GREEN)[+] Plugin $$pluginname.so built!$(RESET)"; \
	done

build-prod:
	@echo -e "$(CYAN)[*] Building $(BOLD)for production...$(RESET)"
	@mkdir -p $(BIN_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) \
		go build \
		-race   \
		-trimpath \
		-ldflags="-s -w" \
		-o $(BIN_DIR)/$(BINARY_NAME) $(CMD_DIR)/$(MAIN_FILE)
	@echo -e "$(GREEN)[+] Production build complete!$(RESET)"

install: tailwindcss-build build-prod build-plugins-prod
	@echo -e "$(CYAN)[*] Installing $(BOLD)...$(RESET)"
	@go install .
	@echo -e "$(GREEN)[+] Install complete!$(RESET)"

run: build build-plugins minify
	@echo -e "$(CYAN)[*] Running $(BOLD)...$(RESET)"
	@$(BIN_DIR)/$(BINARY_NAME) $(ARGS)

run-prod: build-prod tailwindcss-build minify
	@echo -e "$(CYAN)[*] Running $(BOLD)...$(RESET)"
	@$(BIN_DIR)/$(BINARY_NAME) $(ARGS)

version:
	@echo -e "$(GREEN)[+] Version $(BINARY_NAME) v$(VERSION) $(RESET)"

clean:
	@echo -e "$(YELLOW)[*] Cleaning binaries and db...$(RESET)"
	@rm -rf $(BIN_DIR)/*
	@echo -e "$(GREEN)[+] Clean complete!$(RESET)"
	@echo -e "$(YELLOW)[*] Cleaning logs...$(RESET)"
	@rm -rf $(LOGS_DIR)/*
	@echo -e "$(GREEN)[+] Clean complete!$(RESET)"

tailwindcss-build:
	./tools/tailwindcss -c ../tailwind.config.js -i ./assets/css/global.css -o ./public/css/output.css --minify

tailwindcss-watch:
	./tools/tailwindcss -c ../tailwind.config.js -i ./assets/css/global.css -o ./public/css/output.css --watch

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

minify:
	@echo -e "$(CYAN)[*] Minifying $(BOLD)...$(RESET)"
	@uglifyjs assets/js/*.js -o public/js/output.min.js -c -m
	@echo -e "$(GREEN)[+] Minification complete!$(RESET)"


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
