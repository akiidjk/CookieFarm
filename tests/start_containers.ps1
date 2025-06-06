param(
    [Parameter(Mandatory=$true)]
    [int]$num_containers
)

# Check input validity
if ($num_containers -lt 1 -or $num_containers -gt 10) {
    Write-Host "⚠️ Error: Number of containers must be between 1 and 10"
    exit 1
}

# Function to handle cleanup on script exit
function Cleanup {
    Write-Host "`n🧹 Stopping all containers..."
    docker compose down
    Write-Host "✅ All containers stopped!"
}

# Register cleanup on Ctrl+C
[Console]::TreatControlCAsInput = $true
[Console]::CancelKeyPress += {
    Cleanup
    exit
}

# Set the environment variable for Docker Compose
$env:NUM_CONTAINERS = $num_containers

# Start the containers
Write-Host "🚀 Starting $num_containers service containers..."
docker compose up

# Script completed
Write-Host "✅ All containers stopped!"