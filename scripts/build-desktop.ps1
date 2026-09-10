$ErrorActionPreference = 'Stop'

$repoRoot = Split-Path -Parent $PSScriptRoot
$desktopRoot = Join-Path $repoRoot 'cmd\ai-dev-manager-desktop'
$desktopBuildRoot = Join-Path $desktopRoot 'build'
$appIconSource = Join-Path $repoRoot 'assets\icons\ai-dev-manager-app.png'
$appIconTarget = Join-Path $desktopBuildRoot 'appicon.png'
$windowsIconTarget = Join-Path $desktopBuildRoot 'windows\icon.ico'

if (-not (Test-Path $appIconSource)) {
    throw "Desktop app icon source not found: $appIconSource"
}
New-Item -ItemType Directory -Force -Path $desktopBuildRoot | Out-Null
Copy-Item $appIconSource $appIconTarget -Force
# Wails only regenerates build/windows/icon.ico when it is absent. Remove the
# ignored generated file so every build reflects the committed app icon source.
if (Test-Path $windowsIconTarget) {
    Remove-Item $windowsIconTarget -Force
}

Push-Location $desktopRoot
try {
    go run github.com/wailsapp/wails/v2/cmd/wails@v2.15.0 build -skipbindings @args
    if ($LASTEXITCODE -ne 0) {
        exit $LASTEXITCODE
    }
} finally {
    Pop-Location
}
