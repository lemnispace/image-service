package util

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"time"

	"github.com/google/uuid"
)

// ConvertBase64ToImageBytes decodes a base64 string into image bytes.
// It checks if the decoded bytes can be decoded as an image.
// Returns the image bytes if successful, otherwise returns an error.
func ConvertBase64ToImageBytes(base64Str string) ([]byte, error) {
	imageBytes, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}

	// Check if the decoded bytes can be decoded as an image
	_, _, err = image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, errors.New("invalid image format")
	}

	return imageBytes, nil
}

// GenerateID generates a unique ID using UUID.
// Returns the generated ID as a string.
func GenerateID() string {
	// generate unique ID
	id := uuid.New()
	return id.String()
}

func GetCurrentTime() int64 {
	return time.Now().Unix()
}
