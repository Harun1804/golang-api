package controllers

import (
	"galaxy/backend-api/database"
	"galaxy/backend-api/helpers"
	"galaxy/backend-api/models"
	"galaxy/backend-api/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUsers(ctx *gin.Context) {
	// Inisialisasi slice untuk menampung data user
	var users []models.User

	// Ambil semua data user dari database
	if err := database.DB.Find(&users).Error; err != nil {
		// Jika terjadi error, kirimkan response error
		ctx.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to retrieve users",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	// Kirimkan response sukses dengan data user
	ctx.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Users retrieved successfully",
		Data:    users,
	})
}

func CreateUser(ctx *gin.Context) {
	var req = structs.UserCreateRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Invalid request data",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	user := models.User{
		Name: 	 req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: helpers.HashPassword(req.Password),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create user",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	ctx.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "User created successfully",
		Data:    user,
	})
}
