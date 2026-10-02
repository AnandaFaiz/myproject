package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"portal-berita-api/config"
	"portal-berita-api/handlers"
	"portal-berita-api/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// load .env
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: file .env tidak ditemukan")
	}

	//koneksi DB
	config.ConnectDB()

	// setup router Gin
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-type", "Authorization"},
		ExposeHeaders:    []string{"Content-length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	//static uploads
	r.Static("/uploads", "./public/uploads")

	//endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"pesan": "Hi from backend Go",
		})
	})

	//Route API
	api := r.Group("/api")
	{
		//route ADMIN
		admin := api.Group("")
		admin.Use(middleware.AuthMiddleware())
		admin.Use(middleware.AdminMiddleware())
		{
			admin.POST("/berita", handlers.CreateBerita)
			admin.PUT("/berita/:id", handlers.UpdateBerita)
			admin.DELETE("/berita/:id", handlers.DeleteBerita)
			admin.POST("/upload", handlers.UploadGambar)
		}

		//route PUBLIK
		api.GET("/berita", handlers.GetBerita)
		api.GET("/berita/slug/:slug", handlers.GetBeritaBySlug)
		api.GET("/komentar/:berita_id", handlers.GetKomentarByBerita)

		//route auth
		api.POST("/auth/register", handlers.Register)
		api.POST("/auth/login", handlers.Login)
		api.POST("/auth/logout", handlers.Logout)

		user := api.Group("")
		user.Use(middleware.AuthMiddleware())
		{
			user.GET("/auth/me", handlers.Me)
			user.PUT("/auth/profil", handlers.UpdateProfil)
			user.PUT("/auth/password", handlers.GantiPassword)
			user.POST("/komentar", handlers.CreateKomentar)
			user.DELETE("/komentar/:id", handlers.DeleteKomentar)
		}
	}

	// //endpoint health check
	// r.GET("/health", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"status": "OK",
	// 	})
	// })

	//jalankan server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server jalan di http://localhost:%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Gagal menjalankan server:", err)
	}
}
