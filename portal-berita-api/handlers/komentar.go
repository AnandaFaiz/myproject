package handlers

import (
	"context"
	"log"
	"net/http"
	"portal-berita-api/config"
	"portal-berita-api/models"
	"portal-berita-api/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func CreateKomentar(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"sukses": false,
			"pesan": "Belum login",
		})
		return
	}
	claims := userInterface.(*utils.JWTClaims)

	//bind body
	var input struct {
		BeritaID int    `json:"beritaId"`
		Isi      string `json:"isi"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"sukses": false,
			"pesan":  "Data tidak valid",
		})
		return
	}

	log.Println("Input diterima:", input)
	log.Println("Isi:", input.Isi)
	log.Println("BeritaID:", input.BeritaID)

	input.Isi = strings.TrimSpace(input.Isi)

	if input.Isi == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"sukses": false,
			"pesan":  "Komentar tidak boleh kosong",
		})
		return
	}

	if len(input.Isi) > 1000 {
		c.JSON(http.StatusBadRequest, gin.H{
			"sukses": false,
			"pesan":  "Komentar terlalu panjang",
		})
		return
	}

	//Cek apakah berita ada/ tidak
	var beritaID int
	err := config.DB.QueryRow(context.Background(), "SELECT id FROM berita WHERE id = $1", input.BeritaID).Scan(&beritaID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"sukses": false,
			"pesan":  "Berita tidak ditemukan",
		})
		return
	}

	//Insert komentar + return data dengan user_nama
	var k models.Komentar
	err = config.DB.QueryRow(context.Background(), "INSERT INTO komentar (berita_id, user_id, isi) VALUES ($1, $2, $3) RETURNING id, berita_id, user_id, isi, created_at, updated_at", input.BeritaID, claims.ID, input.Isi).Scan(&k.ID, &k.BeritaID, &k.UserID, &k.Isi, &k.CreatedAt, &k.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"sukses": false,
			"pesan":  "Gagal membuat komentar",
			"error":  err.Error(),
		})
		return
	}

	//Ambil nama user dari claims
	k.UserNama = claims.Nama

	c.JSON(http.StatusCreated, gin.H{
		"sukses": true,
		"pesan":  "Komentar berhasil ditambahkan",
		"data":   k,
	})
}

func GetKomentarByBerita(c *gin.Context) {
	// Ambil berita_id dari URL
	beritaIDStr := c.Param("berita_id")
	beritaID, err := strconv.Atoi(beritaIDStr)
	if err != nil || beritaID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"sukses": false,
			"pesan":  "ID berita tidak valid",
		})
		return
	}

	// Query komentar berdasarkan berita_id
	rows, err := config.DB.Query(context.Background(), `SELECT k.id, k.berita_id, k.user_id, k.isi, k.created_at, k.updated_at, u.nama AS user_nama FROM komentar k JOIN users u ON u.id = k.user_id WHERE k.berita_id = $1 ORDER BY k.created_at DESC`, beritaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"sukses": false,
			"pesan":  "Gagal mengambil komentar",
			"error":  err.Error(),
		})
		return
	}
	defer rows.Close()

	daftarKomentar := []models.Komentar{}
	for rows.Next() {
		var k models.Komentar
		err := rows.Scan(&k.ID, &k.BeritaID, &k.UserID, &k.Isi, &k.CreatedAt, &k.UpdatedAt, &k.UserNama)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"sukses": false,
				"pesan":  "Gagal memproses data komentar",
				"error":  err.Error(),
			})
			return
		}
		daftarKomentar = append(daftarKomentar, k)
	}

	c.JSON(http.StatusOK, gin.H{
		"sukses": true,
		"pesan":  "Komentar berhasil diambil",
		"data":   daftarKomentar,
	})
}

func DeleteKomentar(c *gin.Context) {
	//Ambil user dari context (di-set oleh AuthMiddleware)
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"sukses": false,
			"pesan":  "Belum login",
		})
		return
	}
	claims := userInterface.(*utils.JWTClaims)

	//ambil ID komentar dari URL
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"sukses": false,
			"pesan":  "ID komentar tidak valid",
		})
		return
	}

	//ambil data komentar dulu (utk cek kepemilikan)
	var komentarUserID int
	err = config.DB.QueryRow(context.Background(), "SELECT user_id FROM komentar WHERE id = $1", id).Scan(&komentarUserID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{
				"sukses": false,
				"pesan":  "Komentar tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"sukses": false,
			"pesan":  "Gagal membaca data komentar",
			"error":  err.Error(),
		})
		return
	}

	//cek kepemilikan komentar
	adalahAdmin := claims.Role == "admin"
	adalahPemilik := komentarUserID == claims.ID

	if !adalahAdmin && !adalahPemilik {
		c.JSON(http.StatusForbidden, gin.H{
			"sukses": false,
			"pesan":  "Anda tidak memiliki izin untuk menghapus komentar ini",
		})
		return	
	}

	//Hapus komentar
	result, err := config.DB.Exec(context.Background(), "DELETE FROM komentar WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"sukses": false,
			"pesan":  "Gagal menghapus komentar",
			"error":  err.Error(),
		})
		return
	}

	//cek apakah ada baris yg dihapus
	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"sukses": false,
			"pesan":  "Komentar tidak ditemukan",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"sukses": true,
		"pesan":  "Komentar berhasil dihapus",
	})
}