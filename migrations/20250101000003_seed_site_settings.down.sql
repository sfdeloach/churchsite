DELETE FROM site_settings WHERE key IN (
    'church_name',
    'church_address',
    'church_phone',
    'church_email',
    'morning_service_time',
    'morning_service_name',
    'evening_service_time',
    'evening_service_name',
    'wednesday_service_time',
    'wednesday_service_name',
    'sunday_school_time',
    'sunday_school_name',
    'hero_title',
    'hero_subtitle'
);
