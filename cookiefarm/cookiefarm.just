#!/usr/bin/env just --justfile

import? '../env.just'

# === COMMON TARGETS ===

# Build both server and client
[group('build')]
build: server-build client-build
    @{{ ECHO_CMD }} "{{ GREEN }}[+] Build complete for both server and client!{{ RESET }}"

# === SERVER TARGETS ===

# Build server for development
[group('build')]
[working-directory('cookiefarm')]
server-build:
    @echo -e "{{ CYAN }}[*] Building server...{{ RESET }}"
    @{{ MKDIR_CMD }} {{ SERVER_BIN_DIR }}
    @go build  -race  -gcflags="-m" -o {{ SERVER_BIN_DIR }}{{ PATHSEP }}{{ SERVER_BINARY_NAME }} {{ SERVER_CMD_DIR }}{{ PATHSEP }}{{ SERVER_MAIN_FILE }}
    @echo -e "{{ GREEN }}[+] Server build complete!{{ RESET }}"

# Run the server in development mode
[group('build')]
[working-directory('cookiefarm')]
server-run: server-build server-build-plugins
    @{{ SERVER_BIN_DIR }}{{ PATHSEP }}{{ SERVER_BINARY_NAME }} -c -D

# Build server for production
[group('build')]
[working-directory('cookiefarm')]
server-build-prod:
    @echo -e "{{ CYAN }}[*] Building server for production...{{ RESET }}"
    @{{ MKDIR_CMD }} {{ SERVER_BIN_DIR }}
    @GOOS={{ GOOS }} GOARCH={{ GOARCH }} \
      go build -trimpath -ldflags="-s -w -X main.Version={{ VERSION }}" -o {{ SERVER_BIN_DIR }}{{ PATHSEP }}{{ SERVER_BINARY_NAME }} {{ SERVER_CMD_DIR }}{{ PATHSEP }}{{ SERVER_MAIN_FILE }}
    @echo -e "{{ GREEN }}[+] Production build complete!{{ RESET }}"

# Clean server binaries and logs
[group('dev')]
[working-directory('cookiefarm')]
server-clean:
    @rm -rf {{ SERVER_BIN_DIR }}{{ PATHSEP }}* {{ SERVER_LOGS_DIR }}{{ PATHSEP }}*

# Build server plugins
[group('build')]
[working-directory('cookiefarm')]
server-build-plugins:
    @for file in $(find ./pkg/protocols -name '*.go' ! -name 'protocols.go'); do \
     if grep -q '^package main' "$file"; then \
      filename=$(basename $file); \
      pluginname=${filename%.go}; \
      go build -race -gcflags="-m" -buildmode=plugin -o "./pkg/protocols/$pluginname.so" "$file"; \
     else \
      echo "Skipping $file: not a main package"; \
     fi; \
    done

# Build server plugins for production
[group('build')]
[working-directory('cookiefarm')]
server-build-plugins-prod:
    @for file in $(find ./pkg/protocols -name '*.go' ! -name 'protocols.go'); do \
     if grep -q '^package main' "$file"; then \
      filename=$(basename $file); \
      pluginname=${filename%.go}; \
      GOOS={{ GOOS }} GOARCH={{ GOARCH }} go build -trimpath  -ldflags="-s -w -X main.Version={{ VERSION }}" -buildmode=plugin -o "./pkg/protocols/$pluginname.so" "$file"; \
     else \
      echo "Skipping $file: not a main package"; \
     fi; \
    done

# Watch server files and rebuild on changes (requires air)
[group('dev')]
[working-directory('cookiefarm')]
server-watch:
    @if command -v air > /dev/null; then air; else go install github.com/air-verse/air@latest && air; fi

# === CLIENT TARGETS ===

# Build client for development
[group('build')]
[working-directory('cookiefarm')]
client-build:
    @{{ ECHO_CMD }} "{{ CYAN }}[*] Building client...{{ RESET }}"
    @{{ MKDIR_CMD }} {{ CLIENT_BIN_DIR }}
    @go build -race -ldflags="-X 'main.Version={{ VERSION }}'" -gcflags="-m" -o {{ CLIENT_BIN_DIR }}{{ PATHSEP }}{{ CLIENT_BINARY_NAME }} {{ CLIENT_CMD_DIR }}{{ PATHSEP }}{{ CLIENT_MAIN_FILE }}
    @{{ ECHO_CMD }} "{{ GREEN }}[+] Client build complete!{{ RESET }}"

# Build client for Windows
[group('build')]
[working-directory('cookiefarm')]
client-build-windows:
    @{{ ECHO_CMD }} "{{ CYAN }}[*] Building client for Windows...{{ RESET }}"
    @{{ MKDIR_CMD }} {{ CLIENT_BIN_DIR }}
    @GOOS=windows GOARCH=amd64 go build -ldflags="-X 'main.Version={{ VERSION }}'" -gcflags="-m" -o {{ CLIENT_BIN_DIR }}{{ PATHSEP }}{{ CLIENT_BINARY_NAME }} {{ CLIENT_CMD_DIR }}{{ PATHSEP }}{{ CLIENT_MAIN_FILE }}
    @{{ ECHO_CMD }} "{{ GREEN }}[+] Windows build complete!{{ RESET }}"

