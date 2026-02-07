CREATE TABLE events (
    id                    BIGSERIAL PRIMARY KEY,
    title                 VARCHAR(255) NOT NULL,
    description           TEXT,
    event_date            TIMESTAMP NOT NULL,
    end_date              TIMESTAMP,
    location              VARCHAR(255),
    location_details      TEXT,
    is_recurring          BOOLEAN DEFAULT FALSE,
    recurrence_rule       VARCHAR(20) DEFAULT 'none',
    recurrence_end        DATE,
    registration_enabled  BOOLEAN DEFAULT FALSE,
    capacity_limit        INTEGER,
    registration_deadline TIMESTAMP,
    visible_from          TIMESTAMP,
    visible_until         TIMESTAMP,
    is_public             BOOLEAN DEFAULT TRUE,
    ministry_id           INTEGER,
    created_by            BIGINT,
    created_at            TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at            TIMESTAMP
);

CREATE INDEX idx_events_event_date ON events(event_date);
CREATE INDEX idx_events_is_public ON events(is_public);
CREATE INDEX idx_events_deleted_at ON events(deleted_at);
