package helpers

import (
	"errors"
	"galaxy/backend-api/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Nilai secret diambil dari environment variable JWT_SECRET
var jwtKey = []byte(config.GetEnv("JWT_SECRET", "secret"))

func GenerateToken(username string) string {
	// Mengatur waktu kedaluwarsa token, di sini kita set 60 menit dari waktu sekarang
	expirationTime := time.Now().Add(60 * time.Minute)

	// Membuat klaim (claims) JWT
	// Subject berisi username, dan ExpiresAt menentukan waktu expired token
	claims := &jwt.RegisteredClaims{
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	// Membuat token baru dengan klaim yang telah dibuat
	// Menggunakan algoritma HS256 untuk menandatangani token
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtKey)

	// Mengembalikan token dalam bentuk string
	return token
}

func AltGenerateToken(userId uint, username, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":  userId,
		"username": username,
		"email":    email,
		"exp":      time.Now().Add(60 * time.Minute).Unix(),
	})

	return token.SignedString(jwtKey)
}

func VerifyToken(token string) (uint, error) {
	parsedToken, err := parseToken(token)
	if err := validateToken(parsedToken, err); err != nil {
		return 0, err
	}
	// If you want to use claims:
	claims, err := extractClaims(parsedToken)
	if err != nil {
		return 0, err
	}
	userIdVal, ok := claims["userId"]
	if !ok || userIdVal == nil {
		return 0, errors.New("userId claim missing in token")
	}
	userId, ok := userIdVal.(float64)
	if !ok {
		return 0, errors.New("userId claim is not a valid number")
	}
	return uint(userId), nil
}

func parseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtKey, nil
	})
}

func validateToken(parsedToken *jwt.Token, err error) error {
	if err != nil {
		return errors.New("Could not parse token")
	}
	if !parsedToken.Valid {
		return errors.New("Invalid token")
	}
	return nil
}

func extractClaims(parsedToken *jwt.Token) (jwt.MapClaims, error) {
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("Invalid token claims")
	}
	return claims, nil
}
