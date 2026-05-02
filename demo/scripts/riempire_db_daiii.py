import base64
import json
import random

from shitcurl import login, send_post_request


def random_flag_code(length=32):
    random_bytes = bytes(random.randint(0, 255) for _ in range(length))
    return base64.b64encode(random_bytes).decode("utf-8")[:length].upper()


# login
s = login("password")
if not s or "token" not in s.cookies:
    raise RuntimeError("Login failed or missing token cookie")

headers = {"Content-Type": "application/json", "Cookie": f"token={s.cookies['token']}"}

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
                "status": random.randint(0, 2),
                "team_id": random.randint(1, 80),
                "port_service": random.randint(1, 65535),
                "service_name": service_name,
                "flag_code": random_flag_code(50),
                "response_time": submit_time + random.randint(1, 8),
                "submit_time": submit_time,
                "msg": random.choice(
                    ["", "queued", "submitted", "invalid", "duplicate"]
                ),
                "username": f"team{random.randint(1, 80)}",
                "exploit_name": exploit_name,
            }
        )
        total_generated += 1

        if len(flags_batch) >= batch_size:
            body = json.dumps({"flags": flags_batch})
            res = send_post_request("api/v1/submit-flags", headers, body)
            if res and res.status_code == 200:
                print(f"Batch sent successfully! total_generated={total_generated}")
            else:
                print(
                    f"Failed to send batch: {res.status_code if res else 'No Response'}"
                )
            flags_batch = []

# Flush remaining flags
if flags_batch:
    body = json.dumps({"flags": flags_batch})
    res = send_post_request("api/v1/submit-flags", headers, body)
    if res and res.status_code == 200:
        print(f"Final batch sent successfully! total_generated={total_generated}")
    else:
        print(
            f"Failed to send final batch: {res.status_code if res else 'No Response'}"
        )
