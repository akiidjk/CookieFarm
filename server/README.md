# 🛠️ CookieFarm Server Guide

Welcome to the documentation of the **CookieFarm server**!  
This guide describes how to launch, configure, and use the server responsible for flag management.

---

## 🌐 Server Overview

The **CookieFarm server** performs the following tasks:

- Receives and stores flags submitted by clients into a database.
- Sends flags to the `flagchecker` and handles the results.
- Displays all flags via a web interface:
  - Whether they have been submitted.
  - Whether they have been accepted.

The server is written in **Go** and is designed for easy deployment in both development and production environments.

---

## 🚀 Running the Server

To run the server use docker with:
```bash
docker compose up
```


---

## ⚙️ Execution Options

When running the binary (or using `make run`), you can specify the following **optional arguments**:

| Flag              | Description                                                                      | Default      |
|-------------------|----------------------------------------------------------------------------------|--------------|
| `-d`, `--debug`   | Enables debug mode                                                               | `false`      |
| `-c`, `--config`  | Path to a JSON config file (instead of using the web form)                       | N/A          |
| `-p`, `--password`| Password to access the server web interface                                      | `"password"` |
| `-P`, `--port`    | Sets the port the server will listen on                                          | `8080`       |

---

## 🌐 Web Interface

The **web interface** is accessible at:

```
http://<your_server_ip>:<port>
```

Through the UI you can:

- View all received flags.
- Check the submission and acceptance status of each flag.
- Configure the server (unless `--config` mode is enabled).

🔐 **Password-protected access**: set the password via `-p` or use the default `"password"`.

---

## 📂 Example Usage

```bash
./bin/server -d -p SuperSecret -P 9090
```

This command runs the server:

- In debug mode.
- With password `SuperSecret`.
- On port `9090`.

---

## 📝 Notes

- The server is fully compatible with CookieFarm clients.
- Ensure the `flagchecker` configuration is properly set via JSON or the startup form.
- You may use `.env` files or service managers (e.g., `systemd`, `docker`) for advanced deployment needs.

---

Happy flag handling! 🎯
