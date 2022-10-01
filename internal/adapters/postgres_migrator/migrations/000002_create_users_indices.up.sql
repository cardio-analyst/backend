CREATE INDEX IF NOT EXISTS users_login_password_hash_index ON users (login, password_hash);
CREATE INDEX IF NOT EXISTS users_email_password_hash_index ON users (email, password_hash);
