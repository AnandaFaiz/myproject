package handlers

import (
	"context"
	"net/http"
	"portal-berita-api/config"
	"portal-berita-api/models"
	"portal-berita-api/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetBerita(c *gin.Context) {
	rows, err := config.DB.Query(
		context.Background(),
		`SELECT id, judul, slug, konten, penulis, kategori, gambar, created_at, updated_at
		FROM berita
		ORDER BY created_at DESC`,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"sukses": false,
			"pesan":  "Gagal mengambil data berita",
			"error":  err.Error(),
		})
		return
	}

	defer rows.Close()

	//Slice tampung hasil
	daftarBerita := []models.Berita{}

	//Loop setiap baris
	for rows.Next() {
		var b models.Berita
		err := rows.Scan(
			&b.ID,
			&b.Judul,
			&b.Slug,
			&b.Konten,
			&b.Penulis,
			&b.Kategori,
			&b.Gambar,
			&b.CreatedAt,
			&b.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"sukses": false,
				"pesan":  "Gagal membaca data",
				"error":  err.Error(),
			})
			return
		}
		daftarBerita = append(daftarBerita, b)
	}

	//return JSON
	c.JSON(http.StatusOK, gin.H{
		"sukses": true,
		"data":   daftarBerita,
	})
}

func CreateBerita(c *gin.Context) {
	// struct untuk body req
	var input struct {
		Judul    string `json:"judul"`
		Konten   string `json:"konten"`
		Penulis  string `json:"penulis"`
		Kategori string `json:"kategori"`
		Gambar   string `json:"gambar"`
	}

	// Bind JSON dari body
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"sukses": false,
			"pesan":  "Data tidak valid",
			"error":  err.Error(),
		})
		return
	}

	// Validasi wajib
	if input.Judul == "" || input.Konten == "" || input.Penulis == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"sukses": false,
			"pesan":  "Judul, konten, dan penulis wajib diisi",
		})
		return
	}

	// Generate slug
	slug := utils.BuatSlug(input.Judul)

	// Siapkan nilai nullable
	var kategori *string
	if input.Kategori != "" {
		kategori = &input.Kategori
	}

	var gambar *string
	if input.Gambar != "" {
		gambar = &input.Gambar
	}

	// Insert to DB
	var b models.Berita
	err := config.DB.QueryRow(
		context.Background(),
		`INSERT INTO berita (judul, slug, konten, penulis, kategori, gambar)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, judul, slug, konten, penulis, kategori, gambar, created_at, updated_at`,
		input.Judul,
		slug,
		input.Konten,
		input.Penulis,
		kategori,
		gambar,
	).Scan(
		&b.ID,
		&b.Judul,
		&b.Slug,
		&b.Konten,
		&b.Penulis,
		&b.Kategori,
		&b.Gambar,
		&b.CreatedAt,
		&b.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"sukses": false,
			"pesan":  "Gagal menyimpan berita",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"sukses": true,
		"data":   b,
	})
}

func UpdateBerita(c *gin.Context) {
	// ambil ID dari URL
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"sukses": false,
			"pesan":  "ID tidak valid",
		})
		return
	}

	// Bind body
	var input struct {
		Judul    string `json:"judul"`
		Konten   string `json:"konten"`
		Penulis  string `json:"penulis"`
		Kategori string `json:"kategori"`
		Gambar   string `json:"gambar"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"sukses": false,
			"pesan":  "Data tidak valid",
			"error":  err.Error(),
		})
		return
	}

	if input.Judul == "" || input.Konten == "" || input.Penulis == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"sukses": false,
			"pesan":  "Judul, konten, dan penulis wajib diisi",
		})
		return
	}

	slug := utils.BuatSlug(input.Judul)

	var kategori *string
	if input.Kategori != "" {
		kategori = &input.Kategori
	}

	var gambar *string
	if input.Gambar != "" {
		gambar = &input.Gambar
	}

	// update dan return data yang baru
	var b models.Berita
	err = config.DB.QueryRow(
		context.Background(),
		`UPDATE berita 
		 SET judul = $1, slug = $2, konten = $3, penulis = $4, kategori = $5, gambar = $6, updated_at = NOW()
		 WHERE id = $7
		 RETURNING id, judul, slug, konten, penulis, kategori, gambar, created_at, updated_at`,
		input.Judul,
		slug,
		input.Konten,
		input.Penulis,
		kategori,
		gambar,
		id,
	).Scan(
		&b.ID,
		&b.Judul,
		&b.Slug,
		&b.Konten,
		&b.Penulis,
		&b.Kategori,
		&b.Gambar,
		&b.CreatedAt,
		&b.UpdatedAt,
	)

	if err != nil {
		// cek apakah errornya karena baris tidak ditemukan
		if err.Error() == "no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{
				"sukses": false,
				"pesan":  "Berita tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"sukses": false,
			"pesan":  "Gagal memperbarui berita",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sukses": true,
		"data":   b,
	})
}

func DeleteBerita(c *gin.Context) {
	//ambil ID dari URL
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"sukses": false,
			"pesan":  "ID tidak valid",
		})
		return
	}

	//hapus dan return data yang dihapus
	var b models.Berita
	err = config.DB.QueryRow(
		context.Background(),
		`DELETE FROM berita WHERE id = $1
		RETURNING id, judul, slug, konten, penulis, kategori, gambar, created_at, updated_at`,
		id,
	).Scan(
		&b.ID,
		&b.Judul,
		&b.Slug,
		&b.Konten,
		&b.Penulis,
		&b.Kategori,
		&b.Gambar,
		&b.CreatedAt,
		&b.UpdatedAt,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{
				"sukses": false,
				"pesan":  "Berita tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"sukses": false,
			"pesan":  "Gagal menghapus berita",
			"error":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"sukses": true,
		"pesan":  "Berita berhasil dihapus",
		"data":   b,
	})
}

func GetBeritaBySlug(c *gin.Context) {
	//ambil slug dari URL
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"sukses": false,
			"pesan":  "Slug tidak valid",
		})
		return
	}

	//query
	var b models.Berita
	err := config.DB.QueryRow(
		context.Background(),
		`SELECT id, judul, slug, konten, penulis, kategori, gambar, created_at, updated_at
		FROM berita WHERE slug = $1`,
		slug,
	).Scan(
		&b.ID,
		&b.Judul,
		&b.Slug,
		&b.Konten,
		&b.Penulis,
		&b.Kategori,
		&b.Gambar,
		&b.CreatedAt,
		&b.UpdatedAt,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{
				"sukses": false,
				"pesan":  "Berita tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"sukses": false,
			"pesan":  "Gagal mengambil berita",
			"error":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"sukses": true,
		"data":   b,
	})
}
