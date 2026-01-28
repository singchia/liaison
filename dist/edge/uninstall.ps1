# Liaison Edge 卸载脚本 (Windows PowerShell)

param(
    [switch]$Help,
    [switch]$KeepConfig,
    [switch]$KeepLogs
)

$ErrorActionPreference = "Stop"

function Show-Help {
    Write-Host "Usage: uninstall.ps1 [OPTIONS]"
    Write-Host ""
    Write-Host "Options:"
    Write-Host "  -Help        Show this help message"
    Write-Host "  -KeepConfig  Keep configuration files"
    Write-Host "  -KeepLogs    Keep log files"
    Write-Host ""
    Write-Host "Example:"
    Write-Host "  .\uninstall.ps1"
    Write-Host "  .\uninstall.ps1 -KeepConfig -KeepLogs"
    exit 0
}

if ($Help) { Show-Help }

Write-Host "Starting Liaison Edge uninstallation..." -ForegroundColor Yellow

# ---------------- Paths ----------------
$InstallDir = "C:\Program Files\Liaison"
$BinDir     = Join-Path $InstallDir "bin"
$ConfDir    = Join-Path $InstallDir "conf"
$LogDir     = Join-Path $InstallDir "logs"
$BinaryName = "liaison-edge.exe"
$BinaryPath = Join-Path $BinDir $BinaryName

# ---------------- Stop Process ----------------
Write-Host "Stopping Edge process..." -ForegroundColor Yellow

$processes = Get-Process -Name "liaison-edge" -ErrorAction SilentlyContinue
if ($processes) {
    foreach ($proc in $processes) {
        try {
            Write-Host "Stopping process: $($proc.Id)" -ForegroundColor Cyan
            Stop-Process -Id $proc.Id -Force -ErrorAction Stop
            Start-Sleep -Seconds 1
        } catch {
            Write-Host "Warning: Failed to stop process $($proc.Id): $_" -ForegroundColor Yellow
        }
    }
    Write-Host "Edge process stopped" -ForegroundColor Green
} else {
    Write-Host "No running Edge process found" -ForegroundColor Cyan
}

# ---------------- Remove Binary ----------------
Write-Host "Removing binary..." -ForegroundColor Yellow

if (Test-Path $BinaryPath) {
    try {
        Remove-Item -Path $BinaryPath -Force -ErrorAction Stop
        Write-Host "Binary removed: $BinaryPath" -ForegroundColor Green
    } catch {
        Write-Host "Error: Failed to remove binary: $_" -ForegroundColor Red
        Write-Host "Please make sure the Edge process is stopped and you have sufficient permissions" -ForegroundColor Yellow
        exit 1
    }
} else {
    Write-Host "Binary not found: $BinaryPath" -ForegroundColor Cyan
}

# ---------------- Remove Config ----------------
if (-not $KeepConfig) {
    Write-Host "Removing configuration..." -ForegroundColor Yellow
    
    if (Test-Path $ConfDir) {
        try {
            Remove-Item -Path $ConfDir -Recurse -Force -ErrorAction Stop
            Write-Host "Configuration removed: $ConfDir" -ForegroundColor Green
        } catch {
            Write-Host "Warning: Failed to remove configuration: $_" -ForegroundColor Yellow
        }
    } else {
        Write-Host "Configuration directory not found: $ConfDir" -ForegroundColor Cyan
    }
} else {
    Write-Host "Keeping configuration files (as requested)" -ForegroundColor Cyan
}

# ---------------- Remove Logs ----------------
if (-not $KeepLogs) {
    Write-Host "Removing logs..." -ForegroundColor Yellow
    
    if (Test-Path $LogDir) {
        try {
            Remove-Item -Path $LogDir -Recurse -Force -ErrorAction Stop
            Write-Host "Logs removed: $LogDir" -ForegroundColor Green
        } catch {
            Write-Host "Warning: Failed to remove logs: $_" -ForegroundColor Yellow
        }
    } else {
        Write-Host "Log directory not found: $LogDir" -ForegroundColor Cyan
    }
} else {
    Write-Host "Keeping log files (as requested)" -ForegroundColor Cyan
}

# ---------------- Remove Bin Directory ----------------
if (Test-Path $BinDir) {
    $remainingFiles = Get-ChildItem -Path $BinDir -ErrorAction SilentlyContinue
    if ($remainingFiles.Count -eq 0) {
        try {
            Remove-Item -Path $BinDir -Force -ErrorAction Stop
            Write-Host "Bin directory removed: $BinDir" -ForegroundColor Green
        } catch {
            Write-Host "Warning: Failed to remove bin directory: $_" -ForegroundColor Yellow
        }
    } else {
        Write-Host "Bin directory not empty, keeping: $BinDir" -ForegroundColor Cyan
    }
}

# ---------------- Remove Install Directory ----------------
if (Test-Path $InstallDir) {
    $remainingItems = Get-ChildItem -Path $InstallDir -ErrorAction SilentlyContinue
    if ($remainingItems.Count -eq 0) {
        try {
            Remove-Item -Path $InstallDir -Force -ErrorAction Stop
            Write-Host "Install directory removed: $InstallDir" -ForegroundColor Green
        } catch {
            Write-Host "Warning: Failed to remove install directory: $_" -ForegroundColor Yellow
        }
    } else {
        Write-Host "Install directory not empty, keeping: $InstallDir" -ForegroundColor Cyan
    }
}

Write-Host ""
Write-Host "Liaison Edge uninstallation completed!" -ForegroundColor Green
