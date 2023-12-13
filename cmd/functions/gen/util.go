package gen

// Import necessary AWS SDK packages
import (
	"bytes"
	"context"
	"fmt"
	utl "image-service/internal/util"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Function to bulk store images in S3 and metadata in DynamoDB
func storeBulkImagesAndMetadata(ctx context.Context, imgData [][]byte, metadata []utl.ImageMetadata) error {
	// Initialize AWS session
	sess := session.Must(session.NewSession())

	// Store images in S3 concurrently
	err := bulkUploadToS3(ctx, sess, imgData, metadata)
	if err != nil {
		return fmt.Errorf("failed to bulk upload images to S3: %w", err)
	}

	// Store metadata in DynamoDB in batches
	err = bulkWriteToDynamoDB(ctx, sess, metadata)
	if err != nil {
		return fmt.Errorf("failed to bulk write metadata to DynamoDB: %w", err)
	}

	return nil
}

func uploadToS3(sess *session.Session, imgData []byte, metadata utl.ImageMetadata) error {
	s3Uploader := s3manager.NewUploader(sess)
	_, err := s3Uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("your-s3-bucket-name"),
		Key:    aws.String(metadata.ID + ".png"),
		Body:   bytes.NewReader(imgData),
	})
	if err != nil {
		return fmt.Errorf("failed to upload image to S3: %w", err)
	}

	return nil
}

func bulkUploadToS3(ctx context.Context, sess *session.Session, imgData [][]byte, metadata []utl.ImageMetadata) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(imgData))

	for i, data := range imgData {
		wg.Add(1)
		go func(i int, data []byte) {
			defer wg.Done()

			err := uploadToS3(sess, data, metadata[i])
			if err != nil {
				errChan <- err
				return
			}
		}(i, data)
	}

	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func bulkWriteToDynamoDB(ctx context.Context, sess *session.Session, metadata []utl.ImageMetadata) error {
	dynamoDBSvc := dynamodb.New(sess)

	// DynamoDB batch write can handle up to 25 items at a time
	for i := 0; i < len(metadata); i += 25 {
		end := i + 25
		if end > len(metadata) {
			end = len(metadata)
		}

		writeRequests := make([]*dynamodb.WriteRequest, 0, end-i)
		for _, item := range metadata[i:end] {
			av, err := dynamodbattribute.MarshalMap(item)
			if err != nil {
				return fmt.Errorf("failed to marshal metadata: %w", err)
			}

			writeRequests = append(writeRequests, &dynamodb.WriteRequest{
				PutRequest: &dynamodb.PutRequest{
					Item: av,
				},
			})
		}

		_, err := dynamoDBSvc.BatchWriteItem(&dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]*dynamodb.WriteRequest{
				"your-dynamodb-table-name": writeRequests,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to batch write to DynamoDB: %w", err)
		}
	}

	return nil
}
