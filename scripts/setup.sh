#!/bin/bash

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

kitty --title "cookieserver" bash -c "make build && make run ARGS='--debug'; exec bash" &

echo "✅ Server avviato!"
sleep 2
echo "📡 Invio configurazione..."

cd ../../scripts/
chmod +x shitcurl.py
./shitcurl.py

echo "✅ Configurazione inviata!"

# Run Services
echo "🚀 Avvio Servizi..."

chmod +x ../tests/service.py

kitty --title "service" bash -c "${activate_venv} && ../tests/service.py; exec bash" &

echo "🚀 Servizi avviati!"

echo "🎯 Cookie Farm Server pronto all'uso!"

# Attendi input per terminare tutti i terminali kitty
read -p "🔻 Premi INVIO per chiudere tutti i terminali avviati dallo script..."

# Chiudi le finestre kitty con i titoli assegnati
kitty @ close-window --match title:flagchecker
kitty @ close-window --match title:cookieserver
kitty @ close-window --match title:service

echo "🧹 Tutti i terminali sono stati chiusi!"
