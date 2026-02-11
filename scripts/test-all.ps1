# Tight loop: run backend tests + UI build, print logs, exit non-zero if any fail.
# Usage: .\scripts\test-all.ps1   (from repo root)

$ErrorActionPreference = "Stop"
$repoRoot = Split-Path $PSScriptRoot -Parent
if (-not (Test-Path (Join-Path $repoRoot "go.mod"))) {
  Write-Host "Cannot find repo root (go.mod) from $PSScriptRoot" -ForegroundColor Red
  exit 1
}

$failed = $false

Write-Host "=== Backend tests (go test ./...) ===" -ForegroundColor Cyan
Push-Location $repoRoot
try {
  $goOut = go test ./... 2>&1
  $goOut | ForEach-Object { Write-Host $_ }
  if ($LASTEXITCODE -ne 0) {
    $failed = $true
    Write-Host "Backend tests FAILED" -ForegroundColor Red
  } else {
    Write-Host "Backend tests OK" -ForegroundColor Green
  }
} catch {
  $failed = $true
  Write-Host "Backend tests error: $_" -ForegroundColor Red
} finally {
  Pop-Location
}

Write-Host ""
Write-Host "=== UI build (npm run build) ===" -ForegroundColor Cyan
Push-Location (Join-Path $repoRoot "ui")
try {
  $npmOut = npm run build 2>&1
  $npmOut | ForEach-Object { Write-Host $_ }
  if ($LASTEXITCODE -ne 0) {
    $failed = $true
    Write-Host "UI build FAILED" -ForegroundColor Red
  } else {
    Write-Host "UI build OK" -ForegroundColor Green
  }
} catch {
  $failed = $true
  Write-Host "UI build error: $_" -ForegroundColor Red
} finally {
  Pop-Location
}

Write-Host ""
if ($failed) {
  Write-Host "One or more steps failed." -ForegroundColor Red
  exit 1
}
Write-Host "All checks passed." -ForegroundColor Green
exit 0
