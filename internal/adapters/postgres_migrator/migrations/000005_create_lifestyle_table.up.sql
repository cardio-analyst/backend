CREATE TABLE IF NOT EXISTS lifestyles
(
    user_id                  INTEGER PRIMARY KEY,
    family_status            VARCHAR(255) DEFAULT '' NOT NULL,
    events_participation     VARCHAR(255) DEFAULT '' NOT NULL,
    physical_activity        VARCHAR(255) DEFAULT '' NOT NULL,
    work_status              VARCHAR(255) DEFAULT '' NOT NULL,
    significant_value_high   VARCHAR(255) DEFAULT '' NOT NULL,
    significant_value_medium VARCHAR(255) DEFAULT '' NOT NULL,
    significant_value_low    VARCHAR(255) DEFAULT '' NOT NULL,
    CONSTRAINT fk_lifestyles_users FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
