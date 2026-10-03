package models

import "time"

type Bookmark struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	BeritaID  int       `json:"berita_id"`
	CreatedAt time.Time `json:"createdAt"`
}

type BookmarkDenganBerita struct {
	ID int `json:"id"`
	UserID    int       `json:"user_id"`
	BeritaID  int       `json:"berita_id"`
	CreatedAt time.Time `json:"createdAt"`
	
	//detail berita
	Judul     string    `json:"judul"`
	Slug      string    `json:"slug"`
	Konten   string    `json:"konten"`
	Penulis   string    `json:"penulis"`
	Kategori  *string    `json:"kategori"`
	Gambar    *string    `json:"gambar"`
	BeritaCreatedAt time.Time `json:"beritaCreatedAt"`
}