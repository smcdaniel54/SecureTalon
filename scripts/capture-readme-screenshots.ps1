# Capture real README screenshots. Requires: Go, Node, npm, Playwright (npx playwright install chromium in ui/).
# Run from repo root. Backend and UI are started temporarily; screenshots go to docs/screenshots/.

$ErrorActionPreference = "Stop"
$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
$uiDir = Join-Path $repoRoot "ui"
$backendPort = 8090
$uiPort = 5173

# Find a free port for UI (Vite may use 5174, 5175, etc. if 5173 is taken)
$env:ADDR = ":$backendPort"
$env:ADMIN_TOKEN = "demo"

Write-Host "Starting backend on port $backendPort..."
$backendJob = Start-Job -ScriptBlock {
  Set-Location $using:repoRoot
  $env:ADDR = $using:env:ADDR
  $env:ADMIN_TOKEN = $using:env:ADMIN_TOKEN
  go run ./cmd/securetalon 2>&1
}
Start-Sleep -Seconds 5

Write-Host "Starting UI (port may be $uiPort or next available)..."
$uiJob = Start-Job -ScriptBlock {
  Set-Location $using:uiDir
  npm run dev 2>&1
}
Start-Sleep -Seconds 8

# Detect UI port from job output (Vite: "Local: http://localhost:5173")
$uiOutput = Receive-Job $uiJob
if ($uiOutput -is [array]) { $uiOutput = $uiOutput -join "`n" }
$uiUrl = "http://localhost:$uiPort"
if ($uiOutput -and ($uiOutput -match "localhost:(\d+)")) {
  $uiPort = [int]$Matches[1]
  $uiUrl = "http://localhost:$uiPort"
}
Write-Host "UI at $uiUrl"

Write-Host "Capturing screenshots..."
Push-Location $uiDir
try {
  $env:UI_BASE_URL = $uiUrl
  $env:API_BASE_URL = "http://localhost:$backendPort"
  $env:ADMIN_TOKEN = "demo"
  npm run capture-screenshots
  if ($LASTEXITCODE -ne 0) { throw "capture-screenshots failed" }
} finally {
  Pop-Location
}

Stop-Job $backendJob, $uiJob -ErrorAction SilentlyContinue
Remove-Job $backendJob, $uiJob -Force -ErrorAction SilentlyContinue

Write-Host "Screenshots saved to docs/screenshots/"
Get-ChildItem (Join-Path $repoRoot "docs\screenshots\*.png") | Where-Object { $_.Name -notlike "_*" } | ForEach-Object { Write-Host "  $($_.Name) $([math]::Round($_.Length/1KB, 1)) KB" }
