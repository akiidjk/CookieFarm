# 📜 CookieFarm Client Exploitation Guide

## Prerequisites


# Installation Clientqua

For install the CookieFarm client, run the following command in your terminal:

```bash
 bash <(curl -fsSL https://raw.githubusercontent.com/ByteTheCookies/CookieFarm/refs/heads/main/install.sh)
```

Now you can use the `cookieclient` command globally in your terminal.
If you want to uninstall the client, you can run the following command:

```bash
 bash <(curl -fsSL https://raw.githubusercontent.com/ByteTheCookies/CookieFarm/refs/heads/main/uninstall.sh)
```

---

# User Guide Client

## Configuration Client

Executing the command `cookieclient config` will display the available options for configuring the client:

| Command | Description | Detail Section |
|---------|-------------| -------------|
| `cookieclient config login` | Log in to the server with a password | [Login Section](#login-command) |
| `cookieclient config update` | Update the client configuration with server details | [Update Section](#update-command) |
| `cookieclient config show` | Show the current client configuration | [Show Section](#show-command) |
| `cookieclient config reset` | Reset the client configuration to default | [Reset Section](#reset-command) |
| `cookieclient config logout` | Log out from the server | [Logout Section](#logout-command) |

There other commands available 

## Commands Overview

### Login Command
This command allows you to log in to the CookieFarm server using a password.

```bash
cookieclient config login -P <password>
```
_Example:_

```bash
cookieclient config login -P SuperSecret
```

With this command, you login to the server with the password `SuperSecret`.


### Update Command
This command updates the client configuration with the server's IP address, port, username and if you want to use https.
Every fiels is optional, but you must specify at least one of them.

```bash
cookieclient config update -h <server_ip> -p <server_port> -u <username> [-s]
```

_Example:_

```bash
cookieclient config update -h 192.168.1.10 -p 8000 -u CookieMonster
```

This command updates the client configuration to connect to the server at `192.168.1.10` on port `8000` with the username `CookieMonster`.

### Show Command
This command displays the current client configuration, including server IP, port, username, and whether HTTPS is enabled.

```bash
cookieclient config show
```

### Reset Command
This command resets the client configuration to its default state, removing any custom settings.

```bash
cookieclient config reset
```

### Logout Command
This command logs out from the server, clearing the current session.

```bash
cookieclient config logout
```
