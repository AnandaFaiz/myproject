-- +goose Up
CREATE TABLE berita (
    id SERIAL PRIMARY KEY,
    judul TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    konten TEXT NOT NULL,
    penulis TEXT NOT NULL,
    kategori TEXT,
    gambar TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE berita;
