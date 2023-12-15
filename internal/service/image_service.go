package service

import (
	"context"
	u "image-service/internal/util"
)

type ImageService interface {
	GenerateImage(ctx context.Context) (*u.APIGenTextToImageResponse, error)
}

type imageService struct {
}

func NewImageService() (ImageService, error) {
	return &imageService{}, nil
}

func (s *imageService) GenerateImage(ctx context.Context) (*u.APIGenTextToImageResponse, error) {
	// Call another Lambda function to generate an image and return the result
	return nil, nil
}
