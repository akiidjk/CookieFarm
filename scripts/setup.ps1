# Windows setup script for CookieFarm
# PowerShell version of setup.sh

param(
    [Parameter(Mandatory=$true)]
    [int]$numContainers,
    
    [Parameter(Mandatory=$true)]
    [int]$productionMode
)

# === CONFIG ===
$TOOLS_DIR = "..\server\tools"
$VENV_ACTIVATE = "..\venv\Scripts\Activate.ps1"
$FLAGCHECKER_SCRIPT = "..\tests\flagchecker.py"
$SERVER_DIR = "..\server"
$SCRIPTS_DIR = "..\scripts"
$TESTS_DIR = "..\tests"
$REQUIREMENTS = "..\requirements.txt"
$TAILWIND_URL = "https://github.com/tailwindlabs/tailwindcss/releases/download/v4.1.4/tailwindcss-windows-x64.exe"

# === USAGE CHECK ===
if ($numContainers -lt 1 -or $numContainers -gt 10) {
    Write-Host "Usage:`n  .\setup.ps1 <num_containers> <production_mode>`n"
    Write-Host "  num_containers: Number of containers to start (1-10)"
    Write-Host "  production_mode: 0 for development, 1 for production"
    exit 1
}

# === CLEANUP FUNCTION ===
function Cleanup {
    Write-Host "🧹 Cleaning up... Closing terminals and Docker..."
    Get-Process | Where-Object { $_.MainWindowTitle -match "flagchecker|cookieserver|service|frontend" } | Stop-Process -Force
    Set-Location $SERVER_DIR
    docker compose down
    exit
}

# Register cleanup on Ctrl+C
$null = [Console]::TreatControlCAsInput = $true
[Console]::CancelKeyPress += {
    Cleanup
}

# === REQUIREMENTS ===
Write-Host "📦 Installing Python dependencies..."
python -m pip install --upgrade pip | Out-Null
#python -m pip install -r "$REQUIREMENTS" | Out-Null

# === TAILWIND ===
Write-Host "🎨 Checking TailwindCSS..."
if (-not (Test-Path -Path $TOOLS_DIR)) {
    New-Item -ItemType Directory -Path $TOOLS_DIR -Force | Out-Null
}

if (-not (Test-Path -Path "$TOOLS_DIR\tailwindcss.exe")) {
    Invoke-WebRequest -Uri $TAILWIND_URL -OutFile "$TOOLS_DIR\tailwindcss.exe"
    Write-Host "✅ tailwindcss installed."
}

# === MINIFY ===
Write-Host "📦 Checking minify..."
npm install uglify-js -g

# Create function to start processes in new windows
function Start-ProcessInNewWindow {
    param(
        [string]$Title,
        [string]$Command
    )
    
    $encodedCommand = [Convert]::ToBase64String([System.Text.Encoding]::Unicode.GetBytes($Command))
    Start-Process powershell -ArgumentList "-NoExit", "-NoProfile", "-EncodedCommand $encodedCommand" -WindowStyle Normal
}

# === FLAGCHECKER ===
Write-Host "🚩 Starting Flagchecker..."
$flagcheckerCmd = "Set-Location $(Resolve-Path -Path $TESTS_DIR); & $VENV_ACTIVATE; python flagchecker.py"
Start-ProcessInNewWindow -Title "flagchecker" -Command $flagcheckerCmd
Write-Host "✅ Flagchecker launched in a separate terminal! 🎉"

# === SERVER ===
Write-Host "🍪 Starting CookieFarm Server..."
Set-Location $SERVER_DIR

if ($productionMode -eq 1) {
    Write-Host "🔒 Production mode enabled!"
    $serverCmd = "Set-Location '$(Resolve-Path -Path $SERVER_DIR)'; make build-plugins-prod; make run-prod; .\cookieserver.exe"
    Start-ProcessInNewWindow -Title "cookieserver" -Command $serverCmd
} else {
    Write-Host "🔓 Development mode enabled!"
    $serverCmd = "Set-Location '$(Resolve-Path -Path $SERVER_DIR)'; make build-plugins; make run 'ARGS=--config config.yml --debug'"
    Start-ProcessInNewWindow -Title "cookieserver" -Command $serverCmd
}
Write-Host "✅ Server started!"

Start-Sleep -Seconds 3

# # === SENDING CONFIG ===
# Write-Host "📡 Sending configuration..."
# Set-Location $SCRIPTS_DIR
# python shitcurl.py
# Write-Host "✅ Configuration sent!"

# === FRONTEND ===
Write-Host "🌐 Starting Frontend..."
Set-Location $SERVER_DIR
$frontendCmd = "Set-Location '$(Resolve-Path -Path $SERVER_DIR)'; make tailwindcss-build"
Start-ProcessInNewWindow -Title "frontend" -Command $frontendCmd
Write-Host "✅ Frontend started!"

# === SERVICES ===
Write-Host "🚀 Starting Services..."
Set-Location $TESTS_DIR
$serviceCmd = "Set-Location '$(Resolve-Path -Path $TESTS_DIR)'; .\start_containers.ps1 $numContainers"
Start-ProcessInNewWindow -Title "service" -Command $serviceCmd
Write-Host "✅ Services started!"

# === COMPLETION ===
Write-Host "`n🎯 Cookie Farm Server ready to use!"

Write-Host "🔻 Press ENTER to close all terminals launched by the script..."
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
Cleanup