package models

import "time"

type Berita struct {
	ID        int       `json:"id"`
	Judul     string    `json:"judul"`
	Slug      string    `json:"slug"`
	Konten    string    `json:"konten"`
	Penulis   string    `json:"penulis"`
	Kategori  *string   `json:"kategori"`
	Gambar    *string   `json:"gambar"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updateAt"`
}
