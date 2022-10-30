CREATE TABLE IF NOT EXISTS analyses
(
    id                                  SERIAL PRIMARY KEY,
    user_id                             INTEGER                 NOT NULL,
    high_density_cholesterol            DECIMAL(10, 2),
    low_density_cholesterol             DECIMAL(10, 2),
    triglycerides                       DECIMAL(10, 2),
    lipoprotein                         DECIMAL(10, 2),
    highly_sensitive_c_reactive_protein DECIMAL(10, 2),
    atherogenicity_coefficient          DECIMAL(10, 2),
    creatinine                          DECIMAL(10, 2),
    atherosclerotic_plaques_presence    BOOLEAN,
    created_at                          TIMESTAMP DEFAULT NOW() NOT NULL,
    CONSTRAINT fk_analysis_users FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
