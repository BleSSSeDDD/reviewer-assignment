import http from 'k6/http';
import { check, sleep } from 'k6';

const BASE_URL = 'http://host.docker.internal:8080';

export const options = {
  vus: 1,
  iterations: 1,
};

let testTeams = [];
let testUsers = [];
let testPRs = [];

export default function() {
  console.log('=== Starting API Tests ===\n');

  testTeamCreation();
  sleep(1);

  testGetTeams();
  sleep(1);

  testUserActivation();
  sleep(1);

  testPRCreation();
  sleep(1);

  testGetUserReviews();
  sleep(1);

  testReassignReviewers();
  sleep(1);

  testPRMerge();
  sleep(1);

  testErrorCases();
  sleep(1);

  console.log('\n=== All tests completed ===');
}

function testTeamCreation() {
  console.log('1. Testing Team Creation...');
  
  const teams = [];
  for (let teamNum = 1; teamNum <= 5; teamNum++) {
    const teamName = `test_team_${teamNum}`;
    const members = [];
    
    for (let userNum = 1; userNum <= 5; userNum++) {
      const userId = `test_user_${(teamNum - 1) * 5 + userNum}`;
      members.push({
        user_id: userId,
        username: `TestUser${userId}`,
        is_active: true
      });
      testUsers.push({ user_id: userId, team_name: teamName });
    }
    
    const team = {
      team_name: teamName,
      members: members
    };
    teams.push(team);
    testTeams.push(teamName);
  }

  const requests = teams.map(team => ({
    method: 'POST',
    url: `${BASE_URL}/team/add`,
    body: JSON.stringify(team),
    params: { headers: { 'Content-Type': 'application/json' } }
  }));

  const responses = http.batch(requests);
  
  let successCount = 0;
  responses.forEach((res, i) => {
    const checkResult = check(res, {
      [`team ${teams[i].team_name} created`]: (r) => r.status === 201 || r.status === 409
    });
    if (checkResult) successCount++;
  });

  console.log(`   Created ${successCount}/${teams.length} teams\n`);
}

function testGetTeams() {
  console.log('2. Testing Team Retrieval...');
  
  const requests = testTeams.map(teamName => ({
    method: 'GET',
    url: `${BASE_URL}/team/get?team_name=${teamName}`,
    params: { headers: { 'Content-Type': 'application/json' } }
  }));

  const responses = http.batch(requests);
  
  let successCount = 0;
  responses.forEach((res, i) => {
    const checkResult = check(res, {
      [`team ${testTeams[i]} retrieved`]: (r) => r.status === 200
    });
    if (checkResult) successCount++;
  });

  console.log(`   Retrieved ${successCount}/${testTeams.length} teams\n`);
}

function testUserActivation() {
  console.log('3. Testing User Activation...');
  
  const deactivationRequests = [
    { user_id: 'test_user_2', is_active: false },
    { user_id: 'test_user_7', is_active: false }
  ].map(data => ({
    method: 'POST',
    url: `${BASE_URL}/users/setIsActive`,
    body: JSON.stringify(data),
    params: { headers: { 'Content-Type': 'application/json' } }
  }));

  const deactivationResponses = http.batch(deactivationRequests);
  
  let deactivationSuccess = 0;
  deactivationResponses.forEach((res, i) => {
    const checkResult = check(res, {
      [`user ${deactivationRequests[i].body.user_id} deactivated`]: (r) => r.status === 200
    });
    if (checkResult) deactivationSuccess++;
  });

  console.log(`   Deactivated ${deactivationSuccess}/2 users\n`);
}

function testPRCreation() {
  console.log('4. Testing PR Creation with Auto-Assignment...');
  
  const prRequests = [
    {
      pull_request_id: 'test_pr_1',
      pull_request_name: 'Test PR 1 - Feature Implementation',
      author_id: 'test_user_1'
    },
    {
      pull_request_id: 'test_pr_2', 
      pull_request_name: 'Test PR 2 - Bug Fix',
      author_id: 'test_user_3'
    },
    {
      pull_request_id: 'test_pr_3',
      pull_request_name: 'Test PR 3 - Refactoring',
      author_id: 'test_user_6'
    }
  ];

  const requests = prRequests.map(pr => ({
    method: 'POST',
    url: `${BASE_URL}/pullRequest/create`,
    body: JSON.stringify(pr),
    params: { headers: { 'Content-Type': 'application/json' } }
  }));

  const responses = http.batch(requests);
  
  let successCount = 0;
  responses.forEach((res, i) => {
    const prData = prRequests[i];
    testPRs.push(prData.pull_request_id);
    
    const checkResult = check(res, {
      [`PR ${prData.pull_request_id} created`]: (r) => r.status === 201 || r.status === 409
    });
    if (checkResult) successCount++;
  });

  console.log(`   Created ${successCount}/${prRequests.length} PRs\n`);
}

