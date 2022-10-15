CREATE TABLE IF NOT EXISTS sessions
(
    id            SERIAL PRIMARY KEY,
    user_id       INTEGER UNIQUE NOT NULL,
    refresh_token VARCHAR(255)   NOT NULL,
    whitelist     VARCHAR(255)[] NOT NULL,
    CONSTRAINT fk_sessions_users FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
