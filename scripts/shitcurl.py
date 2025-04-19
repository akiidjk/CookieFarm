#!../venv/bin/python3

import requests
import json
import logging


logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
BASE_URL = "http://localhost:8080"

def send_post_request(endpoint, headers=None, data=None, files=None):
    url = f"{BASE_URL}/{endpoint}"
    try:
        logging.info(f"Sending POST request to {url} with data: {data}")
        response = requests.post(url, headers=headers, data=data, files=files)
        response.raise_for_status()
        logging.info(f"Response received: {response.status_code}")
        return response
    except requests.exceptions.HTTPError as http_err:
        logging.error(f"HTTP error occurred: {http_err}")
    except Exception as err:
        logging.error(f"Other error occurred: {err}")
    return None

def login(password):
    logging.info("Logging in with provided password...")
    payload = {'password': password}
    headers = {'Content-Type': 'application/x-www-form-urlencoded'}
    response = send_post_request('api/v1/auth/login', headers=headers, data=payload)

    if response:
        try:
            token = response.json().get('token')
            if token:
                logging.info("Login successful, token received.")
                return token
            else:
                logging.warning("Token not found in the response.")
        except json.JSONDecodeError:
            logging.error("Failed to decode JSON response.")
    return None

def configure(token, config_data):
    logging.info("Configuring the system with provided configuration data...")
    headers = {
        'Authorization': f'Bearer {token}',
        'Content-Type': 'application/json',
    }
    payload = json.dumps(config_data)
    response = send_post_request('api/v1/config', headers=headers, data=payload)

    if response:
        logging.info("Configuration updated successfully.")
        return response.text
    return None

def verify_token(token):
    logging.info("Verifying token...")
    payload = {"token": token}
    response = send_post_request('api/v1/auth/verify', data=payload)

    if response:
        logging.info("Token verification completed.")
        return response.text
    return None

if __name__ == '__main__':
    password = 'password'
    token = login(password)

    if token:
        config_data = {
            "config": {
                "server": {
                    "host_flagchecker": "localhost:3000",
                    "team_token": "4242424242424242424",
                    "submit_flag_checker_time": 15,
                    "max_flag_batch_size": 1000,
                    "protocol": "cc_http"
                },
                "client": {
                    "base_url_server": "http://localhost:8080",
                    "submit_flag_server_time": 15,
                    "services": [
                        {"name": "CCApp", "port": 80},
                       	{"name": "Ticket", "port": 1337},
                       	{"name": "Poll", "port": 8080},
                       	{"name": "COOKIEFLAG", "port": 6969},
                    ],
                    "range_ip_teams": "80",
                    "format_ip_teams": "10.0.0.{}",
                    "my_team_ip": "10.0.0.1"
                }
            }
        }
        configure(token, config_data)
        verify_token(token)
