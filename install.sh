#!/bin/bash

set -euo pipefail
set -o errtrace

# ======================================= COLORS =============================================================

if [ -t 1 ] && [ -n "${TERM:-}" ]; then
  ESC=$'\033['
else
  ESC=""
fi

RESET="${ESC}0m"
BOLD="${ESC}1m"
DIM="${ESC}2m"
ITALIC="${ESC}3m"
UNDERLINE="${ESC}4m"
BLINK="${ESC}5m"
INVERSE="${ESC}7m"
HIDDEN="${ESC}8m"
STRIKE="${ESC}9m"

# -- CookieFarm color scheme (truecolor) --------------------------------------
# Mapped from the Go fang.ColorScheme used in the CLI.
# Requires a truecolor-capable terminal (most modern terminals).
C_TITLE=$'\033[38;2;205;161;87m'       # #CDA157 golden amber  — headers, section titles
C_BASE=$'\033[38;2;233;233;233m'       # #E9E9E9 light gray    — default text
C_DESCRIPTION=$'\033[38;2;217;217;217m' # #D9D9D9 light gray   — descriptions, body
C_DIM=$'\033[38;2;136;136;136m'        # #888888 mid gray      — dimmed, comments, help, dash
C_FLAG=$'\033[38;2;33;150;243m'        # #2196F3 blue          — flags/keys
C_GREEN=$'\033[38;2;33;155;84m'        # #219B54 green         — quoted strings, success values
C_ARGUMENT=$'\033[38;2;237;237;237m'   # #EDEDED near white    — arguments, primary text
C_ERROR_FG=$'\033[38;2;237;237;237m'   # #EDEDED near white    — error header foreground
C_ERROR_BG=$'\033[48;2;231;76;60m'     # #E74C3C red bg        — error header background
C_ERROR=$'\033[38;2;231;76;60m'        # #E74C3C red           — error details

# Gum hex strings (used with --foreground / --border-foreground)
GC_TITLE="#CDA157"
GC_DIM="#888888"
GC_FLAG="#2196F3"
GC_GREEN="#219B54"
GC_ERROR="#E74C3C"
GC_BASE="#E9E9E9"
GC_ARGUMENT="#EDEDED"

# Semantic aliases
ERROR="${BOLD}${C_ERROR_BG}${C_ERROR_FG}"
SUCCESS="${BOLD}${C_GREEN}"
INFO="${BOLD}${C_FLAG}"
WARN="${BOLD}${C_TITLE}"

CLEAR_LINE="${ESC}2K"
MOVE_START_LINE="${ESC}0G"

