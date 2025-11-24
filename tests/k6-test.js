import http from 'k6/http';
import { check, sleep } from 'k6';

const BASE_URL = 'http://host.docker.internal:8080';

export const options = {
  stages: [
    { duration: '60s', target: 3 },
  ],
  thresholds: {
    http_req_failed: ['rate<0.001'],   
    http_req_duration: ['p(95)<300'],  
    http_reqs: ['rate>5'],           
  },
};

export default function() {
  const randomTeamNum = Math.floor(Math.random() * 20) + 1;
  const randomUserNum = Math.floor(Math.random() * 200) + 1;
  
  const randomTeam = `team${randomTeamNum}`;
  const randomUser = `user${randomUserNum}`;
  
  const op = Math.random();
  
  if (op < 0.3) {
    // 30% - получение команды 
    const res = http.get(`${BASE_URL}/team/get?team_name=${randomTeam}`);
    check(res, { 'get team': (r) => r.status === 200 });
    
  } else if (op < 0.6) {
    // 30% - получение ревью 
    const res = http.get(`${BASE_URL}/users/getReview?user_id=${randomUser}`);
    check(res, { 'get reviews': (r) => r.status === 200 });
    
  } else if (op < 0.8) {
    // 20% - изменение активности пользователя
    const payload = {
      user_id: randomUser,
      is_active: Math.random() > 0.5
    };
    const res = http.post(`${BASE_URL}/users/setIsActive`, JSON.stringify(payload), {
      headers: { 'Content-Type': 'application/json' }
    });
    check(res, { 'set active': (r) => r.status === 200 });
    
  } else {
    // 20% - мерж PR
    const randomPRNum = Math.floor(Math.random() * 2) + 1;
    const prId = `pr_team${randomTeamNum}_${randomPRNum}`;
    
    const payload = { pull_request_id: prId };
    const res = http.post(`${BASE_URL}/pullRequest/merge`, JSON.stringify(payload), {
      headers: { 'Content-Type': 'application/json' }
    });
    check(res, { 'merge PR': (r) => r.status === 200 });
  }
  
  sleep(0.2);
}