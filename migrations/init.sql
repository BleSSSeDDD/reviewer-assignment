CREATE TABLE users(
    user_id VARCHAR(100) PRIMARY KEY, -- длина до 100 символов чтобы обезопасить бд от огромных строк
    user_name VARCHAR(100) NOT NULL, 
    is_active BOOLEAN NOT NULL DEFAULT true
);

CREATE TABLE teams(
    team_name VARCHAR(100) PRIMARY KEY
);

-- связывает юзеров с их командами, один юзер может быть в нескольких командах
CREATE TABLE members_of_teams(
    user_id VARCHAR(100) REFERENCES users(user_id),
    team_name VARCHAR(100) REFERENCES teams(team_name),
    PRIMARY KEY (user_id, team_name)
);

-- удаления юзеров и команд не предусмотрено
CREATE TABLE pull_requests(
    request_id VARCHAR(100) PRIMARY KEY,
    title VARCHAR(500) NOT NULL,
    user_id VARCHAR(100) REFERENCES users(user_id), 
    status VARCHAR(6) CHECK (status IN ('OPEN', 'MERGED')) DEFAULT 'OPEN',
    created_at TIMESTAMP DEFAULT NOW(),
    merged_at TIMESTAMP
);

-- связывает пулл реквесты с ревьюервами
-- у одного реквеста может быть много ревьюеров (в коде будет ограничение на макисмум 2)
-- у одного ревьюера может быть много реквестов
CREATE TABLE pull_requests_reviewers(
    request_id VARCHAR(100) REFERENCES pull_requests(request_id),
    reviewer_id VARCHAR(100) REFERENCES users(user_id),
    PRIMARY KEY (request_id, reviewer_id)
);

