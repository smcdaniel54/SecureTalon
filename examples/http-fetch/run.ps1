# HTTP fetch example — allow GET to jsonplaceholder.typicode.com, then send intent.
# Prereq: Backend running: ADDR=:8090 ADMIN_TOKEN=demo go run ./cmd/securetalon

$Base = "http://localhost:8090"
$Token = "demo"
$Headers = @{
    "Authorization" = "Bearer $Token"
    "Content-Type"  = "application/json"
}

Write-Host "1. Create session..."
$r = Invoke-RestMethod -Uri "$Base/v1/sessions" -Method Post -Headers $Headers -Body '{"label":"HTTP fetch example","metadata":{}}'
$SessionId = $r.id
Write-Host "   Session: $SessionId"

Write-Host "2. Set session policy (http.fetch to jsonplaceholder.typicode.com, GET only)..."
$policy = @{
    overrides = @(
        @{
            tool       = "http.fetch"
            allow      = $true
            constraints = @{
                domains   = @("https://jsonplaceholder.typicode.com")
                methods   = @("GET")
                max_bytes = 200000
            }
        }
    )
} | ConvertTo-Json -Depth 5
Invoke-RestMethod -Uri "$Base/v1/sessions/$SessionId/policy" -Method Put -Headers $Headers -Body $policy | Out-Null
Write-Host "   Policy set."

Write-Host "3. Post message with http.fetch intent..."
$msgBody = '{"role":"user","content":"Fetch post 1 from API","intents":[{"tool":"http.fetch","params":{"url":"https://jsonplaceholder.typicode.com/posts/1","method":"GET"}}]}'
$r = Invoke-RestMethod -Uri "$Base/v1/sessions/$SessionId/messages" -Method Post -Headers $Headers -Body $msgBody
$RunId = $r.run_id
Write-Host "   Run: $RunId"

Write-Host "4. Wait and get run..."
Start-Sleep -Seconds 3
$run = Invoke-RestMethod -Uri "$Base/v1/runs/$RunId" -Method Get -Headers $Headers
Write-Host "   Status: $($run.status)"
if ($run.steps) {
    foreach ($s in $run.steps) {
        Write-Host "   Step: $($s.type) $($s.status)"
        if ($s.details) { Write-Host "     Details: $($s.details | ConvertTo-Json -Compress)" }
    }
}

Write-Host "Done. Open UI -> Sessions -> $SessionId -> Run $RunId to see response in step details."
