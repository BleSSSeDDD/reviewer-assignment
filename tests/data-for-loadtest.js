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
    const teamName = `team_${teamNum}`;
    const members = [];
    
    for (let userNum = 1; userNum <= 10; userNum++) {
      const userId = `user_${(teamNum - 1) * 10 + userNum}`;
      members.push({
        user_id: userId,
        username: `User ${userId}`,
        is_active: true
      });
    }
    
    teams.push({
      team_name: teamName,
      members: members
    });
  }

  teams.forEach(team => {
    const res = http.post(`${BASE_URL}/team/add`, JSON.stringify(team), {
      headers: { 'Content-Type': 'application/json' }
    });
    
    check(res, {
      [`team ${team.team_name} created`]: (r) => r.status === 201 || r.status === 400
    });
  });
  
  console.log('Created 20 teams with 200 users');
}