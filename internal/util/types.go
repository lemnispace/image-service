package util

import (
	"github.com/aws/aws-lambda-go/events"
)

// Response represents the structure for the API Gateway proxy response.
type Response events.APIGatewayProxyResponse

// Request represents the structure for the API Gateway proxy request.
type Request events.APIGatewayProxyRequest

// APIGenTextToImageRequest represents the structure for the request to the AI service for generating an image from text.
type APIGenTextToImageRequest struct {
	Prompt         string `json:"prompt" description:"Text prompt with description of the things you want in the image to be generated"`
	NegativePrompt string `json:"negative_prompt,omitempty" description:"Items you don't want in the image"`
	Seed           *int   `json:"seed,omitempty" description:"Seed is used to reproduce results, same seed will give you same image in return again. Pass null for a random number."`
	Steps          *int   `json:"steps,omitempty" description:"Number of steps to run the model for"`
	Samples        *int   `json:"samples,omitempty" description:"Number of samples to generate"`
	Width          *int   `json:"width,omitempty" description:"Width of the image"`
	Height         *int   `json:"height,omitempty" description:"Height of the image"`
}

// APIGenTextToImageResponse represents the structure for the response from the AI service for generating an image from text.
type APIGenTextToImageResponse struct {
	Artifacts []struct {
		Base64       string `json:"base64"`
		Seed         int    `json:"seed"`
		FinishReason string `json:"finishReason"`
	} `json:"artifacts"`
}

// ImageMetadata represents the structure for the metadata of an image.
type ImageMetadata struct {
	Seed         int      `json:"seed"`
	Prompt       string   `json:"prompt"`
	FinishReason string   `json:"finish_reason"`
	CreatedAt    int64    `json:"created_at"`
	CreatedBy    string   `json:"created_by"`
	Styles       []string `json:"styles"`
	Categories   []string `json:"categories"`
	Tags         []string `json:"tags"`
}
