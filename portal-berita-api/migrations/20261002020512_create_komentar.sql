-- +goose Up
CREATE TABLE komentar(
    id SERIAL PRIMARY KEY,
    berita_id INTEGER NOT NULL REFERENCES berita(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    isi TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_komentar_berita_id ON komentar(berita_id);
CREATE INDEX idx_komentar_user_id ON komentar(user_id);

-- +goose Down
DROP TABLE komentar;
