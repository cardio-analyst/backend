CREATE TABLE IF NOT EXISTS very_high_risk_female_not_smoking
(
    age_min                     INTEGER       NOT NULL,
    age_max                     INTEGER       NOT NULL,
    systolic_blood_pressure_min INTEGER       NOT NULL,
    systolic_blood_pressure_max INTEGER       NOT NULL,
    non_hdl_cholesterol_min     DECIMAL(6, 5) NOT NULL,
    non_hdl_cholesterol_max     DECIMAL(6, 5) NOT NULL,
    risk_value                  INTEGER       NOT NULL
);

CREATE TABLE IF NOT EXISTS very_high_risk_male_not_smoking
(
    age_min                     INTEGER       NOT NULL,
    age_max                     INTEGER       NOT NULL,
    systolic_blood_pressure_min INTEGER       NOT NULL,
    systolic_blood_pressure_max INTEGER       NOT NULL,
    non_hdl_cholesterol_min     DECIMAL(6, 5) NOT NULL,
    non_hdl_cholesterol_max     DECIMAL(6, 5) NOT NULL,
    risk_value                  INTEGER       NOT NULL
);

CREATE TABLE IF NOT EXISTS very_high_risk_female_smoking
(
    age_min                     INTEGER       NOT NULL,
    age_max                     INTEGER       NOT NULL,
    systolic_blood_pressure_min INTEGER       NOT NULL,
    systolic_blood_pressure_max INTEGER       NOT NULL,
    non_hdl_cholesterol_min     DECIMAL(6, 5) NOT NULL,
    non_hdl_cholesterol_max     DECIMAL(6, 5) NOT NULL,
    risk_value                  INTEGER       NOT NULL
);

CREATE TABLE IF NOT EXISTS very_high_risk_male_smoking
(
    age_min                     INTEGER       NOT NULL,
    age_max                     INTEGER       NOT NULL,
    systolic_blood_pressure_min INTEGER       NOT NULL,
    systolic_blood_pressure_max INTEGER       NOT NULL,
    non_hdl_cholesterol_min     DECIMAL(6, 5) NOT NULL,
    non_hdl_cholesterol_max     DECIMAL(6, 5) NOT NULL,
    risk_value                  INTEGER       NOT NULL
);

CREATE TABLE IF NOT EXISTS high_risk_female_not_smoking
(
    age_min                     INTEGER       NOT NULL,
    age_max                     INTEGER       NOT NULL,
    systolic_blood_pressure_min INTEGER       NOT NULL,
    systolic_blood_pressure_max INTEGER       NOT NULL,
    non_hdl_cholesterol_min     DECIMAL(6, 5) NOT NULL,
    non_hdl_cholesterol_max     DECIMAL(6, 5) NOT NULL,
    risk_value                  INTEGER       NOT NULL
);

CREATE TABLE IF NOT EXISTS high_risk_male_not_smoking
(
    age_min                     INTEGER       NOT NULL,
    age_max                     INTEGER       NOT NULL,
    systolic_blood_pressure_min INTEGER       NOT NULL,
    systolic_blood_pressure_max INTEGER       NOT NULL,
    non_hdl_cholesterol_min     DECIMAL(6, 5) NOT NULL,
    non_hdl_cholesterol_max     DECIMAL(6, 5) NOT NULL,
    risk_value                  INTEGER       NOT NULL
);

CREATE TABLE IF NOT EXISTS high_risk_female_smoking
(
    age_min                     INTEGER       NOT NULL,
    age_max                     INTEGER       NOT NULL,
    systolic_blood_pressure_min INTEGER       NOT NULL,
    systolic_blood_pressure_max INTEGER       NOT NULL,
    non_hdl_cholesterol_min     DECIMAL(6, 5) NOT NULL,
    non_hdl_cholesterol_max     DECIMAL(6, 5) NOT NULL,
    risk_value                  INTEGER       NOT NULL
);

CREATE TABLE IF NOT EXISTS high_risk_male_smoking
(
    age_min                     INTEGER       NOT NULL,
    age_max                     INTEGER       NOT NULL,
    systolic_blood_pressure_min INTEGER       NOT NULL,
    systolic_blood_pressure_max INTEGER       NOT NULL,
    non_hdl_cholesterol_min     DECIMAL(6, 5) NOT NULL,
    non_hdl_cholesterol_max     DECIMAL(6, 5) NOT NULL,
    risk_value                  INTEGER       NOT NULL
);
