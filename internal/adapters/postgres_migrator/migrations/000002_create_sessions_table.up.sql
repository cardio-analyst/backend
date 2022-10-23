CREATE TABLE IF NOT EXISTS sessions
(
    user_id       INTEGER PRIMARY KEY,
    refresh_token VARCHAR(255)   NOT NULL,
    whitelist     VARCHAR(255)[] NOT NULL,
    CONSTRAINT fk_sessions_users FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
