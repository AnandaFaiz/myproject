package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"uuid"

	"github.com/gin-gonic/gin"
)

var tipeValid = map[string]string{
	"image/jpeg": "jpg",
	"image/png":  "png",
	"image/webp": "webp",
	"image/gif":  "gif",
}

// max 5MB
const maxFileSize = 5 * 1024 * 1024

func UploadGambar(c *gin.Context) {
	// Ambil file dari form-data dengan key "file"
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"sukses": false,
			"pesan":  "File tidak ditemukan (gunakan key 'file')",
			"error":  err.Error(),
		})
		return
	}

	// Cek ukuran file
	if file.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"sukses": false,
			"pesan":  "Ukuran file maksimal 5 MB",
		})
		return
	}

	// Buka file untuk baca tipe
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"sukses": false,
			"pesan":  "Gagal membuka file",
		})
		return
	}
	defer f.Close()

	// Baca 512 byte pertama untuk deteksi MIME type
	buffer := make([]byte, 512)
	_, err = f.Read(buffer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"sukses": false,
			"pesan":  "Gagal membaca file",
		})
		return
	}
	mimeType := http.DetectContentType(buffer)

	// Cek tipe file
	ext, ok := tipeValid[mimeType]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"sukses": false,
			"pesan":  "Format file harus JPG, PNG, WEBP, atau GIF",
		})
		return
	}

	// Buat folder kalau belum ada
	folderUpload := "public/uploads"
	if err := os.MkdirAll(folderUpload, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"sukses": false,
			"pesan":  "Gagal membuat folder upload",
		})
		return
	}

	// Bikin nama file unik
	// uuid v4 + timestamp (jaga-jaga kalau uuid tabrakan)
	namaFile := fmt.Sprintf("%s_%d.%s",
		strings.ReplaceAll(uuid.New().String(), "-", ""),
		time.Now().Unix(),
		ext,
	)

	// Simpan file
	pathFile := filepath.Join(folderUpload, namaFile)
	if err := c.SaveUploadedFile(file, pathFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"sukses": false,
			"pesan":  "Gagal menyimpan file",
			"error":  err.Error(),
		})
		return
	}

	// Return URL publik
	c.JSON(http.StatusOK, gin.H{
		"sukses": true,
		"url":    "/uploads/" + namaFile,
	})
}
