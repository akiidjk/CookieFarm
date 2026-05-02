#!/bin/bash

# Make sure we have directories ready
mkdir -p ../output

# CookieFarm (Polling SQLite database directly)
echo "Starting CookieFarm Timeline generation..."
python3 generate_flag_timeline.py \
    --db ../../../cookiefarm/cookiefarm.db \
    --query "SELECT COUNT(*) FROM flags" \
    --output ../output/cf_flag_count_timeline.txt \
    --duration 1200 &
CF_PID=$!

# DestructiveFarm (Polling SQLite database directly)
echo "Starting DestructiveFarm Timeline generation..."
python3 generate_flag_timeline.py \
    --db /tmp/DestructiveFarm/server/flags.sqlite \
    --query "SELECT COUNT(*) FROM flags" \
    --output ../output/df_flag_count_timeline.txt \
    --duration 1200 &
DF_PID=$!

echo "Both timelines are being generated in the background."
echo "CookieFarm PID: $CF_PID"
echo "DestructiveFarm PID: $DF_PID"

wait $CF_PID
wait $DF_PID
echo "Done."
