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
