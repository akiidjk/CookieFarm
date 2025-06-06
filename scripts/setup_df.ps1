# Windows setup script for DestructiveFarm integration with CookieFarm
# PowerShell version of setup_df.sh

param(
    [Parameter(Mandatory=$true)]
    [int]$numContainers,
    
    [Parameter(Mandatory=$true)]
    [string]$pathDf
)

# Check parameters
if ($numContainers -lt 1 -or $numContainers -gt 10) {
    Write-Host "Usage:`n  .\setup_df.ps1 <num_containers> <path_df>"
    Write-Host "  num_containers: Number of containers to start (1-10)"
    Write-Host "  path_df: Path to DestructiveFarm directory"
    exit 1
}

# === CLEANUP FUNCTION ===
function Cleanup {
    Write-Host "🧹 Cleaning up... Closing terminals and Docker..."
    Get-Process | Where-Object { $_.MainWindowTitle -match "flagchecker|cookieserver|service|frontend|destructivefarm" } | Stop-Process -Force
    Set-Location "..\tests\"
    docker compose down
    Set-Location "..\scripts\"
    exit
}

# Register cleanup on Ctrl+C
$null = [Console]::TreatControlCAsInput = $true
[Console]::CancelKeyPress += {
    Cleanup
}

# Install requirements
Write-Host "📦 Installing Python dependencies..."
python -m pip install --upgrade pip | Out-Null
python -m pip install -r "..\requirements.txt" | Out-Null

$VENV_ACTIVATE = "..\venv\Scripts\Activate.ps1"

# Create function to start processes in new windows
function Start-ProcessInNewWindow {
    param(
        [string]$Title,
        [string]$Command
    )
    
    $encodedCommand = [Convert]::ToBase64String([System.Text.Encoding]::Unicode.GetBytes($Command))
    Start-Process powershell -ArgumentList "-NoExit", "-NoProfile", "-EncodedCommand $encodedCommand" -WindowStyle Normal
}

# Run Flagchecker
Write-Host "🚩 Starting Flagchecker..."
$flagcheckerCmd = "Set-Location '$(Resolve-Path -Path '..\tests\')'; & '$VENV_ACTIVATE'; python flagchecker.py"
Start-ProcessInNewWindow -Title "flagchecker" -Command $flagcheckerCmd
Write-Host "✅ Flagchecker launched in a separate terminal! 🎉"
Write-Host ""

# Run Services
Write-Host "🚀 Starting Services..."
Set-Location "..\tests\"
$serviceCmd = "Set-Location '$(Resolve-Path -Path '..\tests\')'; .\start_containers.ps1 $numContainers"
Start-ProcessInNewWindow -Title "service" -Command $serviceCmd
Write-Host "🚀 Services started!"

# Run DestructiveFarm
Write-Host "🚀 Starting DestructiveFarm..."
Set-Location "..\scripts\"

# Copy configuration to DestructiveFarm
Copy-Item -Path ".\config_df.py" -Destination "$pathDf\server\config.py" -Force
$dfCmd = "Set-Location '$pathDf\server\'; python start_server.py"
Start-ProcessInNewWindow -Title "destructivefarm" -Command $dfCmd
Write-Host "🚀 DestructiveFarm started!"

Write-Host "🎯 Environment for DF ready to use!"

# Wait for input to terminate all terminals
Write-Host "🔻 Press ENTER to close all terminals launched by the script..."
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")

Cleanup