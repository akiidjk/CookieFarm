# Windows Setup Scripts for CookieFarm

This document explains how to use the Windows-adapted setup scripts for CookieFarm.

## Prerequisites

Before using these scripts, ensure you have the following installed:
- [Python](https://www.python.org/downloads/) (3.7+)
- [PowerShell](https://docs.microsoft.com/en-us/powershell/) 5.1+
- [Docker Desktop for Windows](https://www.docker.com/products/docker-desktop)
- Git Bash or similar (optional, for running original bash scripts if needed)

## Setup Scripts

Two main setup scripts are provided for Windows environments:

### 1. Standard Setup (`setup.ps1`)

This script sets up the CookieFarm environment with all necessary components.

#### Usage:
```powershell
.\setup.ps1 <num_containers> <production_mode>
```

#### Parameters:
- `num_containers`: Number of service containers to start (1-10)
- `production_mode`: Use 0 for development mode, 1 for production mode

#### Example:
```powershell
.\setup.ps1 3 0  # Start 3 containers in development mode
```

### 2. DestructiveFarm Setup (`setup_df.ps1`)

This script sets up the environment for running CookieFarm with DestructiveFarm integration.

#### Usage:
```powershell
.\setup_df.ps1 <num_containers> <path_df>
```

#### Parameters:
- `num_containers`: Number of service containers to start (1-10)
- `path_df`: Full path to the DestructiveFarm directory

#### Example:
```powershell
.\setup_df.ps1 3 C:\path\to\DestructiveFarm
```

## Execution Policy

If you receive security errors when running the scripts, you may need to adjust the PowerShell execution policy. Run PowerShell as Administrator and execute:

```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

## Running the Scripts

1. Open PowerShell
2. Navigate to the scripts directory:
   ```powershell
   cd path\to\CookieFarm\scripts
   ```
3. Execute the desired script with appropriate parameters

## Troubleshooting

### Docker Issues
- Ensure Docker Desktop is running before executing the scripts
- If you encounter Docker network errors, try restarting Docker Desktop

### PowerShell Window Management
- If new PowerShell windows don't open correctly, you may need to run the scripts with elevated privileges (as Administrator)
- The scripts attempt to manage PowerShell windows automatically. If this fails, you may need to manually close them after use

### Python Virtual Environment
- If you encounter Python errors, ensure your virtual environment is properly set up
- You may need to manually activate the virtual environment:
  ```powershell
  ..\venv\Scripts\Activate.ps1
  ```

## Differences from Linux Scripts

The Windows PowerShell scripts differ from the original Bash scripts in several ways:
- Use PowerShell's process management instead of kitty terminal
- Windows-compatible path handling
- PowerShell-specific syntax for command execution
- Adapted cleanup procedures for Windows processes
- Windows-compatible download URLs for tools
- Modified environment variable handling

For any issues not covered here, please refer to the project documentation or open an issue on the project repository.
