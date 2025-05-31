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

| Command                     | Description                                     |
|----------------------------|-------------------------------------------------|
| `config login`             | Log in to the server with a password            |
| `config update`            | Update client config with IP, port, user, etc.  |
| `config show`              | Show current configuration                      |
| `config reset`             | Reset configuration to default                  |
| `config logout`            | Log out from the server                         |

### Exploit Management Commands

| Command                      | Description                                        |
|-----------------------------|----------------------------------------------------|
| `exploit create`            | Create a new exploit template                      |
| `exploit run`               | Run an exploit against a target                   |
| `exploit list`              | List all currently running exploits                |
| `exploit remove`            | Remove an exploit template                         |
| `exploit stop`              | Stop a running exploit                             |

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

## 🧾 Detailed Command Usage

### 🔐 `config login`
Authenticate with the server using a password:
```bash
cookieclient config login -P <password>
```

### 🛠️ `config update`
Update the client configuration (all fields optional, at least one required):
```bash
cookieclient config update -h <server_ip> -p <port> -u <username> [-s]
```

### 📋 `config show`
Display the current configuration:
```bash
cookieclient config show
```

### ♻️ `config reset`
Reset the configuration to default:
```bash
cookieclient config reset
```

### 🚪 `config logout`
Log out and clear the current session:
```bash
cookieclient config logout
```

---

### 🧱 `exploit create`
Create a new exploit template:
```bash
cookieclient exploit create -n <exploit_name>
```

Example:
```bash
cookieclient exploit create -n ./my_exploit
```

### 🎯 `exploit run`
Run an exploit:
```bash
cookieclient exploit run -e <exploit_file> -p <port> [-t <timeout>] [-T <threads>] [-d]
```

### 📄 `exploit list`
List all running exploits:
```bash
cookieclient exploit list
```

### 🗑️ `exploit remove`
Remove a saved exploit template:
```bash
cookieclient exploit remove -n <exploit_name>
```

### ⛔ `exploit stop`
Stop a running exploit:
```bash
cookieclient exploit stop -n <exploit_name>
```
