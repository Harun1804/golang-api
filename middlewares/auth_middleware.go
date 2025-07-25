package middlewares

import (
	"galaxy/backend-api/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(config.GetEnv("JWT_SECRET", "secret"))

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil header Authorization dari request
		tokenString := c.GetHeader("Authorization")

		// Jika token kosong, kembalikan respons 401 Unauthorized
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token is required",
			})
			c.Abort() // Hentikan request selanjutnya
			return
		}

		// Hapus prefix "Bearer " dari token
		// Header biasanya berbentuk: "Bearer <token>"
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Buat struct untuk menampung klaim token
		claims := &jwt.RegisteredClaims{}

		// Parse token dan verifikasi tanda tangan dengan jwtKey
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Kembalikan kunci rahasia untuk memverifikasi token
			return jwtKey, nil
		})

		// Jika token tidak valid atau terjadi error saat parsing
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			c.Abort() // Hentikan request
			return
		}

		// Simpan klaim "sub" (username) ke dalam context
		c.Set("username", claims.Subject)

		// Lanjut ke handler berikutnya
		c.Next()
	}
}
