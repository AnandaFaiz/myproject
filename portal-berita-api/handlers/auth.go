package handlers

import (
	"context"
	"net/http"
	"portal-berita-api/config"
	"portal-berita-api/models"
	"portal-berita-api/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var input struct {
		Nama     string `json:"nama"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"sukses": false,
			"pesan":  "Data tidak valid",
			"error":  err.Error(),
		})
		return
	}

	//trim space
	input.Nama = strings.TrimSpace(input.Nama)
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))

	//validasi
	if input.Nama == "" || input.Email == "" || input.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"sukses": false,
			"pesan":  "Nama, email, dan password wajib diisi",
		})
		return
	}

	if len(input.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{
			"sukses": false,
			"pesan":  "Password minimal 6 karakter",
		})
		return
	}

	//cek email
	var existingID int
	err := config.DB.QueryRow(
		context.Background(),
		"SELECT id FROM users WHERE email = $1",
		input.Email,
	).Scan(&existingID)

	if err == nil {
		//tidak error => baris ditemukan
		c.JSON(http.StatusBadRequest, gin.H{
			"sukses": false,
			"pesan":  "Email sudah terdaftar",
		})
		return
	}

	//hash pw
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"sukses": false,
			"pesan":  "Gagal hash password",
			"error":  err.Error(),
		})
		return
	}

	//insert error
	var u models.User
	err = config.DB.QueryRow(
		context.Background(),
		`INSERT INTO users (nama, email, password, role) VALUES ($1, $2, $3, $4)
		RETURNING id, nama, email, role, created_at`,
		input.Nama,
		input.Email,
		string(hashedPassword),
		"user",
	).Scan(
		&u.ID,
		&u.Nama,
		&u.Email,
		&u.Role,
		&u.CreatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"sukses": false,
			"pesan":  "Gagal mendaftar",
			"error":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"sukses": true,
		"pesan":  "Registrasi berhasil",
		"user":   u,
	})
}

func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"sukses": false,
			"pesan":  "Data tidak valid",
			"error":  err.Error(),
		})
		return
	}

	// Trim spasi & lowercase email
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))

	if input.Email == "" || input.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"sukses": false,
			"pesan":  "Email dan password wajib diisi",
		})
		return
	}

	// Cari user berdasarkan email
	var u models.User
	err := config.DB.QueryRow(
		context.Background(),
		`SELECT id, nama, email, password, role, created_at
		 FROM users
		 WHERE email = $1`,
		input.Email,
	).Scan(
		&u.ID,
		&u.Nama,
		&u.Email,
		&u.Password,
		&u.Role,
		&u.CreatedAt,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"sukses": false,
				"pesan":  "Email atau password salah",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"sukses": false,
			"pesan":  "Gagal login",
			"error":  err.Error(),
		})
		return
	}

	// Cek password
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"sukses": false,
			"pesan":  "Email atau password salah",
		})
		return
	}

	// Bikin JWT token
	token, err := utils.BuatToken(u.ID, u.Nama, u.Email, u.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"sukses": false,
			"pesan":  "Gagal membuat token",
			"error":  err.Error(),
		})
		return
	}

	// Clear password sebelum kirim ke frontend
	u.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"sukses": true,
		"pesan":  "Login berhasil",
		"token":  token,
		"user":   u,
	})
}

func Me(c *gin.Context) {
	// Ambil user dari context (di-set oleh AuthMiddleware)
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"sukses": false,
			"pesan":  "Belum login",
		})
		return
	}

	claims, ok := userInterface.(*utils.JWTClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"sukses": false,
			"pesan":  "Gagal membaca data user",
		})
		return
	}

	//ambil dari DB
	var u models.User
	err := config.DB.QueryRow(
		context.Background(),
		`SELECT id, nama, email, role, created_at FROM users WHERE id = $1`,
		claims.ID,
	).Scan(
		&u.ID,
		&u.Nama,
		&u.Email,
		&u.Role,
		&u.CreatedAt,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"sukses": false,
			"pesan":  "User tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sukses": true,
		"user":   u,
	})
}

func UpdateProfil(c *gin.Context) {
	// Ambil user dari context (di-set oleh AuthMiddleware)
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"sukses": false,
			"pesan":  "Belum login",
		})
		return
	}
	claims := userInterface.(*utils.JWTClaims)

	var input struct {
		Nama  string `json:"nama"`
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"sukses": false,
			"pesan":  "Data tidak valid",
			"error":  err.Error(),
		})
		return
	}

	input.Nama = strings.TrimSpace(input.Nama)
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))

	if input.Nama == "" || input.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"sukses": false,
			"pesan":  "Nama dan email wajib diisi",
		})
		return
	}

	// Cek apakah email sudah dipakai user lain (kecuali dirinya sendiri)
	var existingID int
	err := config.DB.QueryRow(
		context.Background(),
		"SELECT id FROM users WHERE email = $1 AND id != $2",
		input.Email,
		claims.ID,
	).Scan(&existingID)

	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"sukses": false,
			"pesan":  "Email sudah dipakai user lain",
		})
		return
	}

	// Update
	var u models.User
	err = config.DB.QueryRow(
		context.Background(),
		`UPDATE users SET nama = $1, email = $2 WHERE id = $3
		 RETURNING id, nama, email, role, created_at`,
		input.Nama,
		input.Email,
		claims.ID,
	).Scan(
		&u.ID,
		&u.Nama,
		&u.Email,
		&u.Role,
		&u.CreatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"sukses": false,
			"pesan":  "Gagal memperbarui profil",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sukses": true,
		"pesan":  "Profil berhasil diperbarui",
		"user":   u,
	})
}

func GantiPassword(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"sukses": false,
			"pesan":  "Belum login",
		})
		return
	}
	claims := userInterface.(*utils.JWTClaims)

	var input struct {
		PasswordLama string `json:"passwordLama"`
		PasswordBaru string `json:"passwordBaru"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"sukses": false,
			"pesan":  "Data tidak valid",
			"error":  err.Error(),
		})
		return
	}

	if input.PasswordLama == "" || input.PasswordBaru == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"sukses": false,
			"pesan":  "Semua field wajib diisi",
		})
		return
	}

	if len(input.PasswordBaru) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{
			"sukses": false,
			"pesan":  "Password baru minimal 6 karakter",
		})
		return
	}

	// Ambil user dari database (untuk cek password lama)
	var passwordHash string
	err := config.DB.QueryRow(
		context.Background(),
		"SELECT password FROM users WHERE id = $1",
		claims.ID,
	).Scan(&passwordHash)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"sukses": false,
			"pesan":  "User tidak ditemukan",
		})
		return
	}

	// Cek password lama
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(input.PasswordLama))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"sukses": false,
			"pesan":  "Password lama salah",
		})
		return
	}

	// Hash password baru
	hashedBaru, err := bcrypt.GenerateFromPassword([]byte(input.PasswordBaru), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"sukses": false,
			"pesan":  "Gagal hash password baru",
			"error":  err.Error(),
		})
		return
	}

	// Update
	_, err = config.DB.Exec(
		context.Background(),
		"UPDATE users SET password = $1 WHERE id = $2",
		string(hashedBaru),
		claims.ID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"sukses": false,
			"pesan":  "Gagal mengganti password",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sukses": true,
		"pesan":  "Password berhasil diganti",
	})
}

func Logout(c *gin.Context) {
	// JWT tidak bisa dihapus dari server.
	// Cukup beri tahu frontend untuk hapus token dari storage.
	c.JSON(http.StatusOK, gin.H{
		"sukses": true,
		"pesan":  "Logout berhasil",
	})
}

