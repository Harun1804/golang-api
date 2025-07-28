package controllers

import (
	"galaxy/backend-api/database"
	"galaxy/backend-api/helpers"
	"galaxy/backend-api/models"
	"galaxy/backend-api/payloads"
	"net/http"

	"github.com/gin-gonic/gin"
)

var filePath = "media/posts"

func GetPosts(ctx *gin.Context) {
	var posts []models.Post
	if err := database.DB.Find(&posts).Error; err != nil {
		helpers.SendError(ctx, http.StatusInternalServerError, "Failed to retrieve posts", err)
		return
	}

	var postResponses []payloads.PostResponse
	for _, post := range posts {
		postResponses = append(postResponses, payloads.ToPostResponse(post))
	}

	helpers.SendSuccess(ctx, http.StatusOK, "Posts retrieved successfully", postResponses)
}

func GetPost(ctx *gin.Context) {
	id := ctx.Param("id")
	var post models.Post

	if err := database.DB.First(&post, id).Error; err != nil {
		helpers.SendError(ctx, http.StatusNotFound, "Post not found", err)
		return
	}

	response := payloads.ToPostResponse(post)
	helpers.SendSuccess(ctx, http.StatusOK, "Post retrieved successfully", response)
}

func CreatePost(ctx *gin.Context) {
	var req payloads.PostCreateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		helpers.SendError(ctx, http.StatusUnprocessableEntity, "Validation error", err)
		return
	}

	imageName := helpers.GenerateFilename(req.Image.Filename)
	imageFile, err := req.Image.Open()
	if err != nil {
			helpers.SendError(ctx, http.StatusBadRequest, "Failed to open image", err)
			return
	}
	defer imageFile.Close()

	err = helpers.SaveMedia(imageFile, imageName, true, filePath)
	if err != nil {
		helpers.SendError(ctx, http.StatusInternalServerError, "Failed to save image", err)
		return
	}

	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		Image:  imageName,
	}

	if err := database.DB.Create(&post).Error; err != nil {
		helpers.SendError(ctx, http.StatusInternalServerError, "Failed to create post", err)
		return
	}

	helpers.SendSuccess(ctx, http.StatusCreated, "Post created successfully", payloads.ToPostResponse(post))
}

func UpdatePost(ctx *gin.Context) {
	id := ctx.Param("id")
	var post models.Post

	if err := database.DB.First(&post, id).Error; err != nil {
		helpers.SendError(ctx, http.StatusNotFound, "Post not found", err)
		return
	}

	var req payloads.PostUpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		helpers.SendError(ctx, http.StatusUnprocessableEntity, "Validation error", err)
		return
	}

	post.Title = req.Title
	post.Content = req.Content

	if req.Image != nil {
		if err := helpers.DeleteMedia(filePath, post.Image); err != nil {
			helpers.SendError(ctx, http.StatusInternalServerError, "Failed to delete post image", err)
			return
		}

		imageName := helpers.GenerateFilename(req.Image.Filename)
		imageFile, err := req.Image.Open()
		if err != nil {
			helpers.SendError(ctx, http.StatusBadRequest, "Failed to open image", err)
			return
		}
		defer imageFile.Close()

		err = helpers.SaveMedia(imageFile, imageName, true, filePath)
		if err != nil {
			helpers.SendError(ctx, http.StatusInternalServerError, "Failed to save image", err)
			return
		}

		post.Image = imageName
	}

	if err := database.DB.Save(&post).Error; err != nil {
		helpers.SendError(ctx, http.StatusInternalServerError, "Failed to update post", err)
		return
	}

	response := payloads.ToPostResponse(post)
	helpers.SendSuccess(ctx, http.StatusOK, "Post updated successfully", response)
}

func DeletePost(ctx *gin.Context) {
	id := ctx.Param("id")
	var post models.Post

	if err := database.DB.First(&post, id).Error; err != nil {
		helpers.SendError(ctx, http.StatusNotFound, "Post not found", err)
		return
	}

	if err := helpers.DeleteMedia(filePath, post.Image); err != nil {
		helpers.SendError(ctx, http.StatusInternalServerError, "Failed to delete post image", err)
		return
	}

	if err := database.DB.Delete(&post).Error; err != nil {
		helpers.SendError(ctx, http.StatusInternalServerError, "Failed to delete post", err)
		return
	}

	helpers.SendSuccess(ctx, http.StatusOK, "Post deleted successfully", nil)
}