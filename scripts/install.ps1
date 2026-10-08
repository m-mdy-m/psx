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

$ErrorActionPreference = 'Stop'

$BinaryName = 'psx.exe'
$Repo = 'm-mdy-m/psx'
$SystemInstallDir = Join-Path $env:ProgramFiles 'PSX'
$UserInstallDir = Join-Path $env:LOCALAPPDATA 'PSX'

function Fail {
    param([string]$Message)
    throw $Message
}

function Test-Administrator {
    $identity = [Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = [Security.Principal.WindowsPrincipal]::new($identity)
    return $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

function Get-InstallDirectory {
    if (Test-Administrator) {
        return $SystemInstallDir
    }
    return $UserInstallDir
}

function Get-LatestVersion {
    try {
        $response = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest" -Method Get
    } catch {
        Fail "Could not fetch latest release information: $($_.Exception.Message)"
    }

    if (-not $response.tag_name) {
        Fail 'Could not determine latest version.'
    }

    return [string]$response.tag_name
}

function Normalize-Version {
    param([string]$Value)
    if ($Value.StartsWith('v')) {
        return $Value
    }
    return "v$Value"
}

function Get-Sha256 {
    param([string]$FilePath)
    return (Get-FileHash -LiteralPath $FilePath -Algorithm SHA256).Hash.ToLowerInvariant()
}

function Download-File {
    param(
        [string]$Uri,
        [string]$OutFile
    )

    try {
        Invoke-WebRequest -Uri $Uri -OutFile $OutFile -UseBasicParsing
    } catch {
        Fail "Download failed: $($_.Exception.Message)"
    }
}

function Verify-Binary {
    param(
        [string]$BinaryPath,
        [string]$BinaryName,
        [string]$ChecksumsPath
    )

    $line = Get-Content -LiteralPath $ChecksumsPath | Where-Object {
        $_ -match "\s+$([regex]::Escape($BinaryName))$"
    } | Select-Object -First 1

    if (-not $line) {
        Fail "No checksum found for $BinaryName"
    }

    $expected = ($line -split '\s+')[0].ToLowerInvariant()
    $actual = Get-Sha256 $BinaryPath

    if ($actual -ne $expected) {
        Fail "Checksum verification failed for $BinaryName"
    }

    Write-Host '✓ Checksum verified' -ForegroundColor Green
}

function Add-ToPath {
    param([string]$Directory)

    $pathScope = if (Test-Administrator) { 'Machine' } else { 'User' }
    $currentPath = [Environment]::GetEnvironmentVariable('Path', $pathScope)

    $parts = @()
    if ($currentPath) {
        $parts = $currentPath -split ';' | Where-Object { $_ -and $_.Trim() }
    }

    $alreadyPresent = $parts | Where-Object { $_.TrimEnd('\\') -ieq $Directory.TrimEnd('\\') }

    if (-not $alreadyPresent) {
        $newPath = (($parts + $Directory) -join ';')
        [Environment]::SetEnvironmentVariable('Path', $newPath, $pathScope)
    }

    # Make the command available in the current PowerShell process too.
    $sessionParts = @()
    if ($env:Path) {
        $sessionParts = $env:Path -split ';' | Where-Object { $_ -and $_.Trim() }
    }

    $sessionPresent = $sessionParts | Where-Object { $_.TrimEnd('\\') -ieq $Directory.TrimEnd('\\') }
    if (-not $sessionPresent) {
        $env:Path = (($sessionParts + $Directory) -join ';')
    }

    Write-Host "✓ Added to PATH ($pathScope)" -ForegroundColor Green
}

function Install-Binary {
    param([string]$SourcePath)

    $installDir = Get-InstallDirectory
    $destination = Join-Path $installDir $BinaryName

    if (-not (Test-Path -LiteralPath $installDir)) {
        New-Item -ItemType Directory -Path $installDir -Force | Out-Null
    }

    Copy-Item -LiteralPath $SourcePath -Destination $destination -Force

    Write-Host ''
    Write-Host '✓ PSX installed successfully!' -ForegroundColor Green
    Write-Host "Location: $destination"

    try {
        & $destination --version
        if ($LASTEXITCODE -ne 0) {
            Fail 'Installed binary returned a non-zero exit code.'
        }
    } catch {
        Fail "Installed binary could not be executed: $($_.Exception.Message)"
    }

    Add-ToPath $installDir

    Write-Host ''
    Write-Host 'Try: psx --version' -ForegroundColor Yellow
}

function Install-FromGitHub {
    param([string]$RequestedVersion)

    Write-Host 'Platform: Windows (x64)'
    Write-Host "Version: $RequestedVersion"
    Write-Host ''

    if ([Environment]::Is64BitOperatingSystem -eq $false) {
        Fail '32-bit Windows is not supported.'
    }

    $version = if ($RequestedVersion -eq 'latest') {
        Write-Host 'Fetching latest version...' -ForegroundColor Yellow
        Get-LatestVersion
    } else {
        Normalize-Version $RequestedVersion
    }

    Write-Host "Release: $version"

    # The Windows release publishes a raw EXE as well as a ZIP. The raw EXE is
    # intentionally used here because it is directly executable and its SHA256
    # is present in checksums.txt.
    $assetName = 'psx-windows-x64.exe'
    $assetUrl = "https://github.com/$Repo/releases/download/$version/$assetName"
    $checksumsUrl = "https://github.com/$Repo/releases/download/$version/checksums.txt"
    $tempDir = Join-Path $env:TEMP ("psx-install-{0}" -f ([Guid]::NewGuid().ToString('N')))
    $binaryPath = Join-Path $tempDir $assetName
    $checksumsPath = Join-Path $tempDir 'checksums.txt'

    try {
        New-Item -ItemType Directory -Path $tempDir -Force | Out-Null

        Write-Host "Downloading $assetName..." -ForegroundColor Yellow
        Download-File $assetUrl $binaryPath

        if (-not (Test-Path -LiteralPath $binaryPath)) {
            Fail 'Downloaded binary was not created.'
        }

        if ((Get-Item -LiteralPath $binaryPath).Length -lt 100KB) {
            Fail 'Downloaded file is unexpectedly small; refusing to install it.'
        }

        $bytes = [System.IO.File]::ReadAllBytes($binaryPath)
        if ($bytes.Length -lt 2 -or $bytes[0] -ne 0x4D -or $bytes[1] -ne 0x5A) {
            Fail 'Downloaded file is not a valid Windows PE executable.'
        }

        Write-Host 'Downloading checksums...' -ForegroundColor Yellow
        Download-File $checksumsUrl $checksumsPath
        Verify-Binary $binaryPath $assetName $checksumsPath

        Install-Binary $binaryPath
    }
    finally {
        if (Test-Path -LiteralPath $tempDir) {
            Remove-Item -LiteralPath $tempDir -Recurse -Force -ErrorAction SilentlyContinue
        }
    }
}

function Install-FromLocal {
    param([string]$BinaryPath)

    if ([string]::IsNullOrWhiteSpace($BinaryPath)) {
        $BinaryPath = 'build\psx.exe'
    }

    if (-not (Test-Path -LiteralPath $BinaryPath)) {
        Fail "Binary not found: $BinaryPath`nBuild first: .\scripts\build.ps1"
    }

    Write-Host "Installing from: $BinaryPath"
    Install-Binary (Resolve-Path -LiteralPath $BinaryPath).Path
}

function Uninstall-PSX {
    $locations = @(
        (Join-Path $SystemInstallDir $BinaryName),
        (Join-Path $UserInstallDir $BinaryName)
    )

    $found = $false
    foreach ($location in $locations) {
        if (Test-Path -LiteralPath $location) {
            try {
                Remove-Item -LiteralPath $location -Force
                Write-Host "✓ Removed $location" -ForegroundColor Green
                $found = $true
            } catch {
                Fail "Failed to remove $location: $($_.Exception.Message)"
            }
        }
    }

    if (-not $found) {
        Write-Host 'PSX is not installed.'
    } else {
        Write-Host 'PSX uninstalled successfully.' -ForegroundColor Green
        Write-Host 'PATH entries are intentionally left untouched.' -ForegroundColor Yellow
    }
}

function Show-Help {
    Write-Host 'Usage: install.ps1 {github|local|uninstall|help} [version/path]'
    Write-Host ''
    Write-Host 'Examples:'
    Write-Host '  irm https://raw.githubusercontent.com/m-mdy-m/psx/main/scripts/install.ps1 | iex'
    Write-Host '  .\scripts\install.ps1 github v0.3.0'
    Write-Host '  .\scripts\install.ps1 local'
    Write-Host '  .\scripts\install.ps1 uninstall'
    Write-Host ''
    Write-Host 'The installer never calls exit so it is safe to run through Invoke-Expression.'
}

Write-Host 'PSX Installation Script' -ForegroundColor Blue
Write-Host '=======================' -ForegroundColor Blue
Write-Host ''

switch ($Command) {
    'github'   { Install-FromGitHub $Version }
    'local'    { Install-FromLocal $Path }
    'uninstall'{ Uninstall-PSX }
    'help'     { Show-Help }
    default    { Fail "Unknown command: $Command" }
}
