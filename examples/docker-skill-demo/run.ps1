# Docker skill demo — ready-to-run. Set $ImageDigest or edit policy.json + message.json.
# Prereq: Backend with Docker: ADDR=:8090 ADMIN_TOKEN=demo go run ./cmd/securetalon

$Base = "http://localhost:8090"
$Token = "demo"
$Headers = @{
    "Authorization" = "Bearer $Token"
    "Content-Type"  = "application/json"
}

$policyPath = Join-Path $PSScriptRoot "policy.json"
$messagePath = Join-Path $PSScriptRoot "message.json"
$policyBody = Get-Content $policyPath -Raw
$messageBody = Get-Content $messagePath -Raw

# Replace placeholder if IMAGE_DIGEST env is set; otherwise require user to have edited the JSON files
if ($env:IMAGE_DIGEST) {
    $policyBody = $policyBody -replace 'YOUR_IMAGE_DIGEST', $env:IMAGE_DIGEST
    $messageBody = $messageBody -replace 'YOUR_IMAGE_DIGEST', $env:IMAGE_DIGEST
} elseif ($policyBody -match 'YOUR_IMAGE_DIGEST') {
    Write-Host "Set IMAGE_DIGEST env or replace YOUR_IMAGE_DIGEST in policy.json and message.json"
    Write-Host "Example: docker inspect hello-world --format '{{index .RepoDigests 0}}'"
    exit 1
}

Write-Host "1. Create session..."
$r = Invoke-RestMethod -Uri "$Base/v1/sessions" -Method Post -Headers $Headers -Body '{"label":"Docker skill demo","metadata":{}}'
$SessionId = $r.id
Write-Host "   Session: $SessionId"

Write-Host "2. Set session policy (docker.run for image digest)..."
Invoke-RestMethod -Uri "$Base/v1/sessions/$SessionId/policy" -Method Put -Headers $Headers -Body $policyBody | Out-Null
Write-Host "   Policy set."

Write-Host "3. Post message with docker.run intent..."
$r = Invoke-RestMethod -Uri "$Base/v1/sessions/$SessionId/messages" -Method Post -Headers $Headers -Body $messageBody
$RunId = $r.run_id
Write-Host "   Run: $RunId"

Write-Host "4. Wait and get run..."
Start-Sleep -Seconds 5
$run = Invoke-RestMethod -Uri "$Base/v1/runs/$RunId" -Method Get -Headers $Headers
Write-Host "   Status: $($run.status)"
if ($run.steps) { foreach ($s in $run.steps) { Write-Host "   Step: $($s.type) $($s.status)" } }

Write-Host "Done. Open UI -> Sessions -> $SessionId -> Run $RunId"
