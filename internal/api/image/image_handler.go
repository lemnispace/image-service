package api

import (
	"fmt"
	"image-service/internal/repository"
	"image-service/internal/service"
	u "image-service/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ImageHandler struct {
	repo    repository.ImageRepository
	service service.ImageService
}

func NewImageHandler(repo repository.ImageRepository) (*ImageHandler, error) {
	service, error := service.NewImageService()
	if error != nil {
		return nil, fmt.Errorf("failed to create image service: %w", error)
	}
	return &ImageHandler{repo: repo, service: service}, nil
}

func (h *ImageHandler) GetImages(c *gin.Context) {
	// Implement pagination and fetching logic
	images, err := h.repo.GetImages(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get images: %v", err)})
	}
	c.JSON(http.StatusOK, images)
}

func (h *ImageHandler) GenerateImage(c *gin.Context) {
	// Parse the request
	var req u.APIGenTextToImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to parse request: %v", err)})
	}
	// generate the image and store it in the database
	image, err := h.service.GenerateImage(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate image: %v", err)})
	}
	var imgsData [][]byte
	for _, img := range image.Artifacts {
		imgbd, err := u.ConvertBase64ToImageBytes(img.Base64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate image: %v", err)})
		}
		imgsData = append(imgsData, imgbd)
	}
	// store the image in the database
	err = h.storeImage(c, req, image, imgsData[0])
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to store image: %v", err)})
	}

	c.JSON(http.StatusOK, imgsData[0])
}

func (h *ImageHandler) storeImage(c *gin.Context, genImgReq u.APIGenTextToImageRequest, imgRes *u.APIGenTextToImageResponse, imgData []byte) error {
	img := u.ImageMetadata{
		Seed:         imgRes.Artifacts[0].Seed,
		Prompt:       genImgReq.Prompt,
		FinishReason: imgRes.Artifacts[0].FinishReason,
		CreatedAt:    u.GetCurrentTime(),
		CreatedBy:    "admin",
		Tags:         []string{},
		Categories:   []string{},
		Styles:       []string{},
	}
	return h.repo.StoreImage(c, img, imgData)
}
