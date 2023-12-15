package repository

import (
	"context"
	m "image-service/internal/model"
	u "image-service/internal/util"
)

type ImageRepository interface {
	// GetImage(ctx context.Context) (m.Image, error)
	GetImages(ctx context.Context) ([]m.Image, error)
	StoreImage(ctx context.Context, image u.ImageMetadata, imgData []byte) error
	// StoreImages(ctx context.Context, image []u.ImageMetadata, imgData [][]byte) error
}
