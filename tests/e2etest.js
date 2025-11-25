import http from 'k6/http';
import { check } from 'k6';

const BASE_URL = 'http://host.docker.internal:8080';

export const options = {
  vus: 1,
  iterations: 1,
};

export default function() {  
  const edgeCaseTeam = {
    team_name: 'edge_case_team',
    members: [
      { user_id: 'user_edge_1', username: 'User Edge 1', is_active: true },
      { user_id: 'user_edge_1', username: 'User Edge 1 Duplicate', is_active: true }, 
    ],
  };

  let res = http.post(
    `${BASE_URL}/team/add`,
    JSON.stringify(edgeCaseTeam),
    { headers: { 'Content-Type': 'application/json' } }
  );
  check(res, {
    'duplicate user_id in team rejected': r => r.status === 400,
  });

  res = http.post(
    `${BASE_URL}/pullRequest/create`,
    JSON.stringify({
      pull_request_id: 'pr_inactive_author',
      pull_request_name: 'PR with inactive author',
      author_id: 'user_edge_1',
    }),
    { headers: { 'Content-Type': 'application/json' } }
  );
  check(res, {
    'PR with inactive author rejected': r => r.status === 404,
  });

  res = http.post(
    `${BASE_URL}/pullRequest/create`,
    JSON.stringify({
      pull_request_id: 'pr_no_reviewers',
      pull_request_name: 'PR with no active reviewers',
      author_id: 'user_1', 
    }),
    { headers: { 'Content-Type': 'application/json' } }
  );
  check(res, {
    'PR created without reviewers': r => r.status === 201,
  });

  res = http.post(
    `${BASE_URL}/pullRequest/reassign`,
    JSON.stringify({
      pull_request_id: 'pr_1',
      old_user_id: 'user_2',
      new_user_id: 'user_edge_1', 
    }),
    { headers: { 'Content-Type': 'application/json' } }
  );
  check(res, {
    'reassignment to inactive user rejected': r => r.status === 409 || r.status === 404,
  });

  res = http.post(
    `${BASE_URL}/pullRequest/reassign`,
    JSON.stringify({
      pull_request_id: 'pr_1',
      old_user_id: 'user_2',
      new_user_id: 'user_1',
    }),
    { headers: { 'Content-Type': 'application/json' } }
  );
  check(res, {
    'reassignment to author rejected': r => r.status === 409 || r.status === 404,
  });

  res = http.post(
    `${BASE_URL}/pullRequest/merge`,
    JSON.stringify({ pull_request_id: 'non_existent_pr' }),
    { headers: { 'Content-Type': 'application/json' } }
  );
  check(res, {
    'merge non-existent PR rejected': r => r.status === 404,
  });

  res = http.get(`${BASE_URL}/users/getReview?user_id=non_existent_user`);
  check(res, {
    'get reviews for non-existent user rejected': r => r.status === 404,
  });
}
