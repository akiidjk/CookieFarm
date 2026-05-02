#!/usr/bin/env python3

import json
import logging
import os

import requests

s = requests.Session()

logging.basicConfig(
    level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s"
)
BASE_URL = "http://localhost:8080"


def send_post_request(endpoint, headers=None, data=None, files=None):
    url = f"{BASE_URL}/{endpoint}"
    try:
        logging.info(f"Sending POST request to {url} with data: {data}")
        response = s.post(
            url, headers=headers, data=data, files=files, cookies=s.cookies
        )
        response.raise_for_status()
        logging.info(f"Response received: {response.status_code}")
        return response
    except requests.exceptions.HTTPError as http_err:
        logging.error(f"HTTP error occurred: {http_err}")
    except Exception as err:
        logging.error(f"Other error occurred: {err}")
    return None


def login(password) -> requests.Response:
    logging.info("Logging in with provided password...")
    payload = {"password": password}
    headers = {"Content-Type": "application/x-www-form-urlencoded"}
    response = send_post_request("api/v1/auth/login", headers=headers, data=payload)

    if response:
        if response.status_code == 200:
            logging.info("Login successful, token received.")
            return response
        else:
            logging.warning("Token not found in the response.")
            return None
    else:
        logging.error("Failed to send login request.")
        return None


def configure(config_data):
    logging.info("Configuring the system with provided configuration data...")
    headers = {
        "Content-Type": "application/json",
        "Cookie": f"token={s.cookies['token']}",
    }
    payload = json.dumps(config_data)
    response = send_post_request("api/v1/config", headers=headers, data=payload)

    if response:
        logging.info("Configuration updated successfully.")
        return response.json()
    return None


def upload_exploit(file_path):
    logging.info(f"Uploading exploit from {file_path}...")
    import os

    try:
        if not os.path.exists(file_path):
            logging.error(f"File does not exist: {file_path}")
            return None

        # Respect the Fiber endpoint's max file size: 10 MB
        max_size = 10 * 1024 * 1024
        file_size = os.path.getsize(file_path)
        if file_size > max_size:
            logging.error("File is too large: %d bytes (max %d)", file_size, max_size)
            return None

        filename = os.path.basename(file_path)
        if not filename:
            logging.error("Invalid file name for path: %s", file_path)
            return None

        # Open the file in a context manager so it is closed after the request
        with open(file_path, "rb") as f:
            # Let requests set the multipart Content-Type; provide the field name 'file'
            files = {"file": (filename, f)}

            # Include token as a Cookie header only if present in the session cookies.
            headers = {}
            token = s.cookies.get("token")
            if token:
                headers["Cookie"] = f"token={token}"

            response = send_post_request(
                "api/v1/exploit/upload", headers=headers, files=files
            )

            if response:
                logging.info("Exploit uploaded successfully.")
                try:
                    return response.json()
                except ValueError:
                    logging.error("Response is not valid JSON")
                    return None

    except Exception as err:
        logging.error(f"Error uploading exploit: {err}")
    return None


if __name__ == "__main__":
    password = "password"
    login(password)
    upload_exploit(f"{os.getenv('HOME')}/.config/cookiefarm/exploits/main.py")

    # config_data = json.load(open('config.json', 'r'))
    # config_data = {
    #     "config": config_data
    # }
    # configure(config_data)
