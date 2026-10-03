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

func CreateBookmark(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"sukses": false, "pesan": "Belum login"})
		return
	}
	claims := userInterface.(*utils.JWTClaims)

	var input struct{
		BeritaID int `json:"beritaId"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"sukses": false, "pesan": "Input tidak valid", "error": err.Error()})
		return
	}
	if input.BeritaID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"sukses": false, "pesan": "Berita ID tidak valid"})
		return
	}

	var beritaID int
	err := config.DB.QueryRow(
		context.Background(),
		"SELECT id FROM berita WHERE id = $1",
		input.BeritaID,
	).Scan(&beritaID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"sukses": false, "pesan": "Berita tidak ditemukan"})
		return
	}

	var existingID int
	err = config.DB.QueryRow(
		context.Background(),
		"SELECT id FROM bookmark WHERE user_id = $1 AND berita_id = $2",
		claims.ID, input.BeritaID,
	).Scan(&existingID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"sukses": false, "pesan": "Berita sudah di-bookmark"})
		return
	}

	var b models.Bookmark
	err = config.DB.QueryRow(
		context.Background(),
		"INSERT INTO bookmark (user_id, berita_id) VALUES ($1, $2) RETURNING id, user_id, berita_id, created_at",
		claims.ID, input.BeritaID,
	).Scan(&b.ID, &b.UserID, &b.BeritaID, &b.CreatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"sukses": false, "pesan": "Gagal menyimpan bookmark", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"sukses": true, "pesan": "Berita berhasil di-bookmark", "data": b})
}

func GetBookmark(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"sukses": false, "pesan": "Belum login"})
		return
	}

	claims := userInterface.(*utils.JWTClaims)

	rows, err := config.DB.Query(
		context.Background(),
		`SELECT b.id, b.user_id, b.berita_id, b.created_at, 
				br.judul, br.slug, br.konten, br.penulis, br.kategori, br.gambar, br.created_at AS berita_created_at
		 FROM bookmark b
		 JOIN berita br ON b.berita_id = br.id
		 WHERE b.user_id = $1
		 ORDER BY b.created_at DESC`,
		claims.ID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"sukses": false, "pesan": "Gagal mengambil bookmark", "error": err.Error()})
		return
	}
	defer rows.Close()

	daftarBookmark := []models.BookmarkDenganBerita{}
	for rows.Next() {
		var b models.BookmarkDenganBerita
		err := rows.Scan(&b.ID, &b.UserID, &b.BeritaID, &b.CreatedAt, &b.Judul, &b.Slug, &b.Konten, &b.Penulis, &b.Kategori, &b.Gambar, &b.BeritaCreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"sukses": false, "pesan": "Gagal memproses data bookmark", "error": err.Error()})
			return
		}
		daftarBookmark = append(daftarBookmark, b)
	}

	c.JSON(http.StatusOK, gin.H{"sukses": true, "pesan": "Bookmark berhasil diambil", "data": daftarBookmark})
}

func DeleteBookmark(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"sukses": false, "pesan": "Belum login"})
		return
	}
	claims := userInterface.(*utils.JWTClaims)

	beritaIDStr := c.Param("berita_id")
	beritaID, err := strconv.Atoi(beritaIDStr)
	if err != nil || beritaID <= 0 {
		c.JSON(http.StatusBadRequest,gin.H{"sukses": false, "message": "Berita ID tidak valid"})
		return
	}
	
	result, err := config.DB.Exec(
		context.Background(),
		"DELETE FROM bookmark WHERE user_id = $1 AND berita_id = $2",
		claims.ID, beritaID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"sukses": false, "message": "Gagal menghapus bookmark", "error": err.Error()})
		return
	}
	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"sukses": false, "message": "Bookmark tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"sukses": true,
		"message": "Bookmark berhasil dihapus",
	})
}

func CekBookmark(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"sukses": false,
			"pesan": "Belum login",
		})
		return
	}
	claims := userInterface.(*utils.JWTClaims)

	beritaIDStr := c.Param("berita_id")
	beritaID, err := strconv.Atoi(beritaIDStr)
	if err != nil || beritaID <= 0{
		c.JSON(http.StatusBadRequest, gin.H{
			"sukses": false,
			"pesan": "Berita ID tidak valid",
		})
		return
	}
	var existingID int
	err = config.DB.QueryRow(
		context.Background(),
		"SELECT id FROM bookmark WHERE user_id = $1 AND berita_id = $2",
		claims.ID,
		beritaID,
	).Scan(&existingID)

	if err == nil{
		c.JSON(http.StatusOK, gin.H{
			"sukses": true, "tersimpan":true,
		})
		return
	}
	if err.Error() == "no rows in result set"{
		c.JSON(http.StatusOK, gin.H{
			"sukses": true, "tersimpan": false,
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"sukses": false, "pesan": "Gagal mengecek bookmark", "error": err.Error(),
	}) 
}