function testGetUserReviews() {
  console.log('5. Testing Get User Reviews...');
  
  const testReviewers = ['test_user_4', 'test_user_5', 'test_user_8'];
  
  const requests = testReviewers.map(userId => ({
    method: 'GET',
    url: `${BASE_URL}/users/getReview?user_id=${userId}`,
    params: { headers: { 'Content-Type': 'application/json' } }
  }));

  const responses = http.batch(requests);
  
  let successCount = 0;
  responses.forEach((res, i) => {
    const checkResult = check(res, {
      [`reviews for user ${testReviewers[i]} retrieved`]: (r) => r.status === 200
    });
    if (checkResult) successCount++;
  });

  console.log(`   Retrieved reviews for ${successCount}/${testReviewers.length} users\n`);
}

function testReassignReviewers() {
  console.log('6. Testing Reviewer Reassignment...');
  
  if (testPRs.length === 0) {
    console.log('   No PRs available for reassignment test\n');
    return;
  }

  const reassignmentRequest = {
    method: 'POST',
    url: `${BASE_URL}/pullRequest/reassign`,
    body: JSON.stringify({
      pull_request_id: 'test_pr_1',
      old_user_id: 'test_user_4'
    }),
    params: { headers: { 'Content-Type': 'application/json' } }
  };

  const response = http.post(
    reassignmentRequest.url, 
    reassignmentRequest.body, 
    reassignmentRequest.params
  );

  const reassignmentSuccess = check(response, {
    'reassignment handled correctly': (r) => r.status === 200 || r.status === 409 || r.status === 404
  });

  console.log(`   Reassignment test completed: ${reassignmentSuccess ? 'PASS' : 'FAIL'}\n`);
}

function testPRMerge() {
  console.log('7. Testing PR Merge...');
  
  if (testPRs.length === 0) {
    console.log('   No PRs available for merge test\n');
    return;
  }

  const mergeRequest = {
    method: 'POST',
    url: `${BASE_URL}/pullRequest/merge`,
    body: JSON.stringify({
      pull_request_id: 'test_pr_2'
    }),
    params: { headers: { 'Content-Type': 'application/json' } }
  };

  const response = http.post(mergeRequest.url, mergeRequest.body, mergeRequest.params);

  const mergeSuccess = check(response, {
    'PR merged successfully': (r) => r.status === 200
  });

  const retryResponse = http.post(mergeRequest.url, mergeRequest.body, mergeRequest.params);
  const idempotencySuccess = check(retryResponse, {
    'repeated merge is idempotent': (r) => r.status === 200
  });

  console.log(`   Merge test: ${mergeSuccess ? 'PASS' : 'FAIL'}`);
  console.log(`   Idempotency test: ${idempotencySuccess ? 'PASS' : 'FAIL'}\n`);
}

function testErrorCases() {
  console.log('8. Testing Error Cases...');
  
  const errorTests = [
    {
      name: 'duplicate team creation',
      request: {
        method: 'POST',
        url: `${BASE_URL}/team/add`,
        body: JSON.stringify({
          team_name: 'test_team_1',
          members: [{ user_id: 'new_user', username: 'New User', is_active: true }]
        }),
        params: { headers: { 'Content-Type': 'application/json' } }
      },
      expectedStatus: 400
    },
    {
      name: 'get non-existent team',
      request: {
        method: 'GET',
        url: `${BASE_URL}/team/get?team_name=non_existent_team`,
        params: { headers: { 'Content-Type': 'application/json' } }
      },
      expectedStatus: 404
    },
    {
      name: 'create PR with non-existent author',
      request: {
        method: 'POST',
        url: `${BASE_URL}/pullRequest/create`,
        body: JSON.stringify({
          pull_request_id: 'error_test_pr',
          pull_request_name: 'Error Test PR',
          author_id: 'non_existent_user'
        }),
        params: { headers: { 'Content-Type': 'application/json' } }
      },
      expectedStatus: 404
    }
  ];

  const requests = errorTests.map(test => test.request);
  const responses = http.batch(requests);
  
  let successCount = 0;
  responses.forEach((res, i) => {
    const test = errorTests[i];
    const checkResult = check(res, {
      [`error case "${test.name}" handled correctly`]: (r) => r.status === test.expectedStatus
    });
    if (checkResult) successCount++;
  });

  console.log(`   Handled ${successCount}/${errorTests.length} error cases correctly\n`);
}