color_text() {
  if [ $# -lt 2 ]; then return 1; fi
  local text last_arg_index=$#
  text="${!last_arg_index}"
  local styles=()
  local prev_idx=$((last_arg_index - 1))
  local prev_arg="${!prev_idx}"
  if [ "$prev_arg" = "--" ]; then
    local style_count=$((last_arg_index - 2))
    [ "$style_count" -gt 0 ] && styles=("${@:1:$style_count}") || styles=()
  else
    styles=("${@:1:$((last_arg_index - 1))}")
  fi
  local seq="" s
  for s in "${styles[@]}"; do seq="${seq}${!s:-}"; done
  if [ -t 1 ]; then printf "%b%s%b" "$seq" "$text" "$RESET"
  else printf "%s" "$text"; fi
}

styled_echo() { color_text "${@:1:$#-1}" "${!#}"; printf "\n"; }

# Section header helper using the golden amber Title color
section_header() {
  printf "%b%b-- %s --%b\n" "$BOLD" "$C_TITLE" "$1" "$RESET"
}

# ========================================== ERROR HANDLING ==================================================

last_command=""
current_command=""
trap 'last_command=$current_command; current_command=$BASH_COMMAND' DEBUG

err_report() {
  local exit_code=$?
  local lineno=${1:-${BASH_LINENO[0]:-0}}
  [ "$exit_code" -eq 0 ] && return 0
  printf "%b ERROR %b %s%b (exit %d, line %d)\n" \
    "${BOLD}${C_ERROR_BG}${C_ERROR_FG}" "$RESET" \
    "${last_command:-unknown}" "$RESET" \
    "$exit_code" "$lineno" >&2
}

on_exit() {
  local code=$?
  [ "$code" -ne 0 ] && err_report
}

trap 'err_report ${LINENO}' ERR
trap on_exit EXIT

require_cmd() {
  local cmd="$1" help_msg="${2:-}"
  if ! command -v "$cmd" >/dev/null 2>&1; then
    if [ -n "$help_msg" ]; then
      printf "%b ERROR %b Missing command: %b%s%b — %s\n" \
        "${BOLD}${C_ERROR_BG}${C_ERROR_FG}" "$RESET" \
        "${BOLD}${C_FLAG}" "$cmd" "$RESET" "$help_msg" >&2
    else
      printf "%b ERROR %b Missing command: %b%s%b\n" \
        "${BOLD}${C_ERROR_BG}${C_ERROR_FG}" "$RESET" \
        "${BOLD}${C_FLAG}" "$cmd" "$RESET" >&2
    fi
    exit 2
  fi
}

# ========================================== GUM SETUP ==========================================================

gum_bin() {
  local repo="charmbracelet/gum"
  local api="https://api.github.com/repos/${repo}/releases/latest"
  local os arch asset version url tmpdir tarball bindir

  os="$(uname -s)"
  arch="$(uname -m)"

  case "$os" in
    Linux)  os="Linux" ;;
    Darwin) os="Darwin" ;;
    *) printf "%b ERROR %b Unsupported OS: %s\n" "${BOLD}${C_ERROR_BG}${C_ERROR_FG}" "$RESET" "$os" >&2; return 1 ;;
  esac

  case "$arch" in
    x86_64|amd64)  arch="x86_64" ;;
    aarch64|arm64) arch="arm64" ;;
    armv7l|armv7)  arch="armv7" ;;
    *) printf "%b ERROR %b Unsupported arch: %s\n" "${BOLD}${C_ERROR_BG}${C_ERROR_FG}" "$RESET" "$arch" >&2; return 1 ;;
  esac

  require_cmd curl  "Install curl to download gum"
  require_cmd tar   "Install tar to extract the archive"
  require_cmd mktemp

  version="$(curl -fsSL "$api" | sed -n 's/.*"tag_name": *"v\([^"]*\)".*/\1/p' | head -n1)"
  if [ -z "${version:-}" ]; then
    printf "%b ERROR %b Unable to fetch latest gum version from GitHub\n" "${BOLD}${C_ERROR_BG}${C_ERROR_FG}" "$RESET" >&2
    return 1
  fi

  asset="gum_${version}_${os}_${arch}.tar.gz"
  url="https://github.com/${repo}/releases/download/v${version}/${asset}"

  tmpdir="$(mktemp -d /tmp/gum.XXXXXX)"
  tarball="${tmpdir}/${asset}"
  bindir="${tmpdir}/bin"

  mkdir -p "$bindir"
  printf "%b  →%b Downloading gum %b%s%b...\n" "$C_DIM" "$RESET" "${BOLD}${C_TITLE}" "$version" "$RESET"
  curl -fL --progress-bar "$url" -o "$tarball"
  tar -xzf "$tarball" -C "$bindir"
  chmod +x "${bindir}/gum"
  printf '%s\n' "${bindir}/gum"
}

ensure_gum() {
  if [ -x /tmp/gum-bin/gum ]; then
    printf '%s\n' /tmp/gum-bin/gum
    return 0
  fi

  local found
  found="$(find /tmp -maxdepth 4 -type f -name gum -executable -path '/tmp/gum.*' -print -quit 2>/dev/null || true)"
  if [ -n "${found:-}" ]; then
    mkdir -p /tmp/gum-bin
    cp "$found" /tmp/gum-bin/gum
    chmod +x /tmp/gum-bin/gum
    printf '%s\n' /tmp/gum-bin/gum
    return 0
  fi

  local bin
  bin="$(gum_bin)"
  if [ -z "${bin:-}" ] || [ ! -x "$bin" ]; then
    printf "%b ERROR %b gum is not available and could not be installed\n" "${BOLD}${C_ERROR_BG}${C_ERROR_FG}" "$RESET" >&2
    return 1
  fi
  mkdir -p /tmp/gum-bin
  cp "$bin" /tmp/gum-bin/gum
  chmod +x /tmp/gum-bin/gum
  printf '%s\n' /tmp/gum-bin/gum
}

# ========================================== MAIN SCRIPT =========================================================

BANNER="CookieFarm Installer"

printf "%b%b%s%b\n" "$BOLD" "$C_TITLE" "$BANNER" "$RESET"

require_cmd bash
require_cmd printf
require_cmd sed
require_cmd grep
require_cmd find
require_cmd cp

GUM_BIN="$(ensure_gum)" || {
  printf "%b ERROR %b Failed to obtain gum binary\n" "${BOLD}${C_ERROR_BG}${C_ERROR_FG}" "$RESET" >&2
  exit 1
}

build="${1:-}"
if [ -z "$build" ]; then
  printf "%b%bUsage:%b %s <true|false>\n" "$BOLD" "$C_FLAG" "$RESET" "$0"
  exit 1
fi

# --------------------------------- Gum themed helpers ---------------------------------

gum_input() {
  local header="$1"; shift
  "$GUM_BIN" input \
    --header "$header" \
    --header.foreground "$GC_TITLE" \
    --prompt "> " \
    --prompt.foreground "$GC_DIM" \
    --cursor.foreground "$GC_TITLE" \
    "$@"
}

gum_confirm() {
  "$GUM_BIN" confirm \
    --selected.background "$GC_TITLE" \
    --selected.foreground "#0A0C0D" \
    --unselected.foreground "$GC_DIM" \
    "$@"
}

gum_spin() {
  # gum_spin "title" cmd args...
  local title="$1"; shift
  "$GUM_BIN" spin \
    --spinner dot \
    --title "$title" \
    --spinner.foreground "$GC_TITLE" \
    --title.foreground "$GC_BASE" \
    -- "$@"
}

gum_style_section() {
  "$GUM_BIN" style \
    --bold \
    --foreground "$GC_BASE" \
    "$1"
}

# --------------------------------- State globals ---------------------------------

PORT=""
PASSWORD=""
CONFIG_FILE=""
DEBUG=""
preconfigured="false"
url_flag_checker=""
team_token=""
submit_flag_checker_time=""
max_flag_batch_size=""
protocol=""
tick_time=""
flag_ttl=""
start_time=""
end_time=""
services_yaml=""
range_ip_teams=""
format_ip_teams=""
my_team_id=""
regex_flag=""
nop_team=""
url_flag_ids=""

# --------------------------------- Section 1: Basic .env ---------------------------------

gum_ask_basic() {
  gum_style_section "[-] BASIC CONFIGURATION [-]"

  PORT="$(gum_input "Server PORT" --value "8080" --placeholder "8080")"
  PASSWORD="$(gum_input "Server PASSWORD" --value "password" --placeholder "password")"

  if gum_confirm --default "Use a config file (CONFIG_FILE)?"; then
    CONFIG_FILE="true"
  else
    CONFIG_FILE="false"
  fi

  if gum_confirm "Enable DEBUG mode?"; then
    DEBUG="true"
  else
    DEBUG="false"
  fi
}

# --------------------------------- Section 2: config.yml ---------------------------------

gum_ask_config() {
  local dest="${1:-cookiefarm/config.yml}"

  if [ -d "$dest" ] || [ "${dest: -1}" = "/" ] || [ "$dest" = "." ]; then
    dest="${dest%/}/config.yml"
  fi

  if [ "${CONFIG_FILE}" = "true" ]; then
    gum_style_section "[-] Configuration file [-]"

    if gum_confirm "Do you already have a preconfigured config.yml?"; then
      cfg_path="$(gum_input "Path to your config.yml" --value "./config.yml" --placeholder "./config.yml")"
      if [ -f "$cfg_path" ]; then
        mkdir -p "$(dirname "$dest")"
        cp "$cfg_path" "$dest"
        printf "%b✔%b Copied config to %b%s%b\n" "$C_GREEN" "$RESET" "${BOLD}${C_ARGUMENT}" "$dest" "$RESET"
        preconfigured="true"
      else
        printf "%b⚠%b Path not found — switching to interactive setup\n" "$C_TITLE" "$RESET"
        preconfigured="false"
      fi
    fi

    if [ "$preconfigured" = "false" ]; then
      gum_style_section "[-] Server settings [-]"

      url_flag_checker="$(gum_input "url_flag_checker" \
        --value "http://localhost:5001/flags" --placeholder "http://<ip>:8081/flags")"
      team_token="$(gum_input "team_token" --placeholder "your team token")"
      submit_flag_checker_time="$(gum_input "submit_flag_checker_time" --value "120" --placeholder "120")"
      max_flag_batch_size="$(gum_input "max_flag_batch_size" --value "1000" --placeholder "1000")"
      protocol="$(gum_input "protocol" --value "cc_http" --placeholder "cc_http")"
      tick_time="$(gum_input "tick_time" --value "120" --placeholder "120")"
      flag_ttl="$(gum_input "flag_ttl (in ticks)" --value "0" --placeholder "0")"

      start_time_default="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"
      start_time="$(gum_input "start_time (ISO 8601 UTC)" \
        --value "$start_time_default" --placeholder "$start_time_default")"

      if command -v python3 >/dev/null 2>&1; then
        end_time_default="$(python3 -c "
from datetime import datetime, timezone, timedelta
print((datetime.now(timezone.utc)+timedelta(hours=8)).strftime('%Y-%m-%dT%H:%M:%SZ'))
")"
      else
        case "$(uname -s)" in
          Linux)  end_time_default="$(date -u -d "+8 hours" +"%Y-%m-%dT%H:%M:%SZ")" ;;
          Darwin) end_time_default="$(date -u -v+8H +"%Y-%m-%dT%H:%M:%SZ")" ;;
          *)      end_time_default="$start_time_default" ;;
        esac
      fi
      end_time="$(gum_input "end_time (ISO 8601 UTC, default now+8h)" \
        --value "$end_time_default" --placeholder "$end_time_default")"

      gum_style_section "[-] Services [-]"

      num_services_str="$(gum_input "How many services to configure?" --value "3" --placeholder "3")"
      if printf '%s' "$num_services_str" | grep -Eq '^[0-9]+$'; then
        num_services="$num_services_str"
      else
        num_services=0
      fi

      services_yaml=""
      if [ "$num_services" -gt 0 ]; then
        i=1
        while [ "$i" -le "$num_services" ]; do
          svc_name="$(gum_input "Service #${i} name" --value "Service${i}" --placeholder "ServiceName")"
          svc_port="$(gum_input "Port for ${svc_name}" --value "808${i}" --placeholder "8080")"
          services_yaml="${services_yaml}    ${svc_name}: ${svc_port}"$'\n'
          i=$((i+1))
        done
      fi

      gum_style_section "[-] Shared / Teams [-]"

      range_ip_teams="$(gum_input "range_ip_teams" --value "40" --placeholder "40")"
      format_ip_teams="$(gum_input "format_ip_teams" \
        --value "10.10.{}.1" --placeholder "10.10.{}.1")"
      my_team_id="$(gum_input "my_team_id" --value "1" --placeholder "1")"
      regex_flag="$(gum_input "regex_flag" --value "[A-Z0-9]{31}=" --placeholder "[A-Z0-9]{31}=")"
      nop_team="$(gum_input "nop_team" --value "0" --placeholder "0")"
      url_flag_ids="$(gum_input "url_flag_ids" \
        --value "http://localhost:5001/flagIds" --placeholder "http://<ip>:8081/flagIds")"

      mkdir -p "$(dirname "$dest")"
      cat > "$dest" <<EOF
