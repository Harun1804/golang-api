package middlewares

import (
	"galaxy/backend-api/config"
	"galaxy/backend-api/helpers"
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
				"error": "Invalid token " + err.Error(),
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

func AltAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userId uint
		var err error

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			helpers.MiddlewareError(c, http.StatusUnauthorized, "Authorization header is required", nil)
			return
		}

		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			userId, err = helpers.VerifyToken(tokenString)
			if err != nil {
				helpers.MiddlewareError(c, http.StatusUnauthorized, "Invalid token", err)
				return
			}
		} else {
			helpers.MiddlewareError(c, http.StatusUnauthorized, "Authorization header must start with Bearer", nil)
			return
		}

		c.Set("userId", userId)
		c.Next()
	}
}
