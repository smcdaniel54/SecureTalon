# Deny then fix — trigger a deny, then add policy override and retry.
# Prereq: Backend: ADDR=:8090 ADMIN_TOKEN=demo. Create work/input.txt (e.g. 'hello' | Set-Content work/input.txt).

$Base = "http://localhost:8090"
$Token = "demo"
$Headers = @{
    "Authorization" = "Bearer $Token"
    "Content-Type"  = "application/json"
}

Write-Host "1. Create session (no file.read policy yet)..."
$r = Invoke-RestMethod -Uri "$Base/v1/sessions" -Method Post -Headers $Headers -Body '{"label":"Deny then fix","metadata":{}}'
$SessionId = $r.id
Write-Host "   Session: $SessionId"

Write-Host "2. Post message with file.read intent (expect DENY)..."
$msgBody = '{"role":"user","content":"Read work/input.txt","intents":[{"tool":"file.read","params":{"path":"work/input.txt"}}]}'
$r = Invoke-RestMethod -Uri "$Base/v1/sessions/$SessionId/messages" -Method Post -Headers $Headers -Body $msgBody
$RunId1 = $r.run_id
Write-Host "   Run: $RunId1"

Start-Sleep -Seconds 2
$run1 = Invoke-RestMethod -Uri "$Base/v1/runs/$RunId1" -Method Get -Headers $Headers
Write-Host "   Status: $($run1.status)"
foreach ($s in $run1.steps) {
    Write-Host "   Step: $($s.type) $($s.status)"
    if ($s.status -eq "denied" -and $s.details) { Write-Host "     Reason: $($s.details | ConvertTo-Json -Compress)" }
}

Write-Host "3. Add policy override (file.read under work)..."
$policy = '{"overrides":[{"tool":"file.read","allow":true,"constraints":{"roots":["work"],"max_bytes":1048576}}]}'
Invoke-RestMethod -Uri "$Base/v1/sessions/$SessionId/policy" -Method Put -Headers $Headers -Body $policy | Out-Null
Write-Host "   Policy set."

Write-Host "4. Retry same message (expect ALLOW + tool_exec ok)..."
$r = Invoke-RestMethod -Uri "$Base/v1/sessions/$SessionId/messages" -Method Post -Headers $Headers -Body $msgBody
$RunId2 = $r.run_id
Write-Host "   Run: $RunId2"

Start-Sleep -Seconds 2
$run2 = Invoke-RestMethod -Uri "$Base/v1/runs/$RunId2" -Method Get -Headers $Headers
Write-Host "   Status: $($run2.status)"
foreach ($s in $run2.steps) { Write-Host "   Step: $($s.type) $($s.status)" }

Write-Host "Done. First run denied; second run allowed after policy fix."
