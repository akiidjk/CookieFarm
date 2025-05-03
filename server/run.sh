#!/bin/sh

# Default port fallback
PORT="${PORT:-8080}"
if [ -n "$CONFIG_FROM_FILE" ]; then
  exec /app/bin/cookieserver -p "$PASSWORD" -P "$PORT" -c "config.json"
else
  exec /app/bin/cookieserver -p "$PASSWORD" -P "$PORT"
fi
