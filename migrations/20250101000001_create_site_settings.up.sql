CREATE TABLE site_settings (
    key         VARCHAR(100) PRIMARY KEY,
    value       TEXT,
    description TEXT,
    updated_by  BIGINT,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
