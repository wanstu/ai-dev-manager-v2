param(
    [Alias('o')]
    [string]$OutputName = 'adm-desktop-windows-amd64.exe',
    [string]$OutputDir = '',
    [switch]$clean,
    [switch]$trimpath
)

$ErrorActionPreference = 'Stop'

$repoRoot = Split-Path -Parent $PSScriptRoot
$desktopRoot = Join-Path $repoRoot 'cmd\ai-dev-manager-desktop'
if ([string]::IsNullOrWhiteSpace($OutputDir)) {
    $OutputDir = Join-Path $repoRoot 'dist'
} elseif (-not [System.IO.Path]::IsPathRooted($OutputDir)) {
    $OutputDir = Join-Path $repoRoot $OutputDir
}

$wailsArgs = @('build', '-skipbindings')
if ($clean) { $wailsArgs += '-clean' }
if ($trimpath) { $wailsArgs += '-trimpath' }
$wailsArgs += @('-o', $OutputName)

Push-Location $desktopRoot
try {
    # wails.json owns Desktop pre-build preparation, including icon refresh.
    & go run github.com/wailsapp/wails/v2/cmd/wails@v2.15.0 @wailsArgs
    if ($LASTEXITCODE -ne 0) {
        exit $LASTEXITCODE
    }
} finally {
    Pop-Location
}

$builtDesktop = Join-Path $desktopRoot (Join-Path 'build\bin' $OutputName)
if (-not (Test-Path -LiteralPath $builtDesktop)) {
    throw "Wails Desktop output not found: $builtDesktop"
}
New-Item -ItemType Directory -Force -Path $OutputDir | Out-Null
$finalDesktop = Join-Path $OutputDir $OutputName
Copy-Item -LiteralPath $builtDesktop -Destination $finalDesktop -Force
Write-Host "Desktop artifact: $finalDesktop"
