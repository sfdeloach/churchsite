ALTER TABLE staff_members
    ADD COLUMN category VARCHAR(50) NOT NULL DEFAULT 'staff';

UPDATE staff_members SET category = 'pastor' WHERE title LIKE '%Pastor%';

CREATE INDEX idx_staff_members_category ON staff_members(category);
