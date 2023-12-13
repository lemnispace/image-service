package util

import (
	"github.com/aws/aws-lambda-go/events"
)

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

// Define the structure for your request to the ai-service
type GenTextToImageRequest struct {
	Prompt         string `json:"prompt" description:"Text prompt with description of the things you want in the image to be generated"`
	NegativePrompt string `json:"negative_prompt,omitempty" description:"Items you don't want in the image"`
	Seed           *int   `json:"seed,omitempty" description:"Seed is used to reproduce results, same seed will give you same image in return again. Pass null for a random number."`
	Steps          *int   `json:"steps,omitempty" description:"Number of steps to run the model for"`
	Samples        *int   `json:"samples,omitempty" description:"Number of samples to generate"`
	Width          *int   `json:"width,omitempty" description:"Width of the image"`
	Height         *int   `json:"height,omitempty" description:"Height of the image"`
}

// Define the structure for the response from the ai-service
type GenTextToImageResponse struct {
	Artifacts []struct {
		Base64       string `json:"base64"`
		Seed         int    `json:"seed"`
		FinishReason string `json:"finishReason"`
	} `json:"artifacts"`
}

// Define the structure for the metadata of the image
type ImageMetadata struct {
	ID           string `json:"id"`
	Prompt       string `json:"prompt"`
	Seed         int    `json:"seed"`
	CreatedAt    string `json:"created_at"`
	CreatedBy    string `json:"created_by"`
	FinishReason string `json:"finish_reason"`
}
