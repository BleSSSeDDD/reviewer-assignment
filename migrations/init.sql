-- таблица для команд нужна, чтобы не возникло ситуации, когда два одновременных запроса с одинаковой
-- командой и разными пользователями не создавали две разные группы пользователей с одинковым названием команды
CREATE TABLE teams(
    team_name VARCHAR(100) PRIMARY KEY
);


CREATE TABLE users(
    user_id VARCHAR(100) PRIMARY KEY, -- длина до 100 символов чтобы обезопасить бд от огромных строк
    user_name VARCHAR(100) NOT NULL, 
    team_name VARCHAR(100) NOT NULL REFERENCES teams(team_name) DEFERRABLE INITIALLY DEFERRED,
    is_active BOOLEAN NOT NULL DEFAULT true
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

