CREATE TABLE ministries (
    id            BIGSERIAL PRIMARY KEY,
    name          VARCHAR(255) NOT NULL,
    slug          VARCHAR(255) UNIQUE NOT NULL,
    description   TEXT,
    leader_id     BIGINT,
    contact_email VARCHAR(255),
    meeting_time  VARCHAR(255),
    location      VARCHAR(255),
    is_active     BOOLEAN DEFAULT TRUE,
    sort_order    INTEGER DEFAULT 0,
    page_content  TEXT,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMP
);

CREATE INDEX idx_ministries_slug ON ministries(slug);
CREATE INDEX idx_ministries_is_active ON ministries(is_active);
CREATE INDEX idx_ministries_sort_order ON ministries(sort_order);
CREATE INDEX idx_ministries_deleted_at ON ministries(deleted_at);
