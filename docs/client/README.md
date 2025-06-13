# 📜 CookieFarm Client Exploitation Guide

## 🔧 Prerequisites

Before you begin, make sure you have:
- **Python 3+**
- **pip** for installing Python modules
- **Modern terminal** that supports ANSI colors and Unicode (for TUI mode)

---

## ⚙️ Client Installation

To install the CookieFarm client, run the following command:

```bash
bash <(curl -fsSL https://raw.githubusercontent.com/ByteTheCookies/CookieFarm/refs/heads/main/install.sh)
```

After installation, the `cookieclient` command will be globally available in your terminal. By default, it will launch in interactive TUI (Text User Interface) mode.

To uninstall the client:

```bash
bash <(curl -fsSL https://raw.githubusercontent.com/ByteTheCookies/CookieFarm/refs/heads/main/uninstall.sh)
```

---

## 📺 Interface Options

The CookieFarm client has two interface modes:

1. **Interactive TUI Mode** (default): A colorful, user-friendly interface with menus and keyboard navigation
   ```bash
   cookieclient
   # TUI starts automatically
   ```

2. **Traditional CLI Mode**: For scripts, automation, or environments where TUI isn't supported
   ```bash
   cookieclient --no-tui
   # OR
   COOKIECLIENT_NO_TUI=1 cookieclient
   ```

> In the rest of this guide, commands will be shown for both TUI and CLI modes. The `-N` flag is used to indicate that the command should run in CLI mode without TUI. If you set the environment variable `COOKIECLIENT_NO_TUI=1`, you can run commands without the `-N` flag.

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
   # In CLI mode (with no environment variable setted):
   cookieclient config login -P SuperSecret -N

   # In TUI mode:
   # Navigate to: Configuration → Login → Enter credentials
   ```

2. **Update configuration** with server details:
   ```bash
   # In CLI mode (with no environment variable setted):
   cookieclient config update -h 192.168.1.10 -p 8000 -u CookieMonster -N

   # In TUI mode:
   # Navigate to: Configuration → Update Config → Fill the form
   ```

3. **Install the helper Python module**:
   ```bash
   pip install cookiefarm-exploiter
   ```
  > For more information about the helper module, check the [cookiefarm-exploiter documentation](https://github.com/ByteTheCookies/CookieFarmExploiter)

4. **Create a new exploit template**:
   ```bash
   # In CLI mode (with no environment variable setted):
   cookieclient exploit create -n my_exploit -N

   # In TUI mode:
   # Navigate to: Exploits → Create Exploit → Enter name
   ```

5. **Write your exploit** in the created file `~/.config/cookiefarm/exploits/my_exploit.py`.

6. **Run the exploit**:
   ```bash
   # In CLI mode (with no environment variable setted):
   cookieclient exploit run -e my_exploit.py -p 1234 -t 120 -T 40 -N

   # In TUI mode:
   # Navigate to: Exploits → Run Exploit → Complete the form
   ```

---

## 🌟 TUI Navigation

The interactive TUI provides easy navigation with the following keyboard shortcuts:

| Key          | Action                         |
|-------------|--------------------------------|
| ↑/↓ or j/k  | Navigate menu items            |
| Enter       | Select item or submit form     |
| ESC         | Go back to previous screen     |
| Tab         | Navigate between input fields  |
| q or Ctrl+C | Quit the application           |

The TUI offers these main views:
- **Main Menu**: Choose between Configuration and Exploit operations
- **Config Menu**: Configuration management commands
- **Exploit Menu**: Exploit management commands
- **Input Forms**: Fill in required parameters for commands
- **Output View**: See command results with syntax highlighting

---

## Detailed Command Usage

### Config Login Command
Authenticate with the server using a password:
```bash
# In CLI mode (with no environment variable setted):
cookieclient config login -P <password> -N

