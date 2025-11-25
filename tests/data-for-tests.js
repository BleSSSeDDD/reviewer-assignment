import http from 'k6/http';
import { check } from 'k6';

const BASE_URL = 'http://host.docker.internal:8080';

export const options = {
  vus: 1,
  iterations: 1,
};

export default function() {
  const teams = [];
  
  for (let teamNum = 1; teamNum <= 20; teamNum++) {
    const teamName = `team${teamNum}`;
    const members = [];
    
    for (let userNum = 1; userNum <= 10; userNum++) {
      const userId = `user${(teamNum - 1) * 10 + userNum}`;
      members.push({
        user_id: userId,
        username: `User${userId}`,
        is_active: Math.random() > 0.3
      });
    }
    
    teams.push({
      team_name: teamName,
      members: members
    });
  }

  const teamRequests = teams.map(team => ({
    method: 'POST',
    url: `${BASE_URL}/team/add`,
    body: JSON.stringify(team),
    params: {
      headers: { 'Content-Type': 'application/json' }
    }
  }));

  const teamResponses = http.batch(teamRequests);
  
  let successCount = 0;
  teamResponses.forEach((res, i) => {
    if (res.status === 201 || res.status === 409) successCount++;
  });
  
  console.log(` Created ${successCount}/20 teams with 200 users`);

  console.log(' Creating test PRs...');
  const prRequests = [];
  
  for (let teamNum = 1; teamNum <= 20; teamNum++) {
    const author1 = `user${(teamNum - 1) * 10 + 1}`; 
    const author2 = `user${(teamNum - 1) * 10 + 2}`;
    
    prRequests.push({
      method: 'POST',
      url: `${BASE_URL}/pullRequest/create`,
      body: JSON.stringify({
        pull_request_id: `pr_team${teamNum}_1`,
        pull_request_name: `Test PR for team ${teamNum} - 1`,
        author_id: author1
      }),
      params: {
        headers: { 'Content-Type': 'application/json' }
      }
    });
    
    prRequests.push({
      method: 'POST',
      url: `${BASE_URL}/pullRequest/create`,
      body: JSON.stringify({
        pull_request_id: `pr_team${teamNum}_2`, 
        pull_request_name: `Test PR for team ${teamNum} - 2`,
        author_id: author2
      }),
      params: {
        headers: { 'Content-Type': 'application/json' }
      }
    });
  }

  const prResponses = http.batch(prRequests);
  let prSuccessCount = 0;
  prResponses.forEach((res, i) => {
    if (res.status === 201 || res.status === 409) prSuccessCount++;
  });

  console.log(` Created ${prSuccessCount}/40 test PRs`);
  console.log(' Test data ready: 20 teams, 200 users, 40 PRs');
}