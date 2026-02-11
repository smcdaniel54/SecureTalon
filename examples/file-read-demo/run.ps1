# File read demo — ready-to-run. Creates work/input.txt if missing, then runs the full flow.
# Prereq: Backend running: ADDR=:8090 ADMIN_TOKEN=demo go run ./cmd/securetalon

$Base = "http://localhost:8090"
$Token = "demo"
$Headers = @{
    "Authorization" = "Bearer $Token"
    "Content-Type"  = "application/json"
}

# Ensure work/input.txt exists (from repo root)
$workDir = Join-Path (Split-Path (Split-Path $PSScriptRoot -Parent) -Parent) "work"
$inputFile = Join-Path $workDir "input.txt"
if (-not (Test-Path $inputFile)) {
    New-Item -ItemType Directory -Force -Path $workDir | Out-Null
    'hello from SecureTalon' | Set-Content $inputFile
    Write-Host "Created $inputFile"
}

Write-Host "1. Create session..."
$r = Invoke-RestMethod -Uri "$Base/v1/sessions" -Method Post -Headers $Headers -Body '{"label":"File read demo","metadata":{}}'
$SessionId = $r.id
Write-Host "   Session: $SessionId"

Write-Host "2. Set session policy (file.read under work)..."
$policyPath = Join-Path $PSScriptRoot "policy.json"
$body = Get-Content $policyPath -Raw
Invoke-RestMethod -Uri "$Base/v1/sessions/$SessionId/policy" -Method Put -Headers $Headers -Body $body | Out-Null
Write-Host "   Policy set."

Write-Host "3. Post message with file.read intent..."
$msgBody = '{"role":"user","content":"Read work/input.txt","intents":[{"tool":"file.read","params":{"path":"work/input.txt"}}]}'
$r = Invoke-RestMethod -Uri "$Base/v1/sessions/$SessionId/messages" -Method Post -Headers $Headers -Body $msgBody
$RunId = $r.run_id
Write-Host "   Run: $RunId"

Write-Host "4. Wait and get run..."
Start-Sleep -Seconds 2
$run = Invoke-RestMethod -Uri "$Base/v1/runs/$RunId" -Method Get -Headers $Headers
Write-Host "   Status: $($run.status)"
if ($run.steps) { foreach ($s in $run.steps) { Write-Host "   Step: $($s.type) $($s.status)" } }

Write-Host "Done. Open UI -> Sessions -> $SessionId -> Run $RunId"
Write-Host "Policy used: $policyPath"
