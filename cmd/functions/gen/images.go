package gen

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	utl "image-service/internal/util"
	"io"
	"net/http"
	"time"
)

type GenerateImageRequest struct {
	Prompt         string `json:"prompt"`
	NegativePrompt string `json:"negative_prompt"`
}

func Handler(ctx context.Context, request utl.Request) (utl.Response, error) {
	fmt.Printf("Request received: %s\n", request.Path)
	if request.HTTPMethod != "POST" || request.Path != "/images/generate" {
		return utl.ClientError(http.StatusMethodNotAllowed)
	}
	// Parse incoming request
	var genReq GenerateImageRequest
	err := json.Unmarshal([]byte(request.Body), &genReq)
	if err != nil {
		return utl.ServerError(err)
	}
	// Make a request to the Python FastAPI Lambda
	resp, err := genTextToImage(genReq)
	if err != nil {
		return utl.ServerError(err)
	}

	var images [][]byte
	var metadata []utl.ImageMetadata
	for i, artifact := range resp.Artifacts {
		imageBytes, err := utl.ConvertBase64ToImageBytes(artifact.Base64)
		if err != nil {
			return utl.ServerError(err)
		}
		metadata[i] = utl.ImageMetadata{
			Seed:         artifact.Seed,
			Prompt:       genReq.Prompt,
			FinishReason: artifact.FinishReason,
			CreatedAt:    fmt.Sprintf("%d", time.Now().Unix()),
			CreatedBy:    "Admin",
		}

		images[i] = imageBytes
	}
	// store image in S3
	err = storeBulkImagesAndMetadata(ctx, images, metadata)
	if err != nil {
		return utl.ServerError(err)
	}

	// Return the generated images to the client
	return utl.Response{
		StatusCode: http.StatusOK,
		Body:       "Images generated successfully",
	}, nil
}

func genTextToImage(req GenerateImageRequest) (*utl.GenTextToImageResponse, error) {
	// Marshal request to JSON
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	// Define the Python Lambda URL (update with actual URL)
	url := "https://url-to-ai-service"

	// Create POST request
	httpClient := &http.Client{Timeout: time.Second * 30}
	response, err := httpClient.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Read response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	// Unmarshal JSON response into GenTextToImageResponse struct
	var resp utl.GenTextToImageResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return &resp, nil
}
