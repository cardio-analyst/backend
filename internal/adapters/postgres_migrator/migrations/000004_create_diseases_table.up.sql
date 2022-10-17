CREATE TABLE IF NOT EXISTS diseases
(
    id                    SERIAL PRIMARY KEY,
    user_id               INTEGER UNIQUE NOT NULL,
    cvds_predisposition   VARCHAR(255)   NOT NULL,
    take_statins          BOOLEAN        NOT NULL,
    ckd                   BOOLEAN        NOT NULL,
    arterial_hypertension BOOLEAN        NOT NULL,
    cardiac_ischemia      BOOLEAN        NOT NULL,
    type_two_diabets      BOOLEAN        NOT NULL,
    infarction_or_stroke  VARCHAR(255)   NOT NULL,
    atherosclerosis       BOOLEAN        NOT NULL,
    other_cvds_diseases   VARCHAR(255),
    CONSTRAINT fk_diseases_users FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
