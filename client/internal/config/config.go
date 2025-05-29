// Package config provides functions to manage the CookieFarm client configuration globally.
package config

import "github.com/ByteTheCookies/cookieclient/internal/models"

var (
	Token         string            // Token stores the global authentication token.
	ServerAddress string            // HostServer holds the global base URL for the server.
	Nickname      string            // Nickname holds the global nickname for the client.
	Protocol      string            // Protocol holds the global protocol (e.g., http, https) for the server connection.
	ServerPort    uint16            // PortServer holds the global port for the server connection.
	ArgsAttack    models.ArgsAttack // Struct holding runtime arguments
	ArgsConfig    models.ArgsConfig // Struct holding configuration arguments
	Current       models.Config     // Current holds the global configuration for the client.
)

const DefaultConfigPath = "~/.config/cookiefarm"

var ExploitTemplate = []byte(`#!/usr/bin/env python3
from cookiefarm_exploiter import exploit_manager

@exploit_manager
def exploit(ip, port, name_service):
    # Run your exploit here
    flag = ""

    # Just print the flag to stdout
    print(flag)
`)

var ConfigTemplate = []byte(`address: "localhost"
port: 8080
https: false
nickname: "guest"
`)