# Build client for Linux
[group('build')]
[working-directory('cookiefarm')]
client-build-linux:
    @{{ ECHO_CMD }} "{{ CYAN }}[*] Building client for Linux...{{ RESET }}"
    @{{ MKDIR_CMD }} {{ CLIENT_BIN_DIR }}
    @GOOS=linux GOARCH=amd64 go build -race -ldflags="-X 'main.Version={{ VERSION }}'" -gcflags="-m" -o {{ CLIENT_BIN_DIR }}{{ PATHSEP }}{{ CLIENT_BINARY_NAME }} {{ CLIENT_CMD_DIR }}{{ PATHSEP }}}{{ CLIENT_MAIN_FILE }}
    @{{ ECHO_CMD }} "{{ GREEN }}[+] Linux build complete!{{ RESET }}"

# Build client for production (Linux)
[group('build')]
[working-directory('cookiefarm')]
client-build-linux-prod:
    @{{ ECHO_CMD }} "{{ CYAN }}[*] Building client for Linux production...{{ RESET }}"
    @{{ MKDIR_CMD }} {{ CLIENT_BIN_DIR }}
    @GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X 'main.Version={{ VERSION }}'" -trimpath -o {{ CLIENT_BIN_DIR }}{{ PATHSEP }}{{ CLIENT_BINARY_NAME }} {{ CLIENT_CMD_DIR }}{{ PATHSEP }}{{ CLIENT_MAIN_FILE }}
    @{{ ECHO_CMD }} "{{ GREEN }}[+] Linux production build complete!{{ RESET }}"

# Build client for production (Windows)
[group('build')]
[working-directory('cookiefarm')]
client-build-windows-prod:
    @{{ ECHO_CMD }} "{{ CYAN }}[*] Building client for Windows production...{{ RESET }}"
    @{{ MKDIR_CMD }} {{ CLIENT_BIN_DIR }}
    @GOOS=windows GOARCH=amd64 go build -trimpath -gcflags="-m" -ldflags="-s -w -X 'main.Version={{ VERSION }}'" -o {{ CLIENT_BIN_DIR }}{{ PATHSEP }}{{ CLIENT_BINARY_NAME }} {{ CLIENT_CMD_DIR }}{{ PATHSEP }}{{ CLIENT_MAIN_FILE }}
    @{{ ECHO_CMD }} "{{ GREEN }}[+] Windows production build complete!{{ RESET }}"

# Build client for production (all platforms)
[group('build')]
[working-directory('cookiefarm')]
client-build-prod:
    @{{ ECHO_CMD }} "{{ CYAN }}[*] Building client for production...{{ RESET }}"
    @just client-build-linux-prod
    @just client-build-windows-prod
    @{{ ECHO_CMD }} "{{ GREEN }}[+] Production build complete!{{ RESET }}"

# Run the client
[group('dev')]
[working-directory('cookiefarm')]
client-run: client-build
    @{{ CLIENT_BIN_DIR }}{{ PATHSEP }}{{ CLIENT_BINARY_NAME }}

# Install the client binary to /usr/local/bin and the virtual environment
[group('dev')]
[working-directory('cookiefarm')]
client-install: client-build
    @sudo cp {{ CLIENT_BIN_DIR }}{{ PATHSEP }}{{ CLIENT_BINARY_NAME }} /usr/local/bin/{{ CLIENT_BINARY_NAME }}
    @sudo cp /usr/local/bin/{{ CLIENT_BINARY_NAME }} ~/.venv/bin/{{ CLIENT_BINARY_NAME }}

# Test
[group('test')]
[working-directory('cookiefarm')]
test:
    @mkdir -p ./coverage
    @gotestsum \
    --post-run-command "notify-send 'Test finished successfully' -a gotestsum -u normal" --format testname \
    -- work -coverprofile=./coverage/coverage.out -v ./... -count=10 \
    && go tool cover -html=./coverage/coverage.out -o ./coverage/coverage.html \
    && xdg-open ./coverage/coverage.html

# Start all the components for run mock tests mode for testing
[group('test')]
setup-tests num_containers="10" production_mode="false":
    cd ./demo/scripts && ./setup.sh {{ num_containers }} {{ production_mode }}

# Generate OpenAPI spec files from Swagger annotations
[group('tools')]
[working-directory('cookiefarm/server/api')]
swagger-generate:
    @echo -e "{{ CYAN }}[*] Generating Swagger/OpenAPI docs...{{ RESET }}"
    @if command -v swag > /dev/null; then \
      swag init -g swagger.go -o docs --parseDependency --parseInternal --outputTypes json,yaml; \
    else \
      echo "swag is not installed. Install it with: go install github.com/swaggo/swag/cmd/swag@latest"; \
      exit 1; \
    fi
    @echo -e "{{ GREEN }}[+] Generated docs/swagger.json and docs/swagger.yaml{{ RESET }}"
