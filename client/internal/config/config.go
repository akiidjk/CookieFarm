// Package config provides functions to manage the CookieFarm client configuration globally.
package config

import "github.com/ByteTheCookies/cookieclient/internal/models"

var (
	Current    models.Config // Current holds the global configuration for the client.
	Token      string        // Token stores the global authentication token.
	HostServer string        // BaseURLServer holds the global base URL for the server.
	Args       models.Args   // Struct holding runtime arguments
)

const DefaultExploitPath = "~/.config/cookiefarm"

var ExploitTemplate = []byte(`#!/usr/bin/env python3
from cookiefarm_exploiter import exploiter_manager
import requests

@exploiter_manager
def exploit(ip, port, name_service):
    # Run your exploit here
    response = requests.get(f"http://{ip}:{port}/")

    # Just print the flag to stdout
    print(response.text)

# Run from the command line with arguments from CookieFarm
# python3 myexploit.py <ip_server> <password> <tick_time> <flag_regex> <thread_number> <port> <name_service>
`)
