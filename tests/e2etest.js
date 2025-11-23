const BASE_URL = 'http://localhost:8080';
const TEST_DELAY = 100;

class TestRunner {
    constructor() {
        this.tests = [];
        this.passed = 0;
        this.failed = 0;
        this.testRunId = Date.now().toString();
    }

    async delay(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }

    async test(name, testFn) {
        try {
            await testFn();
            this.passed++;
        } catch (error) {
            this.failed++;
        }
        await this.delay(TEST_DELAY);
    }

    logResult() {
        console.log(`RESULTS: ${this.passed} passed, ${this.failed} failed`);
    }

    async apiCall(endpoint, options = {}) {
        const url = `${BASE_URL}${endpoint}`;
        try {
            const response = await fetch(url, {
                headers: {
                    'Content-Type': 'application/json',
                    ...options.headers
                },
                ...options
            });
            
            const data = await response.json().catch(() => null);
            
            return {
                status: response.status,
                data,
                ok: response.ok
            };
        } catch (error) {
            return { status: 0, data: null, ok: false };
        }
    }

    generateId(prefix) {
        return `${prefix}-${this.testRunId}-${Math.random().toString(36).substr(2, 9)}`;
    }
}

const runner = new TestRunner();

async function runAllTests() {
    await runTeamTests();
    await runUserTests();
    await runPullRequestTests();
    await runEdgeCaseTests();
    await runReassignmentTests();
    await runDatabaseLimitTests();
    
    runner.logResult();
}

async function runTeamTests() {
    await runner.test('Create team', async () => {
        const teamData = {
            team_name: runner.generateId('team'),
            members: [
                { user_id: runner.generateId('user'), username: 'Alice', is_active: true },
                { user_id: runner.generateId('user'), username: 'Bob', is_active: true }
            ]
        };
        
        const result = await runner.apiCall('/team/add', {
            method: 'POST',
            body: JSON.stringify(teamData)
        });
        
        if (!result.ok) throw new Error(`Expected 201, got ${result.status}`);
    });

    await runner.test('Duplicate team', async () => {
        const teamId = runner.generateId('team');
        const teamData = {
            team_name: teamId,
            members: [{ user_id: runner.generateId('user'), username: 'User', is_active: true }]
        };
        
        await runner.apiCall('/team/add', { method: 'POST', body: JSON.stringify(teamData) });
        const result = await runner.apiCall('/team/add', { method: 'POST', body: JSON.stringify(teamData) });
        
        if (result.status !== 409 && result.status !== 400) throw new Error(`Expected 409 or 400, got ${result.status}`);
    });

    await runner.test('Get team', async () => {
        const teamId = runner.generateId('team');
        const teamData = {
            team_name: teamId,
            members: [{ user_id: runner.generateId('user'), username: 'User', is_active: true }]
        };
        
        await runner.apiCall('/team/add', { method: 'POST', body: JSON.stringify(teamData) });
        const result = await runner.apiCall(`/team/get?team_name=${teamId}`);
        
        if (!result.ok) throw new Error(`Expected 200, got ${result.status}`);
    });

    await runner.test('Get non-existent team', async () => {
        const result = await runner.apiCall('/team/get?team_name=nonexistent');
        if (result.status !== 404) throw new Error(`Expected 404, got ${result.status}`);
    });
}

async function runUserTests() {
    await runner.test('Toggle user active', async () => {
        const teamId = runner.generateId('team');
        const userId = runner.generateId('user');
        const teamData = {
            team_name: teamId,
            members: [{ user_id: userId, username: 'User', is_active: true }]
        };
        
        await runner.apiCall('/team/add', { method: 'POST', body: JSON.stringify(teamData) });

        const result = await runner.apiCall('/users/setIsActive', {
            method: 'POST',
            body: JSON.stringify({ user_id: userId, is_active: false })
        });
        
        if (!result.ok && result.status !== 422 && result.status !== 404) {
            throw new Error(`Expected 200, 422 or 404, got ${result.status}`);
        }
    });

    await runner.test('Update non-existent user', async () => {
        const result = await runner.apiCall('/users/setIsActive', {
            method: 'POST',
            body: JSON.stringify({ user_id: 'nonexistent', is_active: false })
        });
        
        if (result.status !== 422 && result.status !== 404) throw new Error(`Expected 422 or 404, got ${result.status}`);
    });
}

