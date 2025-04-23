#!/bin/bash

if [[ $# -ne 1 ]]; then
    printf "Usage:\n  ./setup.sh <num_containers>"
    exit
fi

cleanup() {
    echo "🧹 Pulizia in corso... Chiudo terminali e Docker..."
    kitty @ close-window --match title:flagchecker
    kitty @ close-window --match title:cookieserver
    kitty @ close-window --match title:service
    kitty @ close-window --match title:frontend
    docker compose down
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

# Run Server
echo "🍪 Avvio CookieFarm Server..."

cd ../server/backend/

kitty --title "cookieserver" bash -c "make build && make run ARGS=''; exec bash" &

echo "✅ Server avviato!"
sleep 3
echo "📡 Invio configurazione..."

cd ../../scripts/
chmod +x shitcurl.py
./shitcurl.py

echo "✅ Configurazione inviata!"

# Run FE
echo "🌐 Start frontend"
cd ../server/frontend/
kitty --title "frontend" bash -c "/bin/bun run dev; exec bash" &
echo "🌐 Frontend started"

# Run Services
echo "🚀 Avvio Servizi..."

cd ../../tests
chmod +x ./start_containers.sh
kitty --title "service" bash -c "./start_containers.sh $1; exec bash" &

echo "🚀 Servizi avviati!"

echo "🎯 Cookie Farm Server pronto all'uso!"

# Attendi input per terminare tutti i terminali kitty
read -p "🔻 Premi INVIO per chiudere tutti i terminali avviati dallo script..."

cleanup()

echo "🧹 Tutti i terminali sono stati chiusi!"
