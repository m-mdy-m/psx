# PSX Build Script for Windows
# Builds binaries for Windows or all platforms

param(
    [Parameter(Position=0)]
    [ValidateSet('current', 'all', 'release', 'clean')]
    [string]$Command = 'current'
)

# Configuration
$Version = if ($env:VERSION) { $env:VERSION } else { 
    try {
        $gitVersion = git describe --tags --always --dirty 2>$null
        if ($gitVersion) { $gitVersion } else { "dev" }
    } catch {
        "dev"
    }
}
$BuildDate = Get-Date -Format "yyyy-MM-dd_HH:mm:ss" -AsUTC
$BinaryName = "psx"
$BuildDir = "build"
$CmdDir = ".\cmd\psx"

# LDFLAGS
$LDFlags = "-s -w -X main.Version=$Version -X main.BuildDate=$BuildDate"

Write-Host "PSX Build Script" -ForegroundColor Blue
Write-Host "================" -ForegroundColor Blue
Write-Host "Version: $Version"
Write-Host "Date: $BuildDate"
Write-Host ""

function Build-Platform {
    param(
        [string]$Os,
        [string]$Arch,
        [string]$Output
    )
    
    Write-Host "Building for $Os/$Arch..." -ForegroundColor Yellow
    
    $env:GOOS = $Os
    $env:GOARCH = $Arch
    $env:CGO_ENABLED = "0"
    
    $buildArgs = @(
        "build",
        "-ldflags", $LDFlags,
        "-o", $Output,
        $CmdDir
    )
    
    & go @buildArgs
    
    if ($LASTEXITCODE -eq 0) {
        $size = (Get-Item $Output).Length / 1MB
        Write-Host "✓ Built: $Output ($([math]::Round($size, 2)) MB)" -ForegroundColor Green
        return $true
    } else {
        Write-Host "✗ Build failed for $Os/$Arch" -ForegroundColor Red
        return $false
    }
}

switch ($Command) {
    'current' {
        Write-Host "Building for current platform..."
        
        if (-not (Test-Path $BuildDir)) {
            New-Item -ItemType Directory -Path $BuildDir | Out-Null
        }
        
        $os = if ($IsLinux) { "linux" } elseif ($IsMacOS) { "darwin" } else { "windows" }
        $arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
        $ext = if ($os -eq "windows") { ".exe" } else { "" }
        
        $success = Build-Platform $os $arch "$BuildDir\$BinaryName$ext"
        
        if ($success) {
            Write-Host ""
            Write-Host "Build complete!" -ForegroundColor Green
            Write-Host "Binary: $BuildDir\$BinaryName$ext"
        } else {
            exit 1
        }
    }
    
    'all' {
        Write-Host "Building for all platforms..."
        
        if (-not (Test-Path $BuildDir)) {
            New-Item -ItemType Directory -Path $BuildDir | Out-Null
        }
        
        $platforms = @(
            @{Os="linux"; Arch="amd64"; Output="$BuildDir\$BinaryName-linux-amd64"},
            @{Os="linux"; Arch="arm64"; Output="$BuildDir\$BinaryName-linux-arm64"},
            @{Os="darwin"; Arch="amd64"; Output="$BuildDir\$BinaryName-darwin-amd64"},
            @{Os="darwin"; Arch="arm64"; Output="$BuildDir\$BinaryName-darwin-arm64"},
            @{Os="windows"; Arch="amd64"; Output="$BuildDir\$BinaryName-windows-amd64.exe"}
        )
        
        $allSuccess = $true
        foreach ($platform in $platforms) {
            $success = Build-Platform $platform.Os $platform.Arch $platform.Output
            if (-not $success) {
                $allSuccess = $false
            }
        }
        
        if ($allSuccess) {
            Write-Host ""
            Write-Host "Generating checksums..." -ForegroundColor Yellow
            
            Get-ChildItem $BuildDir\$BinaryName-* | ForEach-Object {
                $hash = (Get-FileHash $_.FullName -Algorithm SHA256).Hash.ToLower()
                "$hash  $($_.Name)"
            } | Out-File "$BuildDir\SHA256SUMS" -Encoding utf8
            
            Write-Host ""
            Write-Host "All builds complete!" -ForegroundColor Green
            Write-Host "Binaries in: $BuildDir\"
        } else {
            Write-Host ""
            Write-Host "Some builds failed" -ForegroundColor Red
            exit 1
        }
    }
    
    'release' {
        Write-Host "Building release..."
        
        # Clean first
        if (Test-Path $BuildDir) {
            Remove-Item $BuildDir -Recurse -Force
        }
        New-Item -ItemType Directory -Path $BuildDir | Out-Null
        
        # Run tests
        Write-Host "Running tests..." -ForegroundColor Yellow
        go test .\... -v
        
        if ($LASTEXITCODE -ne 0) {
            Write-Host "Tests failed! Aborting release build." -ForegroundColor Red
            exit 1
        }
        
        # Build all platforms
        & $PSCommandPath -Command all
        
        if ($LASTEXITCODE -ne 0) {
            exit 1
        }
        
        # Create archives
        Write-Host ""
        Write-Host "Creating archives..." -ForegroundColor Yellow
        
        Push-Location $BuildDir
        
        # Compress files
        Compress-Archive -Path "$BinaryName-linux-amd64" -DestinationPath "$BinaryName-$Version-linux-amd64.zip"
        Compress-Archive -Path "$BinaryName-linux-arm64" -DestinationPath "$BinaryName-$Version-linux-arm64.zip"
        Compress-Archive -Path "$BinaryName-darwin-amd64" -DestinationPath "$BinaryName-$Version-darwin-amd64.zip"
        Compress-Archive -Path "$BinaryName-darwin-arm64" -DestinationPath "$BinaryName-$Version-darwin-arm64.zip"
        Compress-Archive -Path "$BinaryName-windows-amd64.exe" -DestinationPath "$BinaryName-$Version-windows-amd64.zip"
        
        # Update checksums
        Get-ChildItem *.zip | ForEach-Object {
            $hash = (Get-FileHash $_.FullName -Algorithm SHA256).Hash.ToLower()
            "$hash  $($_.Name)"
        } | Out-File "SHA256SUMS" -Encoding utf8
        
        Pop-Location
        
        Write-Host ""
        Write-Host "Release build complete!" -ForegroundColor Green
        Write-Host "Archives in: $BuildDir\"
    }
    
    'clean' {
        Write-Host "Cleaning build directory..."
        
        if (Test-Path $BuildDir) {
            Remove-Item $BuildDir -Recurse -Force
            Write-Host "Clean complete!" -ForegroundColor Green
        } else {
            Write-Host "Build directory does not exist"
        }
    }
}