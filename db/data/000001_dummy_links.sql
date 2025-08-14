INSERT INTO shortened_links (short_code, original_url, expires_at) VALUES
    ('abc123', 'https://example.com/long-url-1', NULL),
    ('def456', 'https://example.com/long-url-2', NULL),
    ('ghi789', 'https://example.com/long-url-3', '2200-12-31 23:59:59+00'), -- does not expire
    ('jkl012', 'https://example.com/long-url-4', '2020-01-01 00:00:00+00'), -- already expired
    ('mno345', 'https://example.com/long-url-5', NULL);
