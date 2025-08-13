CREATE TABLE IF NOT EXISTS shortened_links (
    id              SERIAL PRIMARY KEY,

    short_code      VARCHAR(255) NOT NULL UNIQUE,
    original_url    TEXT NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    expires_at      TIMESTAMP WITH TIME ZONE,

    attributes      JSONB NOT NULL DEFAULT '{}'
);

CREATE INDEX IF NOT EXISTS idx_shortened_links_short_code ON shortened_links (short_code);
CREATE INDEX IF NOT EXISTS idx_shortened_links_expires_at ON shortened_links (expires_at);
CREATE INDEX IF NOT EXISTS idx_shortened_links_created_at ON shortened_links (created_at);
