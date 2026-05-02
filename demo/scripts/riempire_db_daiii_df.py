import base64
import json
import random
import string
import urllib.request
from urllib.error import HTTPError, URLError


def random_flag_code(length=32):
    return "".join(random.choices(string.ascii_uppercase + string.digits, k=length))


# Basic auth credentials
USERNAME = ""
PASSWORD = "password"

# Headers
headers = {"Content-Type": "application/json"}

# Simulation parameters
# Number of 120-second windows ("tickets") to emulate
tickets_to_emulate = 129
window_seconds = 120
min_flags_per_window = 99
max_flags_per_window = 456

services = ["http", "ssh", "dns", "smtp", "ftp", "redis", "mysql", "postgres"]

exploits = [
    "sqli_blind",
    "rce_template_injection",
    "path_traversal",
    "auth_bypass",
    "deserialization_rce",
    "command_injection",
    "ssrf",
    "buffer_overflow",
]

batch_size = 4_000
base_submit_time = random.randint(1_700_000_000, 1_750_000_000)

BASE_URL = "http://localhost:5000"


def send_post_request(endpoint, req_headers, body, username, password):
    url = f"{BASE_URL}/{endpoint}"
    credentials = base64.b64encode(f"{username}:{password}".encode()).decode("utf-8")
    all_headers = {**req_headers, "Authorization": f"Basic {credentials}"}
    data = body.encode("utf-8") if isinstance(body, str) else body
    req = urllib.request.Request(url, data=data, headers=all_headers, method="POST")
    try:
        with urllib.request.urlopen(req) as response:
            return response
    except HTTPError as e:
        print(f"HTTP error: {e.code} {e.reason}")
        return e
    except URLError as e:
        print(f"URL error: {e.reason}")
        return None


flags_batch = []
total_generated = 0

for ticket in range(tickets_to_emulate):
    submit_time = base_submit_time + (ticket * window_seconds)
    flags_in_window = random.randint(min_flags_per_window, max_flags_per_window)

    for _ in range(flags_in_window):
        service_name = random.choice(services)
        exploit_name = random.choice(exploits)

        flags_batch.append(
            {
                "flag": random_flag_code(50),
                "sploit": exploit_name,
                "team": f"team{random.randint(1, 80)}",
            }
        )
        total_generated += 1

        if len(flags_batch) >= batch_size:
            body = json.dumps(flags_batch)
            res = send_post_request("api/post_flags", headers, body, USERNAME, PASSWORD)
            if res and getattr(res, "status", getattr(res, "code", None)) == 200:
                print(f"Batch sent successfully! total_generated={total_generated}")
            else:
                print(
                    f"Failed to send batch: {getattr(res, 'status', getattr(res, 'code', 'No Response')) if res else 'No Response'}"
                )
            flags_batch = []

# Flush remaining flags
if flags_batch:
    body = json.dumps(flags_batch)
    res = send_post_request("api/post_flags", headers, body, USERNAME, PASSWORD)
    if res and getattr(res, "status", getattr(res, "code", None)) == 200:
        print(f"Final batch sent successfully! total_generated={total_generated}")
    else:
        print(
            f"Failed to send final batch: {getattr(res, 'status', getattr(res, 'code', 'No Response')) if res else 'No Response'}"
        )
