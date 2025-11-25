import http from 'k6/http';
import { check } from 'k6';

const BASE_URL = 'http://host.docker.internal:8080';

export const options = {
  stages: [
    { duration: '3s', target: 1000 }, 
    { duration: '30s', target: 1000 }, 
  ],
  thresholds: {
    http_req_failed: ['rate<0.9'],   
    http_req_duration: ['p(95)<10000'], 
  },
};

export default function() {
  const randomTeamNum = Math.floor(Math.random() * 20) + 1;
  const randomUserNum = Math.floor(Math.random() * 200) + 1;
  
  const randomTeam = `team${randomTeamNum}`;
  const randomUser = `user${randomUserNum}`;

  http.get(`${BASE_URL}/team/get?team_name=${randomTeam}`);
  
  http.get(`${BASE_URL}/users/getReview?user_id=${randomUser}`);
  
  const prId = `pr_${Date.now()}_${Math.random().toString(36).substr(2, 9)}_${__VU}`;
  const payload = {
    pull_request_id: prId,
    pull_request_name: `KILL_SWITCH_${prId}`,
    author_id: randomUser
  };
  http.post(`${BASE_URL}/pullRequest/create`, JSON.stringify(payload), {
    headers: { 'Content-Type': 'application/json' }
  });
  
  const activePayload = {
    user_id: randomUser,
    is_active: Math.random() > 0.5
  };
  http.post(`${BASE_URL}/users/setIsActive`, JSON.stringify(activePayload), {
    headers: { 'Content-Type': 'application/json' }
  });
  
  http.post(`${BASE_URL}/pullRequest/merge`, JSON.stringify({
    pull_request_id: `pr_team${randomTeamNum}_1`
  }), {
    headers: { 'Content-Type': 'application/json' }
  });

  http.post(`${BASE_URL}/pullRequest/reassign`, JSON.stringify({
    pull_request_id: `pr_team${randomTeamNum}_1`,
    old_user_id: randomUser
  }), {
    headers: { 'Content-Type': 'application/json' }
  });

  http.get(`${BASE_URL}/team/get?team_name=${randomTeam}`);
  http.get(`${BASE_URL}/users/getReview?user_id=${randomUser}`);
}