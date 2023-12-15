package api

import (
	"fmt"
	api "image-service/internal/api/image"
	"image-service/internal/repository"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, repo repository.ImageRepository) error {
	imageHandler, err := api.NewImageHandler(repo)
	if err != nil {
		return fmt.Errorf("failed to create image handler: %w", err)
	}
	r.GET("/", imageHandler.GetImages)
	r.POST("/generate", imageHandler.GenerateImage)
	return nil
}
