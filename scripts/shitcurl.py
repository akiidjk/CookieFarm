import requests
import json

url = "http://localhost:8080/api/v1/auth/login"

payload = {'password': 'password'}
files=[

]
headers = {
  'Content-Type': 'application/x-www-form-urlencoded',
}

response = requests.request("POST", url, headers=headers, data=payload, files=files)

print(response.text)

token = response.json()['token']

url = "http://localhost:8080/api/v1/config"

payload = json.dumps({
  "config": {
    "server": {
      "host_flagchecker": "a",
      "team_token": "4242424242424242424",
      "cycle_time": 10
    },
    "client": {
      "base_url_server": "aaaaa",
      "cycle_time": 10,
      "services": None
    }
  }
})
headers = {
  'Authorization': f'Bearer {token}',
  'Content-Type': 'application/json',
}

response = requests.request("POST", url, headers=headers, data=payload)

print(response.text)
