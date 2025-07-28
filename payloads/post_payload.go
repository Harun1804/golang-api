package payloads

import (
	"galaxy/backend-api/helpers"
	"galaxy/backend-api/models"
	"mime/multipart"
)

type PostResponse struct {
	Id        uint   `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Image     string `json:"image"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type PostCreateRequest struct {
	Title   string `form:"title" binding:"required"`
	Content string `form:"content" binding:"required"`
	Image   *multipart.FileHeader `form:"image"`
}

type PostUpdateRequest struct {
	Title   string `form:"title" binding:"required"`
	Content string `form:"content" binding:"required"`
	Image   *multipart.FileHeader `form:"image,omitempty"`
}

func ToPostResponse(post models.Post) PostResponse {
	imageUrl := helpers.GetMediaURL("media/posts", post.Image)
	return PostResponse{
		Id:        post.Id,
		Title:     post.Title,
		Content:   post.Content,
		Image:     imageUrl,
		CreatedAt: post.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: post.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}