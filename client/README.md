# 📜 CookieFarm Exploitation Guide

Welcome to the CookieFarm client documentation! This guide explains how to create, run, and manage your exploits with CookieFarm's client and exploit manager.



# 👨‍💻 Client Overview

The **client** component of CookieFarm is responsible for:
- Retrieving flags through exploits
- Communicating with the server to submit obtained flags
- Handling multithreading and execution timing automatically

It is designed to let you **focus purely on writing exploits**, without worrying about concurrent execution, scheduling, or networking.



# 📊 Exploit Manager

The **exploit manager** is a Python utility that provides a decorator to simplify the exploitation process. It handles:

- Executing your exploit across all adversary machines
- Running exploits periodically (every "tick")
- Managing concurrency automatically with Python coroutines
- Proper output flushing and error handling

## 🔍 How it Works

Here is the basic structure of an exploit:

```python
#!/usr/bin/env python3

from utils.exploiter_manager import exploit_manager
import requests

@exploit_manager
def exploit(ip: str, port: int):
    r = requests.get(f"http://{ip}:{port}/get-flag")
    flag = r.text()
    return flag

if __name__ == "__main__":
    port = 4512  # Example port
    exploit(port=port)
```

Or, if you want to manually specify the server IP (optional advanced usage):

```python
#!/usr/bin/env python3

from utils.exploiter_manager import exploit_manager
import sys

@exploit_manager(server_ip=sys.argv[1])
def exploit(ip: str = "", port: int = 80):
    # Insert your exploit code here
    return  # Return flag

if __name__ == "__main__":
    port = 8081  # Adjust to target service port
    exploit(port=port)
```

✅ The `exploit_manager` automatically:
- Supplies IPs and ports
- Sends retrieved flags to the CookieFarm server
- Manages execution cycles and threads

> **Your only task: write the exploit logic!**

---

# 🚀 Running Your Exploit

Follow these steps to run your exploit with the client:

1. Navigate to the client directory:

   ```bash
   cd CookieFarm/client
   ```

2. Inside the `exploits/` folder, create your exploit script following the structure explained above.

3. Execute the following command:

   ```bash
   make run ARGS="--exploit ./<exploit_name>.py -b <server_url> -p <server_password> -t <tick_time>"
   ```

### 🔍 Command Arguments

| Argument | Description |
|:---------|:------------|
| `-e`, `--exploit` | Path to your exploit file (must be inside `exploits/` folder) |
| `-b`, `--base_url_server` | Base URL and port of the CookieFarm server |
| `-p`, `--password` | Password for server authentication |
| `-t`, `--tick` | Frequency in seconds to re-run the exploit and submit flags |

### 📂 Example Run

```bash
make run ARGS="--exploit ./my_exploit.py -b http://10.10.23.1:8080 -p Str0ng_p4ssw0rd -t 120"
```

This example runs `my_exploit.py` every 120 seconds, sending the obtained flags to `http://10.10.23.1:8080`, using the password `Str0ng_p4ssw0rd`.