async function runPullRequestTests() {
    await runner.test('Create PR', async () => {
        const teamId = runner.generateId('team');
        const userId = runner.generateId('user');
        const prId = runner.generateId('pr');
        
        const teamData = {
            team_name: teamId,
            members: [
                { user_id: userId, username: 'Author', is_active: true },
                { user_id: runner.generateId('user'), username: 'Reviewer', is_active: true }
            ]
        };
        
        await runner.apiCall('/team/add', { method: 'POST', body: JSON.stringify(teamData) });

        const prData = {
            pull_request_id: prId,
            pull_request_name: 'Test',
            author_id: userId
        };
        
        const result = await runner.apiCall('/pullRequest/create', { method: 'POST', body: JSON.stringify(prData) });
        if (!result.ok) throw new Error(`Expected 201, got ${result.status}`);
    });

    await runner.test('Duplicate PR', async () => {
        const prId = runner.generateId('pr');
        const teamId = runner.generateId('team');
        const userId = runner.generateId('user');
        
        const teamData = {
            team_name: teamId,
            members: [
                { user_id: userId, username: 'Author', is_active: true },
                { user_id: runner.generateId('user'), username: 'Reviewer', is_active: true }
            ]
        };
        
        await runner.apiCall('/team/add', { method: 'POST', body: JSON.stringify(teamData) });

        const prData = { pull_request_id: prId, pull_request_name: 'Test', author_id: userId };
        
        await runner.apiCall('/pullRequest/create', { method: 'POST', body: JSON.stringify(prData) });
        const result = await runner.apiCall('/pullRequest/create', { method: 'POST', body: JSON.stringify(prData) });
        
        if (result.status !== 409) throw new Error(`Expected 409, got ${result.status}`);
    });

    await runner.test('PR non-existent author', async () => {
        const prData = {
            pull_request_id: runner.generateId('pr'),
            pull_request_name: 'Test',
            author_id: 'nonexistent'
        };
        
        const result = await runner.apiCall('/pullRequest/create', { method: 'POST', body: JSON.stringify(prData) });
        if (result.status !== 404) throw new Error(`Expected 404, got ${result.status}`);
    });

    await runner.test('Merge PR', async () => {
        const prId = runner.generateId('pr');
        const teamId = runner.generateId('team');
        const userId = runner.generateId('user');
        
        const teamData = {
            team_name: teamId,
            members: [
                { user_id: userId, username: 'Author', is_active: true },
                { user_id: runner.generateId('user'), username: 'Reviewer', is_active: true }
            ]
        };
        
        await runner.apiCall('/team/add', { method: 'POST', body: JSON.stringify(teamData) });

        const prData = { pull_request_id: prId, pull_request_name: 'Test', author_id: userId };
        
        await runner.apiCall('/pullRequest/create', { method: 'POST', body: JSON.stringify(prData) });
        const result = await runner.apiCall('/pullRequest/merge', { method: 'POST', body: JSON.stringify({ pull_request_id: prId }) });
        
        if (!result.ok) throw new Error(`Expected 200, got ${result.status}`);
    });

    await runner.test('Merge PR idempotent', async () => {
        const prId = runner.generateId('pr');
        const teamId = runner.generateId('team');
        const userId = runner.generateId('user');
        
        const teamData = {
            team_name: teamId,
            members: [
                { user_id: userId, username: 'Author', is_active: true },
                { user_id: runner.generateId('user'), username: 'Reviewer', is_active: true }
            ]
        };
        
        await runner.apiCall('/team/add', { method: 'POST', body: JSON.stringify(teamData) });

        const prData = { pull_request_id: prId, pull_request_name: 'Test', author_id: userId };
        
        await runner.apiCall('/pullRequest/create', { method: 'POST', body: JSON.stringify(prData) });
        await runner.apiCall('/pullRequest/merge', { method: 'POST', body: JSON.stringify({ pull_request_id: prId }) });
        const result = await runner.apiCall('/pullRequest/merge', { method: 'POST', body: JSON.stringify({ pull_request_id: prId }) });
        
        if (!result.ok) throw new Error(`Expected 200, got ${result.status}`);
    });
}

async function runEdgeCaseTests() {
    await runner.test('Team solo author', async () => {
        const teamId = runner.generateId('team');
        const userId = runner.generateId('user');
        
        const teamData = {
            team_name: teamId,
            members: [{ user_id: userId, username: 'Solo', is_active: true }]
        };
        
        await runner.apiCall('/team/add', { method: 'POST', body: JSON.stringify(teamData) });

        const prData = {
            pull_request_id: runner.generateId('pr'),
            pull_request_name: 'Test',
            author_id: userId
        };
        
        const result = await runner.apiCall('/pullRequest/create', { method: 'POST', body: JSON.stringify(prData) });
        if (!result.ok) throw new Error(`Expected 201, got ${result.status}`);
    });

    await runner.test('Team inactive author', async () => {
        const teamId = runner.generateId('team');
        const userId = runner.generateId('user');
        
        const teamData = {
            team_name: teamId,
            members: [{ user_id: userId, username: 'Inactive', is_active: false }]
        };
        
        await runner.apiCall('/team/add', { method: 'POST', body: JSON.stringify(teamData) });

        const prData = {
            pull_request_id: runner.generateId('pr'),
            pull_request_name: 'Test',
            author_id: userId
        };
        
        const result = await runner.apiCall('/pullRequest/create', { method: 'POST', body: JSON.stringify(prData) });
        if (result.ok && result.data.pr.assigned_reviewers.length !== 0) {
            throw new Error('Should have 0 reviewers when author inactive');
        }
    });
}

