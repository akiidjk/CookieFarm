#!/bin/bash

set -e

# === CONFIG ===
TOOLS_DIR="../server/tools"
VENV_ACTIVATE="../venv/bin/activate"
FLAGCHECKER_SCRIPT="../tests/flagchecker.py"
SERVER_DIR="../server"
SCRIPTS_DIR="../scripts"
TESTS_DIR="../tests"
REQUIREMENTS="../requirements.txt"
TAILWIND_URL="https://github.com/tailwindlabs/tailwindcss/releases/download/v4.1.4/tailwindcss-linux-x64"
MINIFY_URL="https://github.com/tdewolff/minify/releases/download/v2.23.1/minify_linux_amd64.tar.gz"

# === USAGE CHECK ===
if [[ $# -ne 1 ]]; then
    echo -e "Usage:\n  ./setup.sh <num_containers>"
    exit 1
fi

# === CLEANUP HANDLER ===
cleanup() {
    echo "🧹 Pulizia in corso... Chiudo terminali e Docker..."
    kitty @ close-window --match title:flagchecker || true
    kitty @ close-window --match title:cookieserver || true
    kitty @ close-window --match title:service || true
    kitty @ close-window --match title:frontend || true
    docker compose down
    exit
}
trap cleanup SIGINT

# === REQUIREMENTS ===
echo "📦 Installazione dipendenze Python..."
pip install --upgrade pip > /dev/null
pip install -r "$REQUIREMENTS" > /dev/null

# === TAILWIND ===
echo "🎨 Controllo TailwindCSS..."
mkdir -p "$TOOLS_DIR"
if [ ! -f "$TOOLS_DIR/tailwindcss" ]; then
    wget -q "$TAILWIND_URL" -O "$TOOLS_DIR/tailwindcss"
    chmod +x "$TOOLS_DIR/tailwindcss"
    echo "✅ tailwindcss installato."
fi

# === MINIFY ===
echo "📦 Controllo minify..."
sudo npm install uglify-js -g

# === FLAGCHECKER ===
echo "🚩 Avvio Flagchecker..."
chmod +x "$FLAGCHECKER_SCRIPT"
kitty --title "flagchecker" bash -c "source $VENV_ACTIVATE && $FLAGCHECKER_SCRIPT; exec bash" &
echo "✅ Flagchecker lanciato in un terminale separato! 🎉"

# === SERVER ===
echo "🍪 Avvio CookieFarm Server..."
cd "$SERVER_DIR"
kitty --title "cookieserver" bash -c "make run ARGS='--debug'; exec bash" &
echo "✅ Server avviato!"

sleep 3

# === INVIO CONFIG ===
echo "📡 Invio configurazione..."
cd "$SCRIPTS_DIR"
chmod +x shitcurl.py
./shitcurl.py
echo "✅ Configurazione inviata!"

# === FRONTEND ===
echo "🌐 Avvio Frontend..."
cd "$SERVER_DIR"
kitty --title "frontend" bash -c "make tailwindcss-build; exec bash" &
echo "✅ Frontend avviato!"

# === SERVIZI ===
echo "🚀 Avvio Servizi..."
cd "$TESTS_DIR"
chmod +x ./start_containers.sh
kitty --title "service" bash -c "./start_containers.sh $1; exec bash" &
echo "✅ Servizi avviati!"

# === COMPLETAMENTO ===
echo -e "\n🎯 Cookie Farm Server pronto all'uso!"

read -p "🔻 Premi INVIO per chiudere tutti i terminali avviati dallo script..."
cleanup
