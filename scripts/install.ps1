# PSX Installation Script for Windows

param(
    [Parameter(Position=0)]
    [ValidateSet('github', 'local', 'uninstall', 'help')]
    [string]$Command = 'github',
    
    [Parameter(Position=1)]
    [string]$Version = 'latest',
    
    [Parameter(Position=2)]
    [string]$Path = ''
)

$ErrorActionPreference = "Stop"

# Configuration
$BinaryName = "psx.exe"
$Repo = "m-mdy-m/psx"
$SystemInstallDir = "$env:ProgramFiles\PSX"
$UserInstallDir = "$env:LOCALAPPDATA\PSX"

Write-Host "PSX Installation Script" -ForegroundColor Blue
Write-Host "=======================" -ForegroundColor Blue
Write-Host ""

function Test-Administrator {
    $currentUser = [Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = New-Object Security.Principal.WindowsPrincipal($currentUser)
    return $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

function Add-ToPath {
    param([string]$Directory)
    
    $pathType = if (Test-Administrator) { "Machine" } else { "User" }
    $currentPath = [Environment]::GetEnvironmentVariable("Path", $pathType)
    
    if ($currentPath -notlike "*$Directory*") {
        Write-Host "Adding to PATH..." -ForegroundColor Yellow
        $newPath = if ($currentPath.EndsWith(';')) {
            "$currentPath$Directory"
        } else {
            "$currentPath;$Directory"
        }
        [Environment]::SetEnvironmentVariable("Path", $newPath, $pathType)
        
        # Update current session
        $env:Path = "$env:Path;$Directory"
        
        Write-Host "✓ Added to PATH" -ForegroundColor Green
        Write-Host ""
        Write-Host "Note: New terminals will have PATH updated automatically" -ForegroundColor Yellow
        Write-Host "For current terminal, restart or run:" -ForegroundColor Yellow
        Write-Host "  `$env:Path = [Environment]::GetEnvironmentVariable('Path', '$pathType')" -ForegroundColor Gray
    } else {
        Write-Host "Already in PATH" -ForegroundColor Gray
    }
}

function Test-PSXAvailable {
    try {
        $null = Get-Command psx -ErrorAction Stop
        return $true
    } catch {
        return $false
    }
}

function Install-FromGitHub {
    param([string]$Version)
    
    Write-Host "Platform: Windows (amd64)"
    Write-Host "Version: $Version"
    Write-Host ""
    
    # Get latest version if not specified
    if ($Version -eq 'latest') {
        Write-Host "Fetching latest version..." -ForegroundColor Yellow
        try {
            $response = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
            $Version = $response.tag_name
            Write-Host "Latest version: $Version"
        } catch {
            Write-Host "Could not determine latest version" -ForegroundColor Red
            exit 1
        }
    }
    
    # Build download URL
    $filename = "psx-windows-amd64.exe"
    $url = "https://github.com/$Repo/releases/download/$Version/$filename"
    $tmpFile = "$env:TEMP\psx.exe"
    
    Write-Host ""
    Write-Host "Downloading $filename..." -ForegroundColor Yellow
    
    try {
        Invoke-WebRequest -Uri $url -OutFile $tmpFile
    } catch {
        Write-Host "Download failed: $_" -ForegroundColor Red
        exit 1
    }
    
    # Determine install directory
    $installDir = if (Test-Administrator) { $SystemInstallDir } else { $UserInstallDir }
    
    Write-Host ""
    Write-Host "Installing to $installDir..." -ForegroundColor Yellow
    
    # Create directory if doesn't exist
    if (-not (Test-Path $installDir)) {
        New-Item -ItemType Directory -Path $installDir -Force | Out-Null
    }
    
    # Copy binary
    try {
        Copy-Item $tmpFile "$installDir\$BinaryName" -Force
        Remove-Item $tmpFile -Force
        
        Write-Host "✓ PSX installed successfully!" -ForegroundColor Green
        Write-Host ""
        Write-Host "Location: $installDir\$BinaryName"
        Write-Host "Version: $Version"
        Write-Host ""
        
        # Add to PATH
        Add-ToPath $installDir
        
        # Test if available
        if (Test-PSXAvailable) {
            Write-Host ""
            Write-Host "✓ psx command is available" -ForegroundColor Green
            Write-Host "Try: psx --version"
        } else {
            Write-Host ""
            Write-Host "Note: You may need to restart your terminal" -ForegroundColor Yellow
        }
        
    } catch {
        Write-Host "Installation failed: $_" -ForegroundColor Red
        exit 1
    }
}

function Install-FromLocal {
    param([string]$BinaryPath)
    
    if (-not $BinaryPath) {
        $BinaryPath = "build\psx.exe"
    }
    
    if (-not (Test-Path $BinaryPath)) {
        Write-Host "Binary not found: $BinaryPath" -ForegroundColor Red
        Write-Host "Build first: .\scripts\build.ps1"
        exit 1
    }
    
    Write-Host "Installing from: $BinaryPath"
    
    # Determine install directory
    $installDir = if (Test-Administrator) { $SystemInstallDir } else { $UserInstallDir }
    
    Write-Host "Installing to $installDir..." -ForegroundColor Yellow
    
    # Create directory if doesn't exist
    if (-not (Test-Path $installDir)) {
        New-Item -ItemType Directory -Path $installDir -Force | Out-Null
    }
    
    # Copy binary
    try {
        Copy-Item $BinaryPath "$installDir\$BinaryName" -Force
        
        Write-Host "✓ PSX installed successfully!" -ForegroundColor Green
        Write-Host ""
        Write-Host "Location: $installDir\$BinaryName"
        
        # Show version
        & "$installDir\$BinaryName" --version
        
        Write-Host ""
        
        # Add to PATH
        Add-ToPath $installDir
        
        # Test if available
        if (Test-PSXAvailable) {
            Write-Host ""
            Write-Host "✓ psx command is available" -ForegroundColor Green
        } else {
            Write-Host ""
            Write-Host "Note: You may need to restart your terminal" -ForegroundColor Yellow
        }
        
    } catch {
        Write-Host "Installation failed: $_" -ForegroundColor Red
        exit 1
    }
}

function Uninstall-PSX {
    $locations = @(
        "$SystemInstallDir\$BinaryName",
        "$UserInstallDir\$BinaryName"
    )
    
    $found = $false
    foreach ($location in $locations) {
        if (Test-Path $location) {
            Write-Host "Removing $location..." -ForegroundColor Yellow
            
            try {
                Remove-Item $location -Force
                $found = $true
                Write-Host "✓ Removed" -ForegroundColor Green
            } catch {
                Write-Host "Failed to remove: $_" -ForegroundColor Red
            }
        }
    }
    
    if (-not $found) {
        Write-Host "PSX is not installed"
    } else {
        Write-Host ""
        Write-Host "PSX uninstalled successfully" -ForegroundColor Green
        Write-Host ""
        Write-Host "Note: PATH entries remain. Remove them manually if needed." -ForegroundColor Yellow
    }
}

function Show-Help {
    Write-Host "Usage: .\install.ps1 {github|local|uninstall} [options]"
    Write-Host ""
    Write-Host "Commands:"
    Write-Host "  github [version]  - Install from GitHub releases (default: latest)"
    Write-Host "  local [path]      - Install from local build (default: build\psx.exe)"
    Write-Host "  uninstall         - Remove PSX from system"
    Write-Host ""
    Write-Host "Examples:"
    Write-Host "  .\install.ps1 github              # Install latest from GitHub"
    Write-Host "  .\install.ps1 github v1.0.0       # Install specific version"
    Write-Host "  .\install.ps1 local               # Install from build\psx.exe"
    Write-Host "  .\install.ps1 local C:\path\to\psx.exe  # Install from custom path"
    Write-Host "  .\install.ps1 uninstall           # Remove PSX"
    Write-Host ""
    Write-Host "Notes:"
    Write-Host "  - Run as Administrator for system-wide installation"
    Write-Host "  - Otherwise installs to user directory"
    Write-Host "  - Automatically adds to PATH"
}

# Main execution
switch ($Command) {
    'github' {
        Install-FromGitHub $Version
    }
    'local' {
        Install-FromLocal $Path
    }
    'uninstall' {
        Uninstall-PSX
    }
    'help' {
        Show-Help
    }
    default {
        Write-Host "Unknown command: $Command" -ForegroundColor Red
        Write-Host "Run '.\install.ps1 help' for usage"
        exit 1
    }
}