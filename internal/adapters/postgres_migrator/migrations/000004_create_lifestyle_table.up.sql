CREATE TABLE IF NOT EXISTS lifestyles
(
    user_id                    INTEGER PRIMARY KEY,
    family_status              VARCHAR(255) DEFAULT '',
    events_participation       VARCHAR(255) DEFAULT '',
    physical_activity          VARCHAR(255) DEFAULT '',
    work_status                VARCHAR(255) DEFAULT '',
    significant_value_high     VARCHAR(255) DEFAULT '',
    significant_value_medium   VARCHAR(255) DEFAULT '',
    significant_value_low      VARCHAR(255) DEFAULT '',
    CONSTRAINT fk_lifestyles_users FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
