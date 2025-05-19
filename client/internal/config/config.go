// Package config provides functions to manage the CookieFarm client configuration globally.
package config

import "github.com/ByteTheCookies/cookieclient/internal/models"

var (
	Current    models.Config // Current holds the global configuration for the client.
	Token      string        // Token stores the global authentication token.
	HostServer string        // BaseURLServer holds the global base URL for the server.
	Args       models.Args   // Struct holding runtime arguments
)

var ExploitTemplate = []byte(`#!/usr/bin/env python3

from utils.exploiter_manager import exploit_manager

@exploit_manager
def exploit(ip:str, port:int, name: str):
    # Insert your exploit code here
    return []  # Return flag (or a list of flags, text containing the flags, HTML page containing the flags)


if __name__ == "__main__":
    exploit()
`)
