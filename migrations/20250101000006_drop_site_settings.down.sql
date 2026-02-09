CREATE TABLE site_settings (
    key         VARCHAR(100) PRIMARY KEY,
    value       TEXT,
    description TEXT,
    updated_by  BIGINT,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO site_settings (key, value, description) VALUES
    ('church_name', 'Saint Andrew''s Chapel', 'Name of the church'),
    ('church_address', '123 Main Street, Anytown, USA 12345', 'Physical address'),
    ('church_phone', '(555) 123-4567', 'Main phone number'),
    ('church_email', 'info@sachapel.com', 'Main contact email'),
    ('morning_service_time', '10:30 AM', 'Morning worship service time'),
    ('morning_service_name', 'Morning Worship', 'Name of the morning service'),
    ('evening_service_time', '6:00 PM', 'Evening worship service time'),
    ('evening_service_name', 'Evening Worship', 'Name of the evening service'),
    ('wednesday_service_time', '7:00 PM', 'Wednesday service time'),
    ('wednesday_service_name', 'Prayer Meeting & Bible Study', 'Name of the Wednesday service'),
    ('sunday_school_time', '9:15 AM', 'Sunday School time'),
    ('sunday_school_name', 'Sunday School', 'Name of the Sunday School hour'),
    ('hero_title', 'Welcome to Saint Andrew''s Chapel', 'Homepage hero banner title'),
    ('hero_subtitle', 'A congregation committed to the glory of God and the good of His people.', 'Homepage hero banner subtitle');
