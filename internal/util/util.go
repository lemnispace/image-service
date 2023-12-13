package util

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg" // Import the JPEG package for decoding
	_ "image/png"  // Import the PNG package for decoding
	"net/http"
)

func ClientError(status int) (Response, error) {
	return Response{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func ServerError(err error) (Response, error) {
	fmt.Printf("Internal server error: %v\n", err)
	return Response{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

func ConvertBase64ToImageBytes(base64Str string) ([]byte, error) {
	// Decode the base64 string to bytes
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
