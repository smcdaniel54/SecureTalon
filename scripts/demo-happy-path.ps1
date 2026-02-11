# SecureTalon happy path demo (API only).
# Prereqs: Backend running: ADDR=:8090 ADMIN_TOKEN=demo go run ./cmd/securetalon
# Create work dir and file: New-Item -ItemType Directory -Force work | Out-Null; 'hello' | Set-Content work/input.txt

$Base = "http://localhost:8090"
$Token = "demo"
$Headers = @{
    "Authorization" = "Bearer $Token"
    "Content-Type"  = "application/json"
}

Write-Host "1. Create session..."
$r = Invoke-RestMethod -Uri "$Base/v1/sessions" -Method Post -Headers $Headers -Body '{"label":"Demo","metadata":{}}'
$SessionId = $r.id
Write-Host "   Session: $SessionId"

Write-Host "2. Set session policy (file.read under work)..."
$body = '{"overrides":[{"tool":"file.read","allow":true,"constraints":{"roots":["work"],"max_bytes":1048576}}]}'
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
if ($run.steps) {
    foreach ($s in $run.steps) { Write-Host "   Step: $($s.type) $($s.status)" }
}

Write-Host "5. Query audit (session filter)..."
$events = (Invoke-RestMethod -Uri "$Base/v1/audit?session_id=$SessionId&limit=50" -Method Get -Headers $Headers).events
Write-Host "   Events: $($events.Count)"

Write-Host "6. Validate chain..."
$v = Invoke-RestMethod -Uri "$Base/v1/audit/validate?session_id=$SessionId&limit=500" -Method Get -Headers $Headers
Write-Host "   Chain valid: $($v.valid) (event_count: $($v.event_count))"

Write-Host "7. Load safe replay..."
$replay = Invoke-RestMethod -Uri "$Base/v1/runs/$RunId/replay" -Method Post -Headers $Headers
Write-Host "   Replay run_id: $($replay.run_id) mode: $($replay.mode) events: $($replay.events.Count)"

Write-Host "Done. Open UI: cd ui && npm run dev -> Sessions -> $SessionId -> Run $RunId | Audit | Replay"
