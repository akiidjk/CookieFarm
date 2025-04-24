#!/bin/bash

if [[ $# -ne 2 ]]; then
    printf "Usage:\n  ./setup.sh <num_containers> <path_df>"
    exit
fi

cleanup() {
    echo "🧹 Pulizia in corso... Chiudo terminali e Docker..."
    kitty @ close-window --match title:flagchecker
    kitty @ close-window --match title:cookieserver
    kitty @ close-window --match title:service
    kitty @ close-window --match title:frontend
    cd ../tests/
    docker compose down
    cd ../scripts/
    exit
}

trap cleanup SIGINT

# Install requirements
pip install --upgrade pip > /dev/null
pip install -r ../requirements.txt > /dev/null

activate_venv="source ../venv/bin/activate"

# Run Flagchecker
echo "🚩 Avvio Flagchecker..."

chmod +x ../tests/flagchecker.py
kitty --title "flagchecker" bash -c "${activate_venv} && ../tests/flagchecker.py; exec bash" &

echo "✅ Flagchecker lanciato in un terminale separato! 🎉"
echo ""

# Run Services
echo "🚀 Avvio Servizi..."

cd ../tests
chmod +x ./start_containers.sh
kitty --title "service" bash -c "./start_containers.sh $1; exec bash" &

echo "🚀 Servizi avviati!"

# Run DestructiveFarm
echo "🚀 Avvio DestructiveFarm..."

cd ../scripts/
cat ./config_df.py > $2/server/config.py
chmod +x $2/server/start_server.sh
kitty --title "destructivefarm" bash -c "$2/server/start_server.sh; exec bash" &

echo "🚀 DestructiveFarm avviato!"

echo "🎯 Ambiente per DF pronto all'uso!"

# Attendi input per terminare tutti i terminali kitty
read -p "🔻 Premi INVIO per chiudere tutti i terminali avviati dallo script..."

cleanup()

echo "🧹 Tutti i terminali sono stati chiusi!"
