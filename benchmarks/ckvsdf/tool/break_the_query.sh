#!/bin/bash

echo "Filling the dbs with the data for the query..."

CF_DB="../../../cookiefarm/cookiefarm.db"
DF_DB="/tmp/DestructiveFarm/server/flags.sqlite"
CF_OUT="/dev/null"
DF_OUT="/dev/null"
QUERY="SELECT COUNT(*) FROM flags"
DURATION="${1:-1200}"

python3 generate_flag_timeline.py \
    --source "CF:${CF_DB}:${QUERY}:${CF_OUT}" \
    --source "DF:${DF_DB}:${QUERY}:${DF_OUT}" \
    --interval 0.5 \
    --duration $DURATION &

for i in {1..10}; do
    python3 ../../../demo/scripts/riempire_db_daiii.py
    python3 ../../../demo/scripts/riempire_db_daiii_df.py
    echo "Run $i completed."
done

pkill -f "python3 generate_flag_timeline.py"

./measure_pagination.sh ../output/cf_latency_cold_breakQUERY.json ../output/cf_latency_warm_breakQUERY.json ../output/df_latency_cold_breakQUERY.json ../output/df_latency_warm_breakQUERY.json
