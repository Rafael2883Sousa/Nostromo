-- +goose Up
CREATE TABLE IF NOT EXISTS m365_tenants (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tenant_id TEXT NOT NULL UNIQUE,
    display_name TEXT,
    consented_at TIMESTAMP,
    access_token TEXT,
    refresh_token TEXT,
    expires_at TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS m365_tenants;
