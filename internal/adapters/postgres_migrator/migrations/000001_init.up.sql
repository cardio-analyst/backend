CREATE TABLE IF NOT EXISTS users
(
    id            SERIAL PRIMARY KEY,
    first_name    VARCHAR(255)        NOT NULL,
    last_name     VARCHAR(255)        NOT NULL,
    middle_name   VARCHAR(255),
    region        VARCHAR(255)        NOT NULL,
    birth_date    DATE                NOT NULL,
    login         VARCHAR(255) UNIQUE NOT NULL,
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255)        NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions
(
    user_id       INTEGER PRIMARY KEY,
    refresh_token VARCHAR(255)   NOT NULL,
    whitelist     VARCHAR(255)[] NOT NULL,
    CONSTRAINT fk_sessions_users FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS diseases
(
    user_id                    INTEGER PRIMARY KEY,
    cvd_predisposed            BOOLEAN NOT NULL DEFAULT FALSE,
    takes_statins              BOOLEAN NOT NULL DEFAULT FALSE,
    has_chronic_kidney_disease BOOLEAN NOT NULL DEFAULT FALSE,
    has_arterial_hypertension  BOOLEAN NOT NULL DEFAULT FALSE,
    has_ischemic_heart_disease BOOLEAN NOT NULL DEFAULT FALSE,
    has_type_two_diabetes      BOOLEAN NOT NULL DEFAULT FALSE,
    had_infarction_or_stroke   BOOLEAN NOT NULL DEFAULT FALSE,
    has_atherosclerosis        BOOLEAN NOT NULL DEFAULT FALSE,
    has_other_cvd              BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT fk_diseases_users FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS analyses
(
    id                                  SERIAL PRIMARY KEY,
    user_id                             INTEGER                 NOT NULL,
    high_density_cholesterol            DECIMAL(2, 2),
    low_density_cholesterol             DECIMAL(2, 2),
    triglycerides                       DECIMAL(2, 2),
    lipoprotein                         DECIMAL(2, 2),
    highly_sensitive_c_reactive_protein DECIMAL(2, 2),
    atherogenicity_coefficient          DECIMAL(2, 2),
    creatinine                          DECIMAL(4, 2),
    atherosclerotic_plaques_presence    BOOLEAN,
    created_at                          TIMESTAMP DEFAULT NOW() NOT NULL,
    CONSTRAINT fk_analysis_users FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS lifestyles
(
    user_id                   INTEGER PRIMARY KEY,
    family_status             VARCHAR(255)  NOT NULL DEFAULT '',
    events_participation      VARCHAR(255)  NOT NULL DEFAULT '',
    physical_activity         VARCHAR(255)  NOT NULL DEFAULT '',
    work_status               VARCHAR(255)  NOT NULL DEFAULT '',
    significant_value_high    VARCHAR(255)  NOT NULL DEFAULT '',
    significant_value_medium  VARCHAR(255)  NOT NULL DEFAULT '',
    significant_value_low     VARCHAR(255)  NOT NULL DEFAULT '',
    angina_score              INTEGER       NOT NULL DEFAULT -1,
    adherence_drug_therapy    DECIMAL(4, 2) NOT NULL DEFAULT -1.0,
    adherence_medical_support DECIMAL(4, 2) NOT NULL DEFAULT -1.0,
    adherence_lifestyle_mod   DECIMAL(4, 2) NOT NULL DEFAULT -1.0,
    CONSTRAINT fk_lifestyles_users FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS basic_indicators
(
    id                       SERIAL PRIMARY KEY,
    user_id                  INTEGER                 NOT NULL,
    weight                   DECIMAL(4, 1),
    height                   DECIMAL(4, 1),
    body_mass_index          DECIMAL(4, 1),
    waist_size               DECIMAL(4, 1),
    gender                   VARCHAR(255),
    sbp_level                DECIMAL(4, 1),
    smoking                  BOOLEAN,
    total_cholesterol_level  DECIMAL(4, 1),
    cv_events_risk_value     INTEGER,
    ideal_cardiovascular_age INTEGER,
    created_at               TIMESTAMP DEFAULT NOW() NOT NULL,
    CONSTRAINT fk_basic_indicators_users FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS very_high_risk_female_not_smoking
(
    age_min                     INTEGER       NOT NULL,
    age_max                     INTEGER       NOT NULL,
    systolic_blood_pressure_min INTEGER       NOT NULL,
    systolic_blood_pressure_max INTEGER       NOT NULL,
    non_hdl_cholesterol_min     DECIMAL(3, 1) NOT NULL,
    non_hdl_cholesterol_max     DECIMAL(3, 1) NOT NULL,
    risk_value                  INTEGER       NOT NULL
);

CREATE TABLE IF NOT EXISTS very_high_risk_male_not_smoking
(
    age_min                     INTEGER       NOT NULL,
    age_max                     INTEGER       NOT NULL,
    systolic_blood_pressure_min INTEGER       NOT NULL,
    systolic_blood_pressure_max INTEGER       NOT NULL,
    non_hdl_cholesterol_min     DECIMAL(3, 1) NOT NULL,
    non_hdl_cholesterol_max     DECIMAL(3, 1) NOT NULL,
    risk_value                  INTEGER       NOT NULL
);

CREATE TABLE IF NOT EXISTS very_high_risk_female_smoking
(
    age_min                     INTEGER       NOT NULL,
    age_max                     INTEGER       NOT NULL,
    systolic_blood_pressure_min INTEGER       NOT NULL,
    systolic_blood_pressure_max INTEGER       NOT NULL,
    non_hdl_cholesterol_min     DECIMAL(3, 1) NOT NULL,
    non_hdl_cholesterol_max     DECIMAL(3, 1) NOT NULL,
    risk_value                  INTEGER       NOT NULL
);

CREATE TABLE IF NOT EXISTS very_high_risk_male_smoking
(
    age_min                     INTEGER       NOT NULL,
    age_max                     INTEGER       NOT NULL,
    systolic_blood_pressure_min INTEGER       NOT NULL,
    systolic_blood_pressure_max INTEGER       NOT NULL,
    non_hdl_cholesterol_min     DECIMAL(3, 1) NOT NULL,
    non_hdl_cholesterol_max     DECIMAL(3, 1) NOT NULL,
    risk_value                  INTEGER       NOT NULL
);

CREATE TABLE IF NOT EXISTS high_risk_female_not_smoking
(
    age_min                     INTEGER       NOT NULL,
    age_max                     INTEGER       NOT NULL,
    systolic_blood_pressure_min INTEGER       NOT NULL,
    systolic_blood_pressure_max INTEGER       NOT NULL,
    non_hdl_cholesterol_min     DECIMAL(3, 1) NOT NULL,
    non_hdl_cholesterol_max     DECIMAL(3, 1) NOT NULL,
    risk_value                  INTEGER       NOT NULL
);

CREATE TABLE IF NOT EXISTS high_risk_male_not_smoking
(
    age_min                     INTEGER       NOT NULL,
    age_max                     INTEGER       NOT NULL,
    systolic_blood_pressure_min INTEGER       NOT NULL,
    systolic_blood_pressure_max INTEGER       NOT NULL,
    non_hdl_cholesterol_min     DECIMAL(3, 1) NOT NULL,
    non_hdl_cholesterol_max     DECIMAL(3, 1) NOT NULL,
    risk_value                  INTEGER       NOT NULL
);

CREATE TABLE IF NOT EXISTS high_risk_female_smoking
(
    age_min                     INTEGER       NOT NULL,
    age_max                     INTEGER       NOT NULL,
    systolic_blood_pressure_min INTEGER       NOT NULL,
    systolic_blood_pressure_max INTEGER       NOT NULL,
    non_hdl_cholesterol_min     DECIMAL(3, 1) NOT NULL,
    non_hdl_cholesterol_max     DECIMAL(3, 1) NOT NULL,
    risk_value                  INTEGER       NOT NULL
);

CREATE TABLE IF NOT EXISTS high_risk_male_smoking
(
    age_min                     INTEGER       NOT NULL,
    age_max                     INTEGER       NOT NULL,
    systolic_blood_pressure_min INTEGER       NOT NULL,
    systolic_blood_pressure_max INTEGER       NOT NULL,
    non_hdl_cholesterol_min     DECIMAL(3, 1) NOT NULL,
    non_hdl_cholesterol_max     DECIMAL(3, 1) NOT NULL,
    risk_value                  INTEGER       NOT NULL
);