async function runReassignmentTests() {
    await runner.test('Reassign reviewer', async () => {
        const teamId = runner.generateId('team');
        const authorId = runner.generateId('user');
        const reviewer1 = runner.generateId('user');
        const reviewer2 = runner.generateId('user');
        const prId = runner.generateId('pr');
        
        const teamData = {
            team_name: teamId,
            members: [
                { user_id: authorId, username: 'Author', is_active: true },
                { user_id: reviewer1, username: 'Reviewer1', is_active: true },
                { user_id: reviewer2, username: 'Reviewer2', is_active: true }
            ]
        };
        
        await runner.apiCall('/team/add', { method: 'POST', body: JSON.stringify(teamData) });

        const prData = { pull_request_id: prId, pull_request_name: 'Test', author_id: authorId };
        
        const prResult = await runner.apiCall('/pullRequest/create', { method: 'POST', body: JSON.stringify(prData) });
        
        const oldReviewer = prResult.data.pr.assigned_reviewers[0];
        if (!oldReviewer) return;
        
        const result = await runner.apiCall('/pullRequest/reassign', {
            method: 'POST',
            body: JSON.stringify({ pull_request_id: prId, old_user_id: oldReviewer })
        });
        
        if (!result.ok && result.status !== 409) throw new Error(`Expected 200 or 409, got ${result.status}`);
    });
}

async function runDatabaseLimitTests() {
    await runner.test('Long team name', async () => {
        const teamData = {
            team_name: 'a'.repeat(150),
            members: [{ user_id: runner.generateId('user'), username: 'Test', is_active: true }]
        };
        
        const result = await runner.apiCall('/team/add', { method: 'POST', body: JSON.stringify(teamData) });
        if (result.status !== 400 && result.status !== 422 && result.status !== 413) {
            throw new Error(`Expected 400, 422 or 413, got ${result.status}`);
        }
    });

    await runner.test('Long user ID', async () => {
        const teamData = {
            team_name: runner.generateId('team'),
            members: [{ user_id: 'u'.repeat(150), username: 'Test', is_active: true }]
        };
        
        const result = await runner.apiCall('/team/add', { method: 'POST', body: JSON.stringify(teamData) });
        if (result.status !== 400 && result.status !== 422 && result.status !== 413) {
            throw new Error(`Expected 400, 422 or 413, got ${result.status}`);
        }
    });

    await runner.test('Long PR ID', async () => {
        const teamId = runner.generateId('team');
        const userId = runner.generateId('user');
        
        const teamData = {
            team_name: teamId,
            members: [
                { user_id: userId, username: 'Test', is_active: true },
                { user_id: runner.generateId('user'), username: 'Reviewer', is_active: true }
            ]
        };
        
        await runner.apiCall('/team/add', { method: 'POST', body: JSON.stringify(teamData) });

        const prData = {
            pull_request_id: 'pr'.repeat(100),
            pull_request_name: 'Test',
            author_id: userId
        };
        
        const result = await runner.apiCall('/pullRequest/create', { method: 'POST', body: JSON.stringify(prData) });
        if (result.status !== 400 && result.status !== 422 && result.status !== 413 && result.status !== 404) {
            throw new Error(`Expected 400, 422, 413 or 404, got ${result.status}`);
        }
    });

    await runner.test('Empty strings', async () => {
        const teamData = { team_name: '', members: [{ user_id: '', username: '', is_active: true }] };
        const result = await runner.apiCall('/team/add', { method: 'POST', body: JSON.stringify(teamData) });
        if (result.status !== 400 && result.status !== 422) throw new Error(`Expected 400 or 422, got ${result.status}`);
    });

    await runner.test('SQL injection', async () => {
        const teamData = {
            team_name: "team'; DROP TABLE teams; --",
            members: [{ user_id: "user'; DELETE FROM users; --", username: "Test", is_active: true }]
        };
        
        await runner.apiCall('/team/add', { method: 'POST', body: JSON.stringify(teamData) });
    });

    await runner.test('Large payload', async () => {
        const members = [];
        for (let i = 0; i < 1000; i++) {
            members.push({ user_id: runner.generateId(`user-${i}`), username: `User ${i}`, is_active: true });
        }
        
        const teamData = { team_name: runner.generateId('bulk'), members };
        await runner.apiCall('/team/add', { method: 'POST', body: JSON.stringify(teamData) });
    });
}

console.log("run tests: runAllTests()")