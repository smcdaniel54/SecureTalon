# Create a distribution zip for public release (no .git, no secrets, no build artifacts).
# Run from repo root: .\scripts\make-dist-zip.ps1
# Output: SecureTalon-<date>.zip in the current directory.

$ErrorActionPreference = "Stop"
$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
$date = Get-Date -Format "yyyyMMdd"
$zipName = "SecureTalon-$date.zip"
$tempDir = Join-Path $env:TEMP "SecureTalon-dist-$date"

if (Test-Path $tempDir) { Remove-Item -Recurse -Force $tempDir }
New-Item -ItemType Directory -Path $tempDir | Out-Null

# Copy repo excluding .git and other non-distributable paths
$exclude = @(
    ".git",
    "node_modules",
    "ui/node_modules",
    "ui/dist",
    "data",
    ".env",
    "*.local",
    "*.exe",
    "*.test",
    "*.out",
    "vendor"
)

Get-ChildItem -Path $repoRoot -Force | Where-Object {
    $name = $_.Name
    $excluded = $false
    foreach ($e in $exclude) {
        if ($e -like "**") {
            if ($name -like $e) { $excluded = $true; break }
        } else {
            if ($name -eq $e) { $excluded = $true; break }
        }
    }
    -not $excluded
} | ForEach-Object {
    $dest = Join-Path $tempDir $_.Name
    if ($_.PSIsContainer) {
        Copy-Item -Path $_.FullName -Destination $dest -Recurse -Force
    } else {
        Copy-Item -Path $_.FullName -Destination $dest -Force
    }
}

# Remove .git if any copy slipped (e.g. nested)
Get-ChildItem -Path $tempDir -Recurse -Force -Filter ".git" -Directory -ErrorAction SilentlyContinue | Remove-Item -Recurse -Force

$zipPath = Join-Path (Get-Location) $zipName
if (Test-Path $zipPath) { Remove-Item $zipPath -Force }
Compress-Archive -Path (Join-Path $tempDir "*") -DestinationPath $zipPath -CompressionLevel Optimal
Remove-Item -Recurse -Force $tempDir

Write-Host "Created: $zipPath (no .git, no node_modules/data)"
Write-Host "Ready for public release."
