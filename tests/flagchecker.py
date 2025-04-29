#!/usr/bin/env python3

from flask import Flask, request
import time
import random

app = Flask(__name__)

# Define constants for status codes
ACCEPTED = "flag claimed"
DENIED = ("invalid flag", "flag from nop team", "flag is your own", "flag too old","flag already claimed","the check which dispatched this flag didn't terminate successfully")
RESUBMIT = "the flag is not active yet, wait for next round"
ERROR = "notify the organizers and retry later"

status = {
    "ACCEPTED": ACCEPTED,
    "DENIED": DENIED,
    "RESUBMIT": RESUBMIT,
    "ERROR": ERROR,
}

@app.route("/submit", methods=['PUT'])
def check_flags():
    responses = []

    time.sleep(random.randint(0, 2))

    if request.headers.get('X-Team-Token'):
        data = request.get_json()
        if data:
            for flag in data:
                s = random.choice(list(status.keys()))
                if s == "DENIED":
                    message = random.choice(status["DENIED"])
                elif s == "RESUBMIT":
                    message = status["RESUBMIT"]
                elif s == "ERROR":
                    message = status["ERROR"]
                else:
                    message = status["ACCEPTED"]
                responses.append({
                    "msg": f"[{flag}] {message}",
                    "flag": flag,
                    "status": s
                })
    return responses

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5001)
