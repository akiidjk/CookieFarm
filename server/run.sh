#!/bin/sh

# Default fallback per la porta
PORT="${PORT:-8080}"
DEBUG="${DEBUG:-false}"

# Costruisci il comando base
CMD="/app/bin/cookieserver -p \"$PASSWORD\" -P \"$PORT\""

# Aggiungi config file se esiste
if [ -n "$CONFIG_FROM_FILE" ]; then
    CMD="$CMD -c \"$CONFIG_FROM_FILE\""
fi

# Aggiungi flag debug solo se DEBUG=true
if [ "$DEBUG" = "true" ]; then
    CMD="$CMD -d"
fi

# Esegui
eval exec $CMD
