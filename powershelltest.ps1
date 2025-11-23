# test everything

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

Write-Host "testing done" -ForegroundColor Green

(Invoke-RestMethod -Uri "http://localhost:8080/team/get?team_name=test_mobile" -Method GET).members

(Invoke-RestMethod -Uri "http://localhost:8080/team/get?team_name=test_backend" -Method GET).members
(Invoke-RestMethod -Uri "http://localhost:8080/team/get?team_name=test_frontend" -Method GET).members  
(Invoke-RestMethod -Uri "http://localhost:8080/team/get?team_name=test_mobile" -Method GET).members