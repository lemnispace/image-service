package db

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	m "image-service/internal/model"
	r "image-service/internal/repository"
	u "image-service/internal/util"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	_ "github.com/lib/pq"
)

type ImageRepositoryImpl struct {
	db           *sql.DB
	s3Uploader   *s3manager.Uploader
	s3Downloader *s3manager.Downloader
}

func NewImageRepository() (r.ImageRepository, error) {
	// Connect to the PostgreSQL database
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_CONNECTION_STRING"))
	if err != nil {
		return nil, err
	}

	// Initialize AWS S3 session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		return nil, err
	}
	s3Uploader := s3manager.NewUploader(sess)
	s3Downloader := s3manager.NewDownloader(sess)

	return &ImageRepositoryImpl{
		db:           db,
		s3Uploader:   s3Uploader,
		s3Downloader: s3Downloader,
	}, nil
}

func (repo *ImageRepositoryImpl) GetImages(ctx context.Context) ([]m.Image, error) {
	//get all image metadata from postgres and img data from s3, then combine them
	// First, get the image metadata from PostgreSQL
	rows, err := repo.db.QueryContext(ctx, "SELECT * FROM images")
	if err != nil {
		return nil, fmt.Errorf("failed to get images from PostgreSQL: %w", err)
	}
	defer rows.Close()

	var images []m.Image
	for rows.Next() {
		var image m.Image
		err := rows.Scan(&image.ID, &image.S3URL, &image.Seed, &image.Prompt, &image.Tags, &image.Categories, &image.Styles, &image.CreatedAt, &image.CreatedBy)
		if err != nil {
			return nil, fmt.Errorf("failed to scan image row: %w", err)
		}
		images = append(images, image)
	}

	// Now, get the image data from S3
	var s3Urls []string
	for _, image := range images {
		s3Urls = append(s3Urls, image.S3URL)
	}

	imgData, err := repo.getImagesFromS3(ctx, s3Urls)
	if err != nil {
		return nil, fmt.Errorf("failed to get images from S3: %w", err)
	}

	// Finally, combine the image metadata and image data
	for i, image := range images {
		image.Data = imgData[i]
	}

	return images, nil

}

func (repo *ImageRepositoryImpl) StoreImage(ctx context.Context, image u.ImageMetadata, imgData []byte) error {
	id := u.GenerateID()
	s3url, err := repo.storeImageInS3(ctx, id, imgData)
	if err != nil {
		return fmt.Errorf("failed to store image with id %s in S3: %w", id, err)
	}
	// Now, store the metadata in the PostgreSQL database
	insertQuery := `INSERT INTO images (id, s3_url, seed, prompt, tags, categories, styles, created_at, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err = repo.db.ExecContext(ctx, insertQuery, id, s3url, image.Seed, image.Prompt, image.Tags, image.Categories, image.Styles, image.CreatedAt, image.CreatedBy)
	if err != nil {
		return fmt.Errorf("failed to store image with id %s in PostgreSQL: %w", id, err)
	}

	return nil
}

func (repo *ImageRepositoryImpl) storeImageInS3(ctx context.Context, id string, img []byte) (string, error) {
	// First, upload the image to S3
	key := fmt.Sprintf("images/%s.png", id)
	output, err := repo.s3Uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
		Key:    aws.String(key),
		Body:   bytes.NewReader(img),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image to S3: %w", err)
	}
	return output.Location, nil
}

func (repo *ImageRepositoryImpl) getImagesFromS3(ctx context.Context, s3Urls []string) ([][]byte, error) {
	// First, download the image from S3
	var images [][]byte
	objs := []s3manager.BatchDownloadObject{}
	for i, s3Url := range s3Urls {
		buf := aws.NewWriteAtBuffer(images[i])
		objs = append(objs, s3manager.BatchDownloadObject{
			Object: &s3.GetObjectInput{
				Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
				Key:    aws.String(s3Url),
			},
			Writer: buf,
		})
	}
	iter := &s3manager.DownloadObjectsIterator{Objects: objs}
	if err := repo.s3Downloader.DownloadWithIterator(aws.BackgroundContext(), iter); err != nil {
		return nil, err
	}
	return images, nil
}

// func (repo *ImageRepositoryImpl) getImageFromS3(ctx context.Context, s3Url string) ([]byte, error) {
// 	// First, download the image from S3
// 	buf := aws.NewWriteAtBuffer([]byte{})
// 	_, err := repo.s3Downloader.DownloadWithContext(ctx, buf, &s3.GetObjectInput{
// 		Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
// 		Key:    aws.String(s3Url),
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to download image from S3: %w", err)
// 	}
// 	return buf.Bytes(), nil
// }
