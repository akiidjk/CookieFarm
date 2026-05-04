import argparse
import json
from pprint import pprint

import requests
from shitcurl import login


def test_pagination(base_url, limit, initial_cursor, token, query_params):
    url = f"{base_url}/api/v1/flags/{limit}"
    cookies = {"token": token}
    cursor = initial_cursor  # None = prima pagina

    page = 1
    while True:
        # ✅ params viene ricostruito ad ogni iterazione col cursore aggiornato
        params = {**{"cursor": cursor}, **query_params} if cursor else query_params
        print(params)
        print(f"\n--- Pagina {page} (cursor={cursor!r}) ---")
        try:
            response = requests.get(url, params=params, cookies=cookies)
            print(f"Status Code: {response.status_code}")
            body = response.json()
            pprint(body)

            cursor = body.get("next") or None  # ✅ aggiorna il cursore
            if not cursor:
                print("✅ Ultima pagina raggiunta.")
                break

            page += 1
            input("Press Enter per la prossima pagina...")

        except requests.exceptions.RequestException as e:
            print(f"Request failed: {e}")
            break


def test_get_flags(base_url, limit, query_params, token):
    # Construct the route: /api/v1/flags/:limit
    url = f"{base_url}/api/v1/flags/{limit}"

    # Optional: If the endpoint requires CookieAuth, define it here
    # You may need to replace 'session_token' with your actual cookie name/value
    cookies = {"token": token}

    # Filter out None values to only send provided parameters
    params = {k: v for k, v in query_params.items() if v is not None}

    print(f"Sending GET request to: {url}")
    print(f"With Query Parameters: {params}")

    try:
        response = requests.get(url, params=params, cookies=cookies)
        print(f"\nStatus Code: {response.status_code}")

        try:
            print("Response JSON:")
            pprint(response.json())
        except json.JSONDecodeError:
            print("Response is not JSON. Raw text:")
            print(response.text)

    except requests.exceptions.RequestException as e:
        print(f"Request failed: {e}")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Test the GET /api/v1/flags/:limit endpoint"
    )
    parser.add_argument(
        "--url", default="http://localhost:8080", help="Base URL of the API"
    )
    parser.add_argument("--limit", type=int, default=10, help="Page size (limit)")
    parser.add_argument(
        "--cursor", type=str, default=None, help="Cursore iniziale (base64)"
    )
    parser.add_argument("--status", type=int, help="Status filter (e.g., 0, 1, 2)")
    parser.add_argument("--service", type=str, help="Service filter name")
    parser.add_argument("--team", type=str, help="Team ID filter")
    parser.add_argument("--search", type=str, help="Search string")
    parser.add_argument(
        "--search-field",
        dest="search_field",
        type=str,
        help="Specific field to search in",
    )

    args = parser.parse_args()
    login_response = login("password")
    if not login_response or "token" not in login_response.cookies:
        print("Login failed or missing token cookie. Please check your credentials.")

    test_get_flags(
        args.url,
        args.limit,
        {
            "status": args.status,
            "service": args.service,
            "team": args.team,
            "search": args.search,
            "search_field": args.search_field,
            "cursor": args.cursor,
        },
        login_response.cookies["token"],
    )
    test_pagination(
        args.url,
        args.limit,
        args.cursor,
        login_response.cookies["token"],
        {
            "status": args.status,
            "service": args.service,
            "team": args.team,
            "search": args.search,
            "search_field": args.search_field,
            "cursor": args.cursor,
        },
    )
