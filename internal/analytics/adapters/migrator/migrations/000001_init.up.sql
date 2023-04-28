CREATE TABLE IF NOT EXISTS feedback
(
    user_id          INTEGER      NOT NULL,
    user_first_name  VARCHAR(255) NOT NULL,
    user_last_name   VARCHAR(255) NOT NULL,
    user_middle_name VARCHAR(255),
    user_login       VARCHAR(255) NOT NULL,
    user_email       VARCHAR(255) NOT NULL,
    mark             SMALLINT     NOT NULL,
    message          TEXT
);
