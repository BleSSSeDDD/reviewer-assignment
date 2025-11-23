import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '10s', target: 20 },
    { duration: '30s', target: 30 },
    { duration: '20s', target: 50 },
    { duration: '10s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<300'], 
    http_req_failed: ['rate<0.05'],
  },
};

export default function() {
  const BASE_URL = 'http://host.docker.internal:8080';
  
  const endpoints = [
    '/team/get?team_name=test_backend',
    '/users/getReview?user_id=t1',
    '/users/getReview?user_id=t2', 
    '/team/get?team_name=test_frontend'
  ];
  
  const url = BASE_URL + endpoints[Math.floor(Math.random() * endpoints.length)];
  const res = http.get(url);
  
  check(res, {
    'status is 200': (r) => r.status === 200,
    'response time < 500ms': (r) => r.timings.duration < 500,
  });
  
  sleep(0.05);
}