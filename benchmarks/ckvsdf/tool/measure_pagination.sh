#!/bin/bash

# ==========================================
# CookieFarm Pagination Measurement
# ==========================================

echo "==> Authenticating to CookieFarm..."
CF_TOKEN=$(curl -s -i -X POST http://localhost:8080/api/v1/auth/login -d "password=password" | grep -i "set-cookie: token=" | sed -e 's/.*token=\([^;]*\).*/\1/' | tr -d '\r')

if [ -z "$CF_TOKEN" ]; then
    echo "Failed to get CookieFarm token"
    # exit 1
fi

echo "==> Running CookieFarm Cold Cache test..."
python3 measure_pagination.py \
    --url "http://localhost:8080/api/v1/flags/40" \
    --requests 50 \
    --mode cold \
    --cookie "token=${CF_TOKEN}" \
    --output ../output/cf_latency_cold.json

echo "==> Running CookieFarm Warm Cache test..."
python3 measure_pagination.py \
    --url "http://localhost:8080/api/v1/flags/40" \
    --requests 50 \
    --warmup 10 \
    --mode warm \
    --cookie "token=${CF_TOKEN}" \
    --output ../output/cf_latency_warm.json

# ==========================================
# DestructiveFarm Pagination Measurement
# ==========================================

# Note: DestructiveFarm API endpoint for flags might differ, e.g. /api/get_flags
# Update the URL to match DF's flag endpoint. Here using a placeholder /api/flags/40
DF_URL="http://localhost:5000/ui/show_flags"
DF_DATA="sploit=&team=&flag=&time-since=&time-until=&status=&checksystem_response=&page-number=9838"

echo "==> Running DestructiveFarm Cold Cache test..."
python3 measure_pagination.py \
    --url "$DF_URL" \
    --requests 50 \
    --mode cold \
    --method POST \
    --data "$DF_DATA" \
    --basic-auth ":password" \
    --output ../output/df_latency_cold.json

echo "==> Running DestructiveFarm Warm Cache test..."
python3 measure_pagination.py \
    --url "$DF_URL" \
    --requests 50 \
    --warmup 10 \
    --mode warm \
    --method POST \
    --data "$DF_DATA" \
    --basic-auth ":password" \
    --output ../output/df_latency_warm.json
