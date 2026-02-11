# Docker skill example — register image by digest, allow docker.run, send intent.
# Prereq: Backend with Docker: ADDR=:8090 ADMIN_TOKEN=demo go run ./cmd/securetalon
# Set your image digest (e.g. docker inspect hello-world for RepoDigests).
$ImageDigest = "hello-world@sha256:d6e2aa7d1e80dc6a2e96e64d30c6ee2593e1e7f42e0e2e8e8e8e8e8e8e8e8e8e8"
# If you have real digest, replace above. Example: "hello-world@sha256:1a728cd73f3f2cf7d2df4b1e8f8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8"

$Base = "http://localhost:8090"
$Token = "demo"
$Headers = @{
    "Authorization" = "Bearer $Token"
    "Content-Type"  = "application/json"
}

Write-Host "1. (Optional) Register skill..."
try {
    $skillBody = @{ name = "hello-world"; version = "1.0"; image = $ImageDigest; manifest = @{} } | ConvertTo-Json
    Invoke-RestMethod -Uri "$Base/v1/skills" -Method Post -Headers $Headers -Body $skillBody | Out-Null
    Write-Host "   Skill registered."
} catch {
    Write-Host "   (Skip or already registered)"
}

Write-Host "2. Create session..."
$r = Invoke-RestMethod -Uri "$Base/v1/sessions" -Method Post -Headers $Headers -Body '{"label":"Docker example","metadata":{}}'
$SessionId = $r.id
Write-Host "   Session: $SessionId"

Write-Host "3. Set session policy (docker.run for image digest only, no network/mounts)..."
$policy = @{
    overrides = @(
        @{
            tool       = "docker.run"
            allow      = $true
            constraints = @{
                images          = @($ImageDigest)
                network_allowed = $false
                mounts_allowed  = $false
            }
        }
    )
} | ConvertTo-Json -Depth 5
Invoke-RestMethod -Uri "$Base/v1/sessions/$SessionId/policy" -Method Put -Headers $Headers -Body $policy | Out-Null
Write-Host "   Policy set."

Write-Host "4. Post message with docker.run intent..."
$msgBody = @{
    role    = "user"
    content = "Run container"
    intents = @(@{ tool = "docker.run"; params = @{ image = $ImageDigest; cmd = @() } })
} | ConvertTo-Json -Depth 5
$r = Invoke-RestMethod -Uri "$Base/v1/sessions/$SessionId/messages" -Method Post -Headers $Headers -Body $msgBody
$RunId = $r.run_id
Write-Host "   Run: $RunId"

Write-Host "5. Wait and get run..."
Start-Sleep -Seconds 5
$run = Invoke-RestMethod -Uri "$Base/v1/runs/$RunId" -Method Get -Headers $Headers
Write-Host "   Status: $($run.status)"
if ($run.steps) {
    foreach ($s in $run.steps) { Write-Host "   Step: $($s.type) $($s.status)" }
}

Write-Host "Done. Open UI -> Sessions -> $SessionId -> Run $RunId"
