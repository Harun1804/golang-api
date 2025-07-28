package controllers

import (
	"galaxy/backend-api/database"
	"galaxy/backend-api/helpers"
	"galaxy/backend-api/models"
	"galaxy/backend-api/payloads"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Register menangani proses registrasi user baru
func Register(c *gin.Context) {
	// Inisialisasi struct untuk menangkap data request
	var req = payloads.UserCreateRequest{}

	// Validasi request JSON menggunakan binding dari Gin
	if err := c.ShouldBindJSON(&req); err != nil {
		// Jika validasi gagal, kirimkan response error
		helpers.SendError(c, http.StatusUnprocessableEntity, "Validation Errors", err)
		return
	}

	// Buat data user baru dengan password yang sudah di-hash
	user := models.User{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: helpers.HashPassword(req.Password),
	}

	// Simpan data user ke database
	if err := database.DB.Create(&user).Error; err != nil {
		// Cek apakah error karena data duplikat (misalnya username/email sudah terdaftar)
		if helpers.IsDuplicateEntryError(err) {
			// Jika duplikat, kirimkan response 409 Conflict
			helpers.SendError(c, http.StatusConflict, "Duplicate entry error", err)
		} else {
			// Jika error lain, kirimkan response 500 Internal Server Error
			helpers.SendError(c, http.StatusInternalServerError, "Failed to create user", err)
		}
		return
	}

	helpers.SendSuccess(c, http.StatusCreated, "User created successfully", payloads.ToUserResponse(user, ""))
}

func Login(c *gin.Context) {
// Inisialisasi struct untuk menampung data dari request
	var req = payloads.UserLoginRequest{}
	var user = models.User{}

	// Validasi input dari request body menggunakan ShouldBindJSON
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.SendError(c, http.StatusUnprocessableEntity, "Validation Errors", err)
		return
	}

	// Cari user berdasarkan username yang diberikan di database
	// Jika tidak ditemukan, kirimkan respons error Unauthorized
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		helpers.SendError(c, http.StatusUnauthorized, "User Not Found", err)
		return
	}

	// Bandingkan password yang dimasukkan dengan password yang sudah di-hash di database
	// Jika tidak cocok, kirimkan respons error Unauthorized
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		helpers.SendError(c, http.StatusUnauthorized, "Invalid Password", err)
		return
	}

	// Jika login berhasil, generate token untuk user
	token := helpers.GenerateToken(user.Username)

	helpers.SendSuccess(c, http.StatusOK, "Login Success", payloads.ToUserResponse(user, token))
}