# In TUI mode:
# Navigate to: Configuration → Login → Enter password
```
Parameters:
- `-P <password>`: The password for the server. This is required for authentication.

### Config Update Command
Update the client configuration (all fields optional, at least one required):
```bash
# In CLI mode (with no environment variable setted):
cookieclient config update -h <server_ip> -p <port> -u <username> [-s] -N

# In TUI mode:
# Navigate to: Configuration → Update Config → Fill the form
```
Parameters:
- `-h <server_ip>`: IP address of the server.
- `-p <port>`: Port of the server.
- `-u <username>`: Username for the client. (default is `guest`)
- `-s`: Use secure connection (HTTPS).

### Config Show Command
Display the current configuration:
```bash
# In CLI mode (with no environment variable setted):
cookieclient config show -N

# In TUI mode:
# Navigate to: Configuration → Show Config
```

### Config Reset Command
Reset the configuration to default:
```bash
# In CLI mode (with no environment variable setted):
cookieclient config reset -N

# In TUI mode:
# Navigate to: Configuration → Reset Config
```

### Config Logout Command
Log out and clear the current session:
```bash
# In CLI mode (with no environment variable setted):
cookieclient config logout -N

# In TUI mode:
# Navigate to: Configuration → Logout
```

---

### Exploit Create Command
Create a new exploit template:
```bash
# In CLI mode (with no environment variable setted):
cookieclient exploit create -n <exploit_name> -N

# In TUI mode:
# Navigate to: Exploits → Create Exploit → Enter name
```
Parameters:
- `-n <exploit_name>`: Name of the exploit template. This can be a path to a file or just a name.
> If user does not specify the path in the exploit name, it will be created in the `~/.config/cookiefarm/exploits/` directory. Otherwise, it will be created in the specified path.

*Example:*
```bash
# In CLI mode (with no environment variable setted):
cookieclient exploit create -n ./my_exploit -N

# In TUI mode:
# Navigate to: Exploits → Create Exploit → Enter name as `./my_exploit`
```
In this case, the exploit will be created in the current directory.

### Exploit Run Command
Run an exploit:
```bash
# In CLI mode (with no environment variable setted):
cookieclient exploit run -e <exploit_file> -p <port> [-t <timeout>] [-T <threads>] [-d] -N

# In TUI mode:
# Navigate to: Exploits → Run Exploit → Fill the form
```
Parameters:
- `-e <exploit_file>`: Path to the exploit file (Python script).
- `-p <port>`: Port to run the exploit on.
- `-t <timeout>`: Timeout for the exploit in seconds (default is 120).
- `-T <threads>`: Number of threads to use (default is 10).
- `-d`: Enable debug mode for more verbose output.

> [!IMPORTANT]
> When you run an exploit in TUI mode, it will run in the background, allowing you to continue using the client while monitoring the exploit's progress.
*Example:*
```bash
# In CLI mode (with no environment variable setted):
cookieclient exploit run -e my_exploit.py -p 1234 -t 120 -T 40 -N
# This will return the PID of the running exploit.
#
# In TUI mode:
# # Navigate to: Exploits → Run Exploit → Enter the exploit file, port, timeout, and threads
# # The exploit will run in the background, and you can monitor its progress.
```

### Exploit List Command
List all running exploits:
```bash
# In CLI mode (with no environment variable setted):
cookieclient exploit list -N

# In TUI mode:
# Navigate to: Exploits → List Running Exploits
```

### Exploit Remove Command
Remove a saved exploit template:
```bash
# In CLI mode (with no environment variable setted):
cookieclient exploit remove -n <exploit_name> -N

# In TUI mode:
# Navigate to: Exploits → Remove Exploit → Enter exploit name
```
Parameters:
- `-n <exploit_name>`: Name of the exploit template to remove. This can be a path to a file or just a name.

### Exploit Stop Command
Stop a running exploit:
```bash
# In CLI mode (with no environment variable setted):
cookieclient exploit stop -p <pid> -N

# In TUI mode:
# Navigate to: Exploits → Stop Exploit → Enter PID
```
Parameters:
- `-p <pid>`: Process ID of the exploit to stop.
