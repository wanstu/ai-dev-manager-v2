param(
    [Parameter(Mandatory = $true)]
    [string]$Version
)

$ErrorActionPreference = 'Stop'

if ([string]::IsNullOrWhiteSpace($Version)) {
    throw 'Version is required, for example v1.0.0-rc.2'
}

$repoRoot = Split-Path -Parent $PSScriptRoot
$distRoot = Join-Path $repoRoot 'dist'
$desktopBinRoot = Join-Path $repoRoot 'cmd\ai-dev-manager-desktop\build\bin'
$cliName = "ai-dev-manager-v2-$Version-windows-amd64.exe"
$desktopName = "ai-dev-manager-v2-desktop-$Version-windows-amd64.exe"
$cliPath = Join-Path $distRoot $cliName
$desktopPath = Join-Path $distRoot $desktopName
$checksumPath = Join-Path $distRoot "SHA256SUMS-$Version.txt"

New-Item -ItemType Directory -Force -Path $distRoot | Out-Null

Push-Location $repoRoot
try {
    & go build -trimpath -o $cliPath ./cmd/ai-dev-manager
    if ($LASTEXITCODE -ne 0) {
        exit $LASTEXITCODE
    }

    & (Join-Path $PSScriptRoot 'build-desktop.ps1') -clean -trimpath -o $desktopName
    if ($LASTEXITCODE -ne 0) {
        exit $LASTEXITCODE
    }

    $builtDesktop = Join-Path $desktopBinRoot $desktopName
    if (-not (Test-Path -LiteralPath $builtDesktop)) {
        throw "Wails Desktop output not found: $builtDesktop"
    }
    Copy-Item -LiteralPath $builtDesktop -Destination $desktopPath -Force

    function Get-Sha256Hex([string]$Path) {
        $sha = [System.Security.Cryptography.SHA256]::Create()
        try {
            $stream = [System.IO.File]::OpenRead($Path)
            try {
                $hash = $sha.ComputeHash($stream)
            } finally {
                $stream.Dispose()
            }
        } finally {
            $sha.Dispose()
        }
        return ([System.BitConverter]::ToString($hash)).Replace('-', '').ToLowerInvariant()
    }

    $lines = @(
        "$(Get-Sha256Hex $cliPath)  $cliName",
        "$(Get-Sha256Hex $desktopPath)  $desktopName"
    )
    [System.IO.File]::WriteAllLines($checksumPath, $lines, (New-Object System.Text.UTF8Encoding($false)))

    Write-Host "RC artifacts:"
    Write-Host "  $cliPath"
    Write-Host "  $desktopPath"
    Write-Host "  $checksumPath"
} finally {
    Pop-Location
}
