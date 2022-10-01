CREATE TABLE IF NOT EXISTS users
(
    id            SERIAL PRIMARY KEY,
    first_name    VARCHAR(255)        NOT NULL,
    last_name     VARCHAR(255)        NOT NULL,
    middle_name   VARCHAR(255),
    region        VARCHAR(255)        NOT NULL,
    birth_date    DATE                NOT NULL,
    gender        VARCHAR(255)        NOT NULL,
    login         VARCHAR(255) UNIQUE NOT NULL,
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255)        NOT NULL
);
