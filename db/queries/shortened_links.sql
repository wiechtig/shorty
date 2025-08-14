-- name: ResolveShortenedLink :one
SELECT id, original_url FROM shortened_links
WHERE
    LOWER(short_code) = LOWER($1)
    AND (expires_at IS NULL OR expires_at > NOW())
LIMIT 1;
