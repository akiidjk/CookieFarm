// Package config provides functions to manage the CookieFarm client configuration globally.
package config

var (
	Token              string     // Token stores the global authentication token.
	ServerAddress      string     // HostServer holds the global base URL for the server.
	ArgsAttackInstance ArgsAttack // Struct holding runtime arguments
	ArgsConfigInstance ArgsConfig // Struct holding configuration arguments
	Current            Config     // Current holds the global configuration for the client.
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
nickname: "cookieguest"
`)
