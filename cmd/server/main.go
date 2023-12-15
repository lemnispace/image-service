package main

import (
	"image-service/internal/api"
	db "image-service/internal/db"

	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	repo, err := db.NewImageRepository()
	if err != nil {
		panic(err)
	}
	err = api.SetupRoutes(r, repo)
	if err != nil {
		panic(err)
	}
	ginLambda := ginadapter.New(r)
	lambda.Start(ginLambda.ProxyWithContext)
}
