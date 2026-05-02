#!/bin/bash

# Make sure we have directories ready
mkdir -p ../output

CF_DB="../../../cookiefarm/cookiefarm.db"
DF_DB="/tmp/DestructiveFarm/server/flags.sqlite"
CF_OUT="../output/cf_flag_count_timeline.txt"
DF_OUT="../output/df_flag_count_timeline.txt"
QUERY="SELECT COUNT(*) FROM flags"
DURATION="${1:-1200}"

echo "==> Starting synchronized flag count timeline generation..."
echo "    CF DB: $CF_DB"
echo "    DF DB: $DF_DB"
echo "    Duration: $DURATION seconds"

python3 generate_flag_timeline.py \
    --source "CF:${CF_DB}:${QUERY}:${CF_OUT}" \
    --source "DF:${DF_DB}:${QUERY}:${DF_OUT}" \
    --interval 0.5 \
    --duration $DURATION

echo "Done."
