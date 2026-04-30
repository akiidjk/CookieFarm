CONFIG = {
    "TEAMS": {"Team #{}".format(i): "10.10.{}.1".format(i) for i in range(0, 39 + 1)},
    "FLAG_FORMAT": r"[A-Z0-9]{31}=",
    "SYSTEM_PROTOCOL": "ructf_http",
    "SYSTEM_URL": "http://localhost:5001/submit",
    "SYSTEM_TOKEN": "password",
    "SUBMIT_FLAG_LIMIT": 100,
    "SUBMIT_PERIOD": 5,
    "FLAG_LIFETIME": 5 * 60,
    "SERVER_PASSWORD": "password",
    "ENABLE_API_AUTH": False,
    "API_TOKEN": "00000000000000000000",
}
