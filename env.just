#!/usr/bin/env just --justfile
# === COMMON VARIABLES ===

VERSION := `git describe --tags --abbrev=0`
RESET := "\\033[0m"

# === SERVER VARIABLES ===

SERVER_BIN_DIR := "../bin"
SERVER_CMD_DIR := "./server"
SERVER_LOGS_DIR := "./logs"
SERVER_MAIN_FILE := "main.go"
SERVER_BINARY_NAME := "cks"
GOOS := "linux"
GOARCH := "amd64"

# === CLIENT VARIABLES ===

CLIENT_BIN_DIR := "../bin"
CLIENT_CMD_DIR := "./client"
CLIENT_MAIN_FILE := "main.go"

# CROSS PLATFORM BUILD VARIABLES

PATHSEP := if os() == "windows" { "\\" } else { "/" }
MKDIR_CMD := if os() == "windows" { "mkdir" } else { "mkdir -p" }
ECHO_CMD := if os() == "windows" { "echo" } else { "echo -e" }
CLIENT_BINARY_NAME := if os() == "windows" { "ckc.exe" } else { "ckc" }
