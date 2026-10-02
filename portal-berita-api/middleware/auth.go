package middleware

import (
	"net/http"
	"portal-berita-api/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//ambil token dari header Auth: Bearer
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"sukses": false,
				"pesan":  "Token tidak ditemukan",
			})
			c.Abort()
			return
		}

		//format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"sukses": false,
				"pesan":  "Format token tidak valid",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		//verify
		claims, err := utils.VerifikasiToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"sukses": false,
				"pesan":  "Token tidak valid",
			})
			c.Abort()
			return
		}

		//simpan claims
		c.Set("user", claims)
		c.Next()
	}
}

// cek user role = admin?
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"sukses": false,
				"pesan":  "Belum login",
			})
			c.Abort()
			return
		}

		claims, ok := userInterface.(*utils.JWTClaims)
		if !ok || claims.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"sukses": false,
				"pesan":  "Akses hanya untuk admin",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
