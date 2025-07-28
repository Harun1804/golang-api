package controllers

import (
	"galaxy/backend-api/database"
	"galaxy/backend-api/helpers"
	"galaxy/backend-api/models"
	"galaxy/backend-api/payloads"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUsers(ctx *gin.Context) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		helpers.SendError(ctx, http.StatusInternalServerError, "Failed to retrieve users", err)
		return
	}
	var userResponses []payloads.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, payloads.ToUserResponse(user, ""))
	}
	helpers.SendSuccess(ctx, http.StatusOK, "Users retrieved successfully", userResponses)
}

func GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		helpers.SendError(ctx, http.StatusNotFound, "User not found", err)
		return
	}
	helpers.SendSuccess(ctx, http.StatusOK, "User retrieved successfully", payloads.ToUserResponse(user, ""))
}

func CreateUser(ctx *gin.Context) {
	var req payloads.UserCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.SendError(ctx, http.StatusUnprocessableEntity, "Validation error", err)
		return
	}
	user := models.User{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: helpers.HashPassword(req.Password),
	}
	if err := database.DB.Create(&user).Error; err != nil {
		helpers.SendError(ctx, http.StatusInternalServerError, "Failed to create user", err)
		return
	}
	helpers.SendSuccess(ctx, http.StatusCreated, "User created successfully", payloads.ToUserResponse(user, ""))
}

func UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		helpers.SendError(ctx, http.StatusNotFound, "User not found", err)
		return
	}
	var req payloads.UserUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.SendError(ctx, http.StatusUnprocessableEntity, "Validation error", err)
		return
	}
	user.Name = req.Name
	user.Username = req.Username
	user.Email = req.Email
	if req.Password != "" {
		user.Password = helpers.HashPassword(req.Password)
	}
	if err := database.DB.Save(&user).Error; err != nil {
		helpers.SendError(ctx, http.StatusInternalServerError, "Failed to update user", err)
		return
	}
	helpers.SendSuccess(ctx, http.StatusOK, "User updated successfully", payloads.ToUserResponse(user, ""))
}

func DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		helpers.SendError(ctx, http.StatusNotFound, "User not found", err)
		return
	}
	if err := database.DB.Delete(&user).Error; err != nil {
		helpers.SendError(ctx, http.StatusInternalServerError, "Failed to delete user", err)
		return
	}
	helpers.SendSuccess(ctx, http.StatusOK, "User deleted successfully", nil)
}