CREATE TABLE staff_members (
    id             BIGSERIAL PRIMARY KEY,
    user_id        BIGINT,
    name           VARCHAR(255) NOT NULL,
    title          VARCHAR(255) NOT NULL,
    bio            TEXT,
    email          VARCHAR(255),
    phone          VARCHAR(20),
    photo_url      VARCHAR(500),
    display_order  INTEGER DEFAULT 0,
    is_active      BOOLEAN DEFAULT TRUE,
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at     TIMESTAMP
);

CREATE INDEX idx_staff_members_is_active ON staff_members(is_active);
CREATE INDEX idx_staff_members_display_order ON staff_members(display_order);
CREATE INDEX idx_staff_members_deleted_at ON staff_members(deleted_at);
