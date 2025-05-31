<div align="center">
  <img width="360px" height="auto" src="assets/logo_mucca.png" alt="CookieFarm Logo">
</div>

<p align="center">
  <img src="https://img.shields.io/badge/version-1.0.0-blue" alt="Version">
  <img src="https://img.shields.io/badge/languages-Go%20%7C%20Python-yellowgreen" alt="Languages">
  <img src="https://img.shields.io/badge/keywords-CTF%2C%20Exploiting%2C%20Attack%20Defense-red" alt="Keywords">
</p>

# 🍪 CookieFarm

**CookieFarm** is an *Attack/Defense CTF* framework inspired by DestructiveFarm, developed by the Italian team **ByteTheCookies**.
Its strength lies in a hybrid **Go + Python** architecture and a **zero-distraction philosophy**:
> 🎯 *Your only task is to write the exploit!*

CookieFarm automates exploit distribution, flag submission, and result monitoring — allowing you to focus entirely on building powerful exploits.

---

## 🔧 Prerequisites

Make sure you have the following installed:

- ✅ Python 3+
- ✅ Docker

---

## 📁 Repository Structure

| Directory       | Description |
|------------------|-------------|
| [`client/`](./client/) | Handles exploit creation and flag submission |
| [`server/`](./server/) | Manages exploit distribution, flag collection, and monitoring |

---

## ⚙️ Architecture Overview

<div align="center">
  <img width="800px" height="auto" src="assets/arch_farm.png" alt="Architecture Diagram">
</div>

---

## ▶️ Getting Started

### 🖥️ Starting the Server

1. Move into the `server/` directory:
   ```bash
   cd server/
   ```

2. Create an `.env` file in the server directory to configure the environment settings:

    ```bash
      # Server configuration
      DEBUG=false                   # Enable debug mode for verbose logging
      SERVER_PASSWORD=SuperSecret  # Set a strong password for authentication
      CONFIG_FROM_FILE=config.yml  # Set if the server takes the config from config.yml in the filesystem; otherwise, do not set the variable
      SERVER_PORT=8080            # Define the port the server will listen on
    ```

  > ⚠️ For production environments, set `DEBUG=false` and use a strong, unique password

3. Launch the server using Docker Compose:

   ```bash
   docker compose up
   ```

   > ⚠️ In production, keep `DEBUG=false` and set a strong, unique password.

3. Start the server with Docker Compose:
   ```bash
   docker compose up --build
   ```

📘 For more configuration details, refer to the [server documentation](./server/README.md).

---

### 💻 Using the Client & Running Exploits

1. Run the installation script:
   ```bash
   bash <(curl -fsSL https://raw.githubusercontent.com/ByteTheCookies/CookieFarm/refs/heads/main/install.sh)
   ```

   > After installation, the `cookieclient` command is globally accessible.

2. Log in and configure the client:
   ```bash
   cookieclient config login -P SuperSecret
   cookieclient config update -h 192.168.1.10 -p 8000 -u your_username
   ```

3. Install the Python helper module and create a new exploit template:
   ```bash
   pip install cookiefarm-exploiter
   cookieclient create -n your_exploit_name
   ```

   This will generate `your_exploit_name.py` in `~/.cookiefarm/exploits/`.

4. Run your exploit:
   ```bash
   cookieclient attack -e your_exploit_name.py -p 1234 -t 120 -T 40
   ```

📘 For more usage examples, check out the [client documentation](./client/README.md).

---

## 🤝 Contributing

We welcome contributions, suggestions, and bug reports!
See [CONTRIBUTING.md](./CONTRIBUTING.md) for details on how to get involved.

---

<div align="center">
  Built with ❤️ by <strong>ByteTheCookies</strong>
</div>
