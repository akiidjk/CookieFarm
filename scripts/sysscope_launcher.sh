#!/bin/bash

if [ -z "$1" ]; then
    echo "Usage: $0 <duration>"
    exit 1
fi

if command -v sysscope >/dev/null 2>&1; then
    echo "SysScope installed"
else
    echo "SysScope not installed"
    git clone https://github.com/akiidjk/SysScope.git
    cd SysScope
    cargo install --path .
fi

pidServer=$(pgrep -f ./bin/cookieserver )
pidClient=$(pgrep -f ./bin/./cookieclient )
pidExploit=$(pgrep -f "python3.*WorkSpace.*example_.\.py")


kitty -e bash -c "sysscope -p $pidServer -d $1; exec bash " &
kitty -e bash -c "sysscope -p $pidClient -d $1; exec bash" &
kitty -e bash -c "sysscope -p $pidExploit -d $1; exec bash" &
