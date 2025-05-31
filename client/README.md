# 📜 CookieFarm Client Exploitation Guide

## 🔧 Prerequisites

Before you begin, make sure you have:
- **Python 3+**
- **pip** for installing Python modules

---

## ⚙️ Client Installation

To install the CookieFarm client, run the following command:

```bash
bash <(curl -fsSL https://raw.githubusercontent.com/ByteTheCookies/CookieFarm/refs/heads/main/install.sh)
```

After installation, the `cookieclient` command will be globally available in your terminal.

To uninstall the client:

```bash
bash <(curl -fsSL https://raw.githubusercontent.com/ByteTheCookies/CookieFarm/refs/heads/main/uninstall.sh)
```

---

## 🚀 Client Command Overview

### Configuration Commands

| Command                     | Description                                     | Detailed Usage |
|----------------------------|-------------------------------------------------| -------------------------------------------------|
| `config login`             | Log in to the server with a password            | [Section](#config-login-command) |
| `config update`            | Update client config with IP, port, user, etc.  | [Section](#config-update-command) |
| `config show`              | Show current configuration                      | [Section](#config-show-command) |
| `config reset`             | Reset configuration to default                  | [Section](#config-reset-command) |
| `config logout`            | Log out from the server                         | [Section](#config-logout-command) |

### Exploit Management Commands

| Command                      | Description                                        | Detailed Usage |
|-----------------------------|----------------------------------------------------| -------------------------------------------------|
| `exploit create`            | Create a new exploit template                      | [Section](#exploit-create-command) |
| `exploit run`               | Run an exploit against a target                   | [Section](#exploit-run-command) |
| `exploit list`              | List all currently running exploits                | [Section](#exploit-list-command) |
| `exploit remove`            | Remove an exploit template                         | [Section](#exploit-remove-command) |
| `exploit stop`              | Stop a running exploit                             | [Section](#exploit-stop-command) |

---

## 🧪 Exploitation Workflow

1. **Log in** to the server:
   ```bash
   cookieclient config login -P SuperSecret
   ```

2. **Update configuration** with server details:
   ```bash
   cookieclient config update -h 192.168.1.10 -p 8000 -u CookieMonster
   ```

3. **Install the helper Python module**:
   ```bash
   pip install cookiefarm-exploiter
   ```

4. **Create a new exploit template**:
   ```bash
   cookieclient exploit create -n my_exploit
   ```

5. **Run the exploit**:
   ```bash
   cookieclient exploit run -e my_exploit.py -p 1234 -t 120 -T 40
   ```

---

## Detailed Command Usage

### Config Login Command
Authenticate with the server using a password:
```bash
cookieclient config login -P <password>
```
Parameters:
- `-P <password>`: The password for the server. This is required for authentication.

### Config Update Command
Update the client configuration (all fields optional, at least one required):
```bash
cookieclient config update -h <server_ip> -p <port> -u <username> [-s]
```
Parameters:
- `-h <server_ip>`: IP address of the server.
- `-p <port>`: Port of the server.
- `-u <username>`: Username for the client. (default is `guest`)
- `-s`: Use secure connection (HTTPS).

### Config Show Command
Display the current configuration:
```bash
cookieclient config show
```

### Config Reset Command
Reset the configuration to default:
```bash
cookieclient config reset
```

### Config Logout Command
Log out and clear the current session:
```bash
cookieclient config logout
```

---

### Exploit Create Command
Create a new exploit template:
```bash
cookieclient exploit create -n <exploit_name>
```
Parameters:
- `-n <exploit_name>`: Name of the exploit template. This can be a path to a file or just a name.
> If user does not specify the path in the exploit name, it will be created in the `~/.cookieclient/exploits/` directory. Otherwise, it will be created in the specified path.

Example:
```bash
cookieclient exploit create -n ./my_exploit
```
In this case, the exploit will be created in the current directory.

### Exploit Run Command
Run an exploit:
```bash
cookieclient exploit run -e <exploit_file> -p <port> [-t <timeout>] [-T <threads>] [-d]
```
Parameters:
- `-e <exploit_file>`: Path to the exploit file (Python script).
- `-p <port>`: Port to run the exploit on.
- `-t <timeout>`: Timeout for the exploit in seconds (default is 120).
- `-T <threads>`: Number of threads to use (default is 10).
- `-d`: Enable debug mode for more verbose output.

Example:
```bash
cookieclient exploit run -e my_exploit.py -p 1234 -t 120 -T 40
# This will return the PID of the running exploit.
```

### Exploit List Command
List all running exploits:
```bash
cookieclient exploit list
```

### Exploit Remove Command
Remove a saved exploit template:
```bash
cookieclient exploit remove -n <exploit_name>
```
Parameters:
- `-n <exploit_name>`: Name of the exploit template to remove. This can be a path to a file or just a name.

### Exploit Stop Command
Stop a running exploit:
```bash
cookieclient exploit stop -p <pid>
```
Parameters:
- `-p <pid>`: Process ID of the exploit to stop.
