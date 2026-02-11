# HTTP fetch demo — ready-to-run. Uses policy.json and sends sample intent.
# Prereq: Backend running: ADDR=:8090 ADMIN_TOKEN=demo go run ./cmd/securetalon

$Base = "http://localhost:8090"
$Token = "demo"
$Headers = @{
    "Authorization" = "Bearer $Token"
    "Content-Type"  = "application/json"
}

Write-Host "1. Create session..."
$r = Invoke-RestMethod -Uri "$Base/v1/sessions" -Method Post -Headers $Headers -Body '{"label":"HTTP fetch demo","metadata":{}}'
$SessionId = $r.id
Write-Host "   Session: $SessionId"

Write-Host "2. Set session policy (http.fetch to jsonplaceholder.typicode.com, GET only)..."
$policyPath = Join-Path $PSScriptRoot "policy.json"
$body = Get-Content $policyPath -Raw
Invoke-RestMethod -Uri "$Base/v1/sessions/$SessionId/policy" -Method Put -Headers $Headers -Body $body | Out-Null
Write-Host "   Policy set."

Write-Host "3. Post message with http.fetch intent..."
$msgBody = Get-Content (Join-Path $PSScriptRoot "message.json") -Raw
$r = Invoke-RestMethod -Uri "$Base/v1/sessions/$SessionId/messages" -Method Post -Headers $Headers -Body $msgBody
$RunId = $r.run_id
Write-Host "   Run: $RunId"

Write-Host "4. Wait and get run..."
Start-Sleep -Seconds 3
$run = Invoke-RestMethod -Uri "$Base/v1/runs/$RunId" -Method Get -Headers $Headers
Write-Host "   Status: $($run.status)"
if ($run.steps) { foreach ($s in $run.steps) { Write-Host "   Step: $($s.type) $($s.status)" } }

Write-Host "Done. Open UI -> Sessions -> $SessionId -> Run $RunId"
Write-Host "Policy: $policyPath"
