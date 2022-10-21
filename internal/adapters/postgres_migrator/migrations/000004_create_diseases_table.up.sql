CREATE TABLE IF NOT EXISTS diseases
(
    id                    SERIAL PRIMARY KEY,
    user_id               INTEGER UNIQUE NOT NULL,
    cvds_predisposition   BOOLEAN        DEFAULT FALSE,
    take_statins          BOOLEAN        DEFAULT FALSE,
    ckd                   BOOLEAN        DEFAULT FALSE,
    arterial_hypertension BOOLEAN        DEFAULT FALSE,
    cardiac_ischemia      BOOLEAN        DEFAULT FALSE,
    type_two_diabets      BOOLEAN        DEFAULT FALSE,
    infarction_or_stroke  BOOLEAN        DEFAULT FALSE,
    atherosclerosis       BOOLEAN        DEFAULT FALSE,
    other_cvds_diseases   BOOLEAN        DEFAULT FALSE,
    CONSTRAINT fk_diseases_users FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
