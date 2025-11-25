import http from 'k6/http';
import { check } from 'k6';

const BASE_URL = 'http://host.docker.internal:8080';

export const options = {
  vus: 1,
  iterations: 1,
};

export default function() {
  const headers = { 'Content-Type': 'application/json' };
  const timestamp = Date.now();
  
  const teamPayload = {
    team_name: `e2e_team_${timestamp}`,
    members: [
      { user_id: 'e2e_user1', username: 'E2E_User1', is_active: true },
      { user_id: 'e2e_user2', username: 'E2E_User2', is_active: true },
    ],
  };
  
  let res = http.post(`${BASE_URL}/team/add`, JSON.stringify(teamPayload), { headers });
  check(res, { 'создание команды': (r) => r.status === 201 });
  
  res = http.get(`${BASE_URL}/team/get?team_name=${teamPayload.team_name}`);
  check(res, { 'получение команды': (r) => r.status === 200 });
  
  const prPayload = {
    pull_request_id: `e2e_pr_${timestamp}`,
    pull_request_name: 'E2E Test Feature',
    author_id: 'e2e_user1',
  };
  
  res = http.post(`${BASE_URL}/pullRequest/create`, JSON.stringify(prPayload), { headers });
  const prCreated = check(res, { 'создание pr': (r) => r.status === 201 });
  
  let assignedReviewers = [];
  if (prCreated && res.status === 201) {
    const prData = JSON.parse(res.body);
    assignedReviewers = prData.pr.assigned_reviewers;
  }
  
  if (assignedReviewers.length > 0) {
    res = http.get(`${BASE_URL}/users/getReview?user_id=${assignedReviewers[0]}`);
    check(res, { 'получение pr для ревьювера': (r) => r.status === 200 });
  }
  
  if (assignedReviewers.length > 0) {
    res = http.post(`${BASE_URL}/pullRequest/reassign`, JSON.stringify({
      pull_request_id: prPayload.pull_request_id,
      old_user_id: assignedReviewers[0]
    }), { headers });
    
    check(res, { 'ошибка no_candidate при переназначении': (r) => r.status === 409 });
    
    if (res.status === 409) {
      const errorData = JSON.parse(res.body);
      check(errorData, {
        'код ошибки no_candidate': (e) => e.error.code === 'NO_CANDIDATE'
      });
    }
  }
  
  res = http.post(`${BASE_URL}/users/setIsActive`, JSON.stringify({
    user_id: 'e2e_user2',
    is_active: false
  }), { headers });
  check(res, { 'обновление статуса пользователя': (r) => r.status === 200 });
  
  res = http.post(`${BASE_URL}/pullRequest/merge`, JSON.stringify({
    pull_request_id: prPayload.pull_request_id
  }), { headers });
  check(res, { 'merge pr': (r) => r.status === 200 });
  
  res = http.post(`${BASE_URL}/team/add`, JSON.stringify(teamPayload), { headers });
  check(res, { 'ошибка дублирования команды': (r) => r.status === 409 });
  
  res = http.get(`${BASE_URL}/team/get?team_name=nonexistent_team`);
  check(res, { 'ошибка несуществующей команды': (r) => r.status === 404 });
  
  res = http.post(`${BASE_URL}/pullRequest/create`, JSON.stringify({
    pull_request_id: `e2e_pr_invalid_${timestamp}`,
    pull_request_name: 'Invalid PR',
    author_id: 'non_existent_user',
  }), { headers });
  check(res, { 'ошибка несуществующего автора': (r) => r.status === 404 });
  
  res = http.post(`${BASE_URL}/pullRequest/reassign`, JSON.stringify({
    pull_request_id: prPayload.pull_request_id,
    old_user_id: 'e2e_user2'
  }), { headers });
  check(res, { 'ошибка переназначения на merged pr': (r) => r.status === 409 });
}