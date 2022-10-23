CREATE TABLE IF NOT EXISTS diseases
(
    user_id                    INTEGER PRIMARY KEY,
    cvd_predisposed            BOOLEAN DEFAULT FALSE,
    takes_statins              BOOLEAN DEFAULT FALSE,
    has_chronic_kidney_disease BOOLEAN DEFAULT FALSE,
    has_arterial_hypertension  BOOLEAN DEFAULT FALSE,
    has_ischemic_heart_disease BOOLEAN DEFAULT FALSE,
    has_type_two_diabetes      BOOLEAN DEFAULT FALSE,
    had_infarction_or_stroke   BOOLEAN DEFAULT FALSE,
    has_atherosclerosis        BOOLEAN DEFAULT FALSE,
    has_other_cvd              BOOLEAN DEFAULT FALSE,
    CONSTRAINT fk_diseases_users FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
