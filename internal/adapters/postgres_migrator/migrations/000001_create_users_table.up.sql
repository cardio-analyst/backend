CREATE TABLE IF NOT EXISTS users
(
    id            SERIAL PRIMARY KEY,
    first_name    VARCHAR(100) NOT NULL,
    last_name     VARCHAR(100) NOT NULL,
    middle_name   VARCHAR(100),
    region        VARCHAR(100) NOT NULL,
    birth_date    DATE         NOT NULL,
    gender        VARCHAR(100) NOT NULL,
    login         VARCHAR(100) NOT NULL,
    email         VARCHAR(100) NOT NULL,
    password_hash VARCHAR(100) NOT NULL
);
