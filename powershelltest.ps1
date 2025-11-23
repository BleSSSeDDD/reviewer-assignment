# test everything

docker compose down -v
docker compose up -d

$baseUrl = "http://localhost:8080"

Write-Host "start testing" -ForegroundColor Green

# create teams
Write-Host "create backend team" -ForegroundColor Cyan
Invoke-RestMethod -Uri "$baseUrl/team/add" -Method POST -ContentType "application/json" -Body '{"team_name": "test_backend", "members": [{"user_id": "t1", "username": "TestAlice", "is_active": true}, {"user_id": "t2", "username": "TestBob", "is_active": true}, {"user_id": "t3", "username": "TestCharlie", "is_active": true}]}'

Write-Host "create frontend team" -ForegroundColor Cyan
Invoke-RestMethod -Uri "$baseUrl/team/add" -Method POST -ContentType "application/json" -Body '{"team_name": "test_frontend", "members": [{"user_id": "t4", "username": "TestDavid", "is_active": true}, {"user_id": "t5", "username": "TestEve", "is_active": true}]}'

Write-Host "create mobile team" -ForegroundColor Cyan
Invoke-RestMethod -Uri "$baseUrl/team/add" -Method POST -ContentType "application/json" -Body '{"team_name": "test_mobile", "members": [{"user_id": "t6", "username": "TestFrank", "is_active": true}, {"user_id": "t7", "username": "TestGrace", "is_active": true}]}'

# check teams
Write-Host "check backend team members" -ForegroundColor Cyan
(Invoke-RestMethod -Uri "$baseUrl/team/get?team_name=test_backend" -Method GET).members

Write-Host "check frontend team members" -ForegroundColor Cyan
(Invoke-RestMethod -Uri "$baseUrl/team/get?team_name=test_frontend" -Method GET).members

Write-Host "check mobile team members" -ForegroundColor Cyan
(Invoke-RestMethod -Uri "$baseUrl/team/get?team_name=test_mobile" -Method GET).members

# test users
Write-Host "set user t1 active true" -ForegroundColor Cyan
Invoke-RestMethod -Uri "$baseUrl/users/setIsActive" -Method POST -ContentType "application/json" -Body '{"user_id": "t1", "is_active": true}'

Write-Host "set user t2 active true" -ForegroundColor Cyan
Invoke-RestMethod -Uri "$baseUrl/users/setIsActive" -Method POST -ContentType "application/json" -Body '{"user_id": "t2", "is_active": true}'

Write-Host "check PR for user t1" -ForegroundColor Cyan
(Invoke-RestMethod -Uri "$baseUrl/users/getReview?user_id=t1" -Method GET).pull_requests

Write-Host "check PR for user t2" -ForegroundColor Cyan
(Invoke-RestMethod -Uri "$baseUrl/users/getReview?user_id=t2" -Method GET).pull_requests

