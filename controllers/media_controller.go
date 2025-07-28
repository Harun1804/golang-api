package controllers

import (
	"galaxy/backend-api/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UploadHandler(c *gin.Context) {
    // Get file from form
    file, header, err := c.Request.FormFile("media")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Generate filename
    filename := helpers.GenerateFilename(header.Filename)

    // Save file (no resize for non-image files)
    err = helpers.SaveMedia(file, filename, false, "media")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Build URL
    url := helpers.GetMediaURL("media", filename)

    c.JSON(http.StatusOK, gin.H{
        "filename": filename,
        "url":      url,
    })
}

func DeleteHandler(c *gin.Context) {
    filename := c.Param("filename")
    err := helpers.DeleteMedia("media", filename)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}