configured: true

# Server
server:
  url_flag_checker: "${url_flag_checker}"
  team_token: "${team_token}"
  submit_flag_checker_time: ${submit_flag_checker_time}
  max_flag_batch_size: ${max_flag_batch_size}
  protocol: "${protocol}"
  tick_time: ${tick_time}
  flag_ttl: ${flag_ttl}
  start_time: "${start_time}"
  end_time: "${end_time}"

# Client
shared:
  services:
${services_yaml}  range_ip_teams: ${range_ip_teams}
  format_ip_teams: "${format_ip_teams}"
  my_team_id: ${my_team_id}
  regex_flag: "${regex_flag}"
  nop_team: ${nop_team}
  url_flag_ids: "${url_flag_ids}"
EOF
      printf "%b✔%b Wrote %b%s%b\n" "$C_GREEN" "$RESET" "${BOLD}${C_ARGUMENT}" "$dest" "$RESET"
    fi
  fi
}

# --------------------------------- Write .env ---------------------------------

write_env() {
  local dest="${1:-cookiefarm/.env}"

  if [ -d "$dest" ] || [ "${dest: -1}" = "/" ]; then
    dest="${dest%/}/.env"
  fi

  mkdir -p "$(dirname "$dest")"
  cat > "$dest" <<EOF
PORT=${PORT}
PASSWORD=${PASSWORD}
CONFIG_FILE=${CONFIG_FILE}
DEBUG=${DEBUG}
EOF
  printf "%b✔%b .env written → %b%s%b  PORT=%b%s%b  CONFIG_FILE=%b%s%b  DEBUG=%b%s%b\n" \
    "$C_GREEN" "$RESET" \
    "${BOLD}${C_ARGUMENT}" "$dest" "$RESET" \
    "$C_FLAG" "$PORT" "$RESET" \
    "$C_FLAG" "$CONFIG_FILE" "$RESET" \
    "$C_FLAG" "$DEBUG" "$RESET"
}

# --------------------------------- Deploy ---------------------------------

if [ "$build" = "true" ]; then
  if [ -d CookieFarm ]; then
    printf "%b⚠%b CookieFarm directory exists — skipping clone\n" "$C_TITLE" "$RESET"
  else
    gum_spin "Cloning CookieFarm..." \
      git clone https://github.com/ByteTheCookies/CookieFarm.git
  fi

  cd CookieFarm

  gum_ask_basic
  write_env
  gum_ask_config

  cd cookiefarm
  gum_spin "Building and starting containers..." \
    sh -c 'docker compose -f docker-compose.yml up --build -d 2>&1'
else
  mkdir -p CookieFarm
  cd CookieFarm

  gum_spin "Downloading docker-compose.yml..." \
    wget -q https://raw.githubusercontent.com/ByteTheCookies/CookieFarm/refs/heads/main/cookiefarm/docker-compose.yml

  gum_ask_basic
  write_env .
  gum_ask_config .

  gum_spin "Starting containers..." \
    sh -c 'docker compose -f docker-compose.yml up --build -d 2>&1'
fi

"$GUM_BIN" style \
  --border rounded \
  --border-foreground "$GC_TITLE" \
  --foreground "$GC_ARGUMENT" \
  --bold \
  --padding "1 3" \
  "🍪  CookieFarm is up!"
