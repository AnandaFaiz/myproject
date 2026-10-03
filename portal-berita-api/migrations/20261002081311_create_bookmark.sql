-- +goose Up
CREATE TABLE bookmark (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    berita_id INTEGER NOT NULL REFERENCES berita(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, berita_id)
);
CREATE INDEX idx_bookmark_user_id ON bookmark(user_id);
CREATE INDEX idx_bookmark_berita_id ON bookmark(berita_id);

-- +goose Down
DROP TABLE bookmark;
