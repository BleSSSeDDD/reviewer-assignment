import http from 'k6/http';
import { check } from 'k6';

const BASE_URL = 'http://host.docker.internal:8080';

export const options = {
  stages: [
    { duration: '15s', target: 5 },   
    { duration: '30s', target: 15 }, 
    { duration: '15s', target: 5 },   
  ],
  thresholds: {
    http_req_failed: ['rate<0.001'],  
    http_req_duration: ['p(95)<300'], 
  },
};

let prCounter = 1;

export default function() {
  const randomTeam = Math.floor(Math.random() * 20) + 1;
  const randomUser = Math.floor(Math.random() * 200) + 1;
  
  const teamName = `team_${randomTeam}`;
  const userId = `user_${randomUser}`;
  
  const op = Math.random();
  
  if (op < 0.3) {
    const res = http.get(`${BASE_URL}/team/get?team_name=${teamName}`);
    check(res, { 'get team': (r) => r.status === 200 || r.status === 404 });
    
  } else if (op < 0.5) {
    const res = http.get(`${BASE_URL}/users/getReview?user_id=${userId}`);
    check(res, { 'get reviews': (r) => r.status === 200 });
    
  } else if (op < 0.7) {
    const prId = `pr_${Date.now()}_${prCounter++}`;
    const payload = {
      pull_request_id: prId,
      pull_request_name: `PR ${prId}`,
      author_id: userId
    };
    const res = http.post(`${BASE_URL}/pullRequest/create`, JSON.stringify(payload), {
      headers: { 'Content-Type': 'application/json' }
    });
    check(res, { 'create PR': (r) => r.status === 201 || r.status === 404 || r.status === 409 });
    
  } else if (op < 0.85) {
    const payload = {
      user_id: userId,
      is_active: Math.random() > 0.5
    };
    const res = http.post(`${BASE_URL}/users/setIsActive`, JSON.stringify(payload), {
      headers: { 'Content-Type': 'application/json' }
    });
    check(res, { 'set active': (r) => r.status === 200 || r.status === 404 });
    
  } else if (op < 0.95) {
    const prId = prCounter > 1 ? `pr_${Date.now()}_${Math.floor(Math.random() * prCounter)}` : `pr_${Math.floor(Math.random() * 1000)}`;
    const payload = { pull_request_id: prId };
    const res = http.post(`${BASE_URL}/pullRequest/merge`, JSON.stringify(payload), {
      headers: { 'Content-Type': 'application/json' }
    });
    check(res, { 'merge PR': (r) => r.status === 200 || r.status === 404 });
    
  } else {
    const prId = prCounter > 1 ? `pr_${Date.now()}_${Math.floor(Math.random() * prCounter)}` : `pr_${Math.floor(Math.random() * 1000)}`;
    const payload = {
      pull_request_id: prId,
      old_user_id: userId
    };
    const res = http.post(`${BASE_URL}/pullRequest/reassign`, JSON.stringify(payload), {
      headers: { 'Content-Type': 'application/json' }
    });
    check(res, { 'reassign': (r) => r.status === 200 || r.status === 404 || r.status === 409 });
  }
}