# error cases
Write-Host "try get nonexistent team" -ForegroundColor Cyan
try {
    Invoke-RestMethod -Uri "$baseUrl/team/get?team_name=nonexistent" -Method GET
} catch {
    Write-Host "error (expected): $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "try change nonexistent user" -ForegroundColor Cyan
try {
    Invoke-RestMethod -Uri "$baseUrl/users/setIsActive" -Method POST -ContentType "application/json" -Body '{"user_id": "ghost", "is_active": true}'
} catch {
    Write-Host "error (expected): $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "check PR for nonexistent user" -ForegroundColor Cyan
(Invoke-RestMethod -Uri "$baseUrl/users/getReview?user_id=ghost" -Method GET).pull_requests

# final check
Write-Host "all teams and members:" -ForegroundColor Cyan

Write-Host "backend team:" -ForegroundColor Green
(Invoke-RestMethod -Uri "$baseUrl/team/get?team_name=test_backend" -Method GET).members

Write-Host "frontend team:" -ForegroundColor Green
(Invoke-RestMethod -Uri "$baseUrl/team/get?team_name=test_frontend" -Method GET).members

Write-Host "mobile team:" -ForegroundColor Green
(Invoke-RestMethod -Uri "$baseUrl/team/get?team_name=test_mobile" -Method GET).members

# NEW TESTS FOR PULL REQUESTS

Write-Host "=== PULL REQUEST TESTS ===" -ForegroundColor Yellow

# Test creating PR with automatic reviewer assignment
Write-Host "create PR for user t1 (backend team)" -ForegroundColor Cyan
$pr1 = Invoke-RestMethod -Uri "$baseUrl/pullRequest/create" -Method POST -ContentType "application/json" -Body '{"pull_request_id": "pr-001", "pull_request_name": "Add authentication", "author_id": "t1"}'
Write-Host "Created PR: $($pr1.pr | ConvertTo-Json -Depth 3)" -ForegroundColor Green

Write-Host "create PR for user t4 (frontend team)" -ForegroundColor Cyan
$pr2 = Invoke-RestMethod -Uri "$baseUrl/pullRequest/create" -Method POST -ContentType "application/json" -Body '{"pull_request_id": "pr-002", "pull_request_name": "Update UI components", "author_id": "t4"}'
Write-Host "Created PR: $($pr2.pr | ConvertTo-Json -Depth 3)" -ForegroundColor Green

# Test duplicate PR creation
Write-Host "try create duplicate PR" -ForegroundColor Cyan
try {
    Invoke-RestMethod -Uri "$baseUrl/pullRequest/create" -Method POST -ContentType "application/json" -Body '{"pull_request_id": "pr-001", "pull_request_name": "Duplicate PR", "author_id": "t1"}'
} catch {
    Write-Host "error (expected): $($_.Exception.Message)" -ForegroundColor Red
}

# Test PR for non-existent user
Write-Host "try create PR for non-existent user" -ForegroundColor Cyan
try {
    Invoke-RestMethod -Uri "$baseUrl/pullRequest/create" -Method POST -ContentType "application/json" -Body '{"pull_request_id": "pr-003", "pull_request_name": "Ghost PR", "author_id": "ghost"}'
} catch {
    Write-Host "error (expected): $($_.Exception.Message)" -ForegroundColor Red
}

# Test merging PR
Write-Host "merge PR pr-001" -ForegroundColor Cyan
$mergedPr = Invoke-RestMethod -Uri "$baseUrl/pullRequest/merge" -Method POST -ContentType "application/json" -Body '{"pull_request_id": "pr-001"}'
Write-Host "Merged PR: $($mergedPr.pr | ConvertTo-Json -Depth 3)" -ForegroundColor Green

# Test merging non-existent PR
Write-Host "try merge non-existent PR" -ForegroundColor Cyan
try {
    Invoke-RestMethod -Uri "$baseUrl/pullRequest/merge" -Method POST -ContentType "application/json" -Body '{"pull_request_id": "pr-999"}'
} catch {
    Write-Host "error (expected): $($_.Exception.Message)" -ForegroundColor Red
}

# Test reassigning reviewer - ФИКС ЕБУЧЕГО ПАВЕРШЕЛЛА
Write-Host "reassign reviewer for PR pr-002" -ForegroundColor Cyan
if ($pr2.pr.assigned_reviewers -and $pr2.pr.assigned_reviewers.Count -gt 0) {
    $oldUserId = $pr2.pr.assigned_reviewers[0]
    $reassignBody = @{
        pull_request_id = "pr-002"
        old_user_id = $oldUserId
    } | ConvertTo-Json
    
    $reassigned = Invoke-RestMethod -Uri "$baseUrl/pullRequest/reassign" -Method POST -ContentType "application/json" -Body $reassignBody
    Write-Host "Reassigned PR: $($reassigned.pr | ConvertTo-Json -Depth 3)" -ForegroundColor Green
    Write-Host "Replaced by: $($reassigned.replaced_by)" -ForegroundColor Green
} else {
    Write-Host "No reviewers to reassign" -ForegroundColor Yellow
}

# Test reassignment errors
Write-Host "try reassign non-assigned reviewer" -ForegroundColor Cyan
try {
    $reassignBody = @{
        pull_request_id = "pr-002"
        old_user_id = "t1"
    } | ConvertTo-Json
    
    Invoke-RestMethod -Uri "$baseUrl/pullRequest/reassign" -Method POST -ContentType "application/json" -Body $reassignBody
} catch {
    Write-Host "error (expected): $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "try reassign on merged PR" -ForegroundColor Cyan
if ($pr1.pr.assigned_reviewers -and $pr1.pr.assigned_reviewers.Count -gt 0) {
    $oldUserId = $pr1.pr.assigned_reviewers[0]
    $reassignBody = @{
        pull_request_id = "pr-001"
        old_user_id = $oldUserId
    } | ConvertTo-Json
    
    try {
        Invoke-RestMethod -Uri "$baseUrl/pullRequest/reassign" -Method POST -ContentType "application/json" -Body $reassignBody
    } catch {
        Write-Host "error (expected): $($_.Exception.Message)" -ForegroundColor Red
    }
} else {
    Write-Host "No reviewers in merged PR to reassign" -ForegroundColor Yellow
}

# Test user review assignments
Write-Host "check PR assignments for all users" -ForegroundColor Cyan
foreach ($user in @("t1", "t2", "t3", "t4", "t5", "t6", "t7")) {
    $userReviews = Invoke-RestMethod -Uri "$baseUrl/users/getReview?user_id=$user" -Method GET
    Write-Host "User $user has $($userReviews.pull_requests.Count) PRs to review" -ForegroundColor Magenta
    if ($userReviews.pull_requests.Count -gt 0) {
        Write-Host "PRs: $($userReviews.pull_requests | ConvertTo-Json -Depth 2)" -ForegroundColor Magenta
    }
}

# Test deactivating user and creating new PR - ФИКС JSON
Write-Host "deactivate user t2 and create new PR" -ForegroundColor Cyan
$deactivateBody = '{"user_id": "t2", "is_active": false}'
try {
    Invoke-RestMethod -Uri "$baseUrl/users/setIsActive" -Method POST -ContentType "application/json" -Body $deactivateBody
} catch {
    Write-Host "Error deactivating user: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "create another PR for backend team" -ForegroundColor Cyan
$pr3 = Invoke-RestMethod -Uri "$baseUrl/pullRequest/create" -Method POST -ContentType "application/json" -Body '{"pull_request_id": "pr-003", "pull_request_name": "Fix database issue", "author_id": "t3"}'
Write-Host "Created PR: $($pr3.pr | ConvertTo-Json -Depth 3)" -ForegroundColor Green

# Test edge case - team with only one active member
Write-Host "create team with only one active member" -ForegroundColor Cyan
Invoke-RestMethod -Uri "$baseUrl/team/add" -Method POST -ContentType "application/json" -Body '{"team_name": "test_solo", "members": [{"user_id": "s1", "username": "SoloUser", "is_active": true}]}'

Write-Host "create PR for solo team" -ForegroundColor Cyan
$pr4 = Invoke-RestMethod -Uri "$baseUrl/pullRequest/create" -Method POST -ContentType "application/json" -Body '{"pull_request_id": "pr-004", "pull_request_name": "Solo feature", "author_id": "s1"}'
Write-Host "Created PR: $($pr4.pr | ConvertTo-Json -Depth 3)" -ForegroundColor Green

# Test reassignment when no candidates available
Write-Host "try reassign in solo team (should fail)" -ForegroundColor Cyan
if ($pr4.pr.assigned_reviewers -and $pr4.pr.assigned_reviewers.Count -gt 0) {
    $oldUserId = $pr4.pr.assigned_reviewers[0]
    $reassignBody = @{
        pull_request_id = "pr-004"
        old_user_id = $oldUserId
    } | ConvertTo-Json
    
    try {
        Invoke-RestMethod -Uri "$baseUrl/pullRequest/reassign" -Method POST -ContentType "application/json" -Body $reassignBody
    } catch {
        Write-Host "error (expected): $($_.Exception.Message)" -ForegroundColor Red
    }
} else {
    Write-Host "No reviewers in solo PR to reassign (expected)" -ForegroundColor Yellow
}

# Final comprehensive check
Write-Host "=== FINAL STATUS CHECK ===" -ForegroundColor Yellow

Write-Host "All teams:" -ForegroundColor Green
$teams = @("test_backend", "test_frontend", "test_mobile", "test_solo")
foreach ($team in $teams) {
    try {
        $teamData = Invoke-RestMethod -Uri "$baseUrl/team/get?team_name=$team" -Method GET
        Write-Host "$team : $($teamData.members.Count) members" -ForegroundColor Cyan
    } catch {
        Write-Host "$team : not found" -ForegroundColor Red
    }
}

Write-Host "All PRs status:" -ForegroundColor Green
$prs = @("pr-001", "pr-002", "pr-003", "pr-004")
foreach ($prId in $prs) {
    try {
        # Get PR status through user reviews (since there's no direct PR get endpoint)
        $allUsers = @("t1", "t2", "t3", "t4", "t5", "t6", "t7", "s1")
        foreach ($user in $allUsers) {
            $reviews = Invoke-RestMethod -Uri "$baseUrl/users/getReview?user_id=$user" -Method GET
            $userPr = $reviews.pull_requests | Where-Object { $_.pull_request_id -eq $prId }
            if ($userPr) {
                Write-Host "PR $prId : $($userPr.status) (assigned to $user)" -ForegroundColor Cyan
                break
            }
        }
    } catch {
        Write-Host "PR $prId : error checking" -ForegroundColor Red
    }
}

Write-Host "testing done" -ForegroundColor Green