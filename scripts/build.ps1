# PSX Build Script for Windows

param(
    [Parameter(Position=0)]
    [ValidateSet('current', 'all', 'release', 'clean')]
    [string]$Command = 'current'
)

$ErrorActionPreference = 'Stop'

$Version = if ($env:VERSION) {
    $env:VERSION
} else {
    try {
        $gitVersion = git describe --tags --always --dirty 2>$null
        if ($gitVersion) { $gitVersion } else { 'dev' }
    } catch {
        'dev'
    }
}

$BuildDate = (Get-Date).ToUniversalTime().ToString('yyyy-MM-dd_HH:mm:ss')
$BinaryName = 'psx'
$BuildDir = 'build'
$CmdDir = '.\cmd\psx'

# Version is defined in internal/command, matching the release workflow.
$LDFlags = "-s -w -X github.com/m-mdy-m/psx/internal/command.Version=$Version"

function Build-Platform {
    param(
        [string]$Os,
        [string]$Arch,
        [string]$Output
    )

    Write-Host "Building for $Os/$Arch..." -ForegroundColor Yellow

    $env:GOOS = $Os
    $env:GOARCH = $Arch
    $env:CGO_ENABLED = '0'

    & go build -trimpath -ldflags $LDFlags -o $Output $CmdDir
    if ($LASTEXITCODE -ne 0) {
        throw "Build failed for $Os/$Arch"
    }

    $size = (Get-Item -LiteralPath $Output).Length / 1MB
    Write-Host "✓ Built: $Output ($([math]::Round($size, 2)) MB)" -ForegroundColor Green
}

switch ($Command) {
    'current' {
        if (-not (Test-Path -LiteralPath $BuildDir)) {
            New-Item -ItemType Directory -Path $BuildDir | Out-Null
        }

        $os = if ($IsLinux) { 'linux' } elseif ($IsMacOS) { 'darwin' } else { 'windows' }
        $arch = if ([Environment]::Is64BitOperatingSystem) { 'amd64' } else { '386' }
        $ext = if ($os -eq 'windows') { '.exe' } else { '' }

        Build-Platform $os $arch "$BuildDir\$BinaryName$ext"
    }

    'all' {
        if (-not (Test-Path -LiteralPath $BuildDir)) {
            New-Item -ItemType Directory -Path $BuildDir | Out-Null
        }

        $platforms = @(
            @{ Os='linux';   Arch='amd64'; Output="$BuildDir\$BinaryName-linux-x64" },
            @{ Os='linux';   Arch='arm64'; Output="$BuildDir\$BinaryName-linux-arm64" },
            @{ Os='darwin';  Arch='amd64'; Output="$BuildDir\$BinaryName-darwin-x64" },
            @{ Os='darwin';  Arch='arm64'; Output="$BuildDir\$BinaryName-darwin-arm64" },
            @{ Os='windows'; Arch='amd64'; Output="$BuildDir\$BinaryName-windows-x64.exe" }
        )

        foreach ($platform in $platforms) {
            Build-Platform $platform.Os $platform.Arch $platform.Output
        }

        Get-ChildItem -LiteralPath $BuildDir -Filter "$BinaryName-*" -File |
            Where-Object { $_.Name -notmatch '\.(zip|tar\.gz|sha256|txt)$' } |
            ForEach-Object {
                $hash = (Get-FileHash -LiteralPath $_.FullName -Algorithm SHA256).Hash.ToLowerInvariant()
                "$hash  $($_.Name)"
            } | Set-Content -LiteralPath "$BuildDir\checksums.txt" -Encoding ascii

        Write-Host "✓ All builds complete" -ForegroundColor Green
    }

    'release' {
        if (Test-Path -LiteralPath $BuildDir) {
            Remove-Item -LiteralPath $BuildDir -Recurse -Force
        }
        New-Item -ItemType Directory -Path $BuildDir | Out-Null

        Write-Host 'Running tests...' -ForegroundColor Yellow
        & go test .\...
        if ($LASTEXITCODE -ne 0) {
            throw 'Tests failed! Aborting release build.'
        }

        & $PSCommandPath -Command all
        if ($LASTEXITCODE -ne 0) {
            throw 'Release build failed.'
        }

        Push-Location $BuildDir
        try {
            if (Get-Command tar -ErrorAction SilentlyContinue) {
                & tar -czf "$BinaryName-$Version-linux-x64.tar.gz" "$BinaryName-linux-x64"
                & tar -czf "$BinaryName-$Version-linux-arm64.tar.gz" "$BinaryName-linux-arm64"
                & tar -czf "$BinaryName-$Version-darwin-x64.tar.gz" "$BinaryName-darwin-x64"
                & tar -czf "$BinaryName-$Version-darwin-arm64.tar.gz" "$BinaryName-darwin-arm64"
            } else {
                throw 'tar is required to create Unix release archives on Windows.'
            }

            Compress-Archive -Path "$BinaryName-windows-x64.exe" -DestinationPath "$BinaryName-$Version-windows-x64.zip" -Force
        }
        finally {
            Pop-Location
        }

        Write-Host '✓ Release build complete' -ForegroundColor Green
    }

    'clean' {
        if (Test-Path -LiteralPath $BuildDir) {
            Remove-Item -LiteralPath $BuildDir -Recurse -Force
        }
        Write-Host '✓ Clean complete' -ForegroundColor Green
    }
}
