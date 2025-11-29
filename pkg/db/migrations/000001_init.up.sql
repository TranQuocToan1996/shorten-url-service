CREATE TABLE IF NOT EXISTS shorten_urls (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    code VARCHAR(255) Unique NULL,
    status VARCHAR(255) NOT NULL DEFAULT 'submitted', -- submitted, encoded, failed
    algo VARCHAR(255) NULL,
    clean_url varchar(2048) Unique NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_shorten_urls_status_created_at ON shorten_urls(status, created_at DESC);