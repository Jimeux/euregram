package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/Jimeux/euregram/sls/data"
)

var (
	imageDomain = os.Getenv("IMAGE_DOMAIN")
	imageTable  = os.Getenv("IMAGE_TABLE")
	ddb         *dynamodb.Client
	repository  *data.ImageRepository
)

type response struct {
	Images []data.Image `json:"images"`
}

type res events.APIGatewayV2HTTPResponse

func main() {
	cfg, _ := config.LoadDefaultConfig(context.Background())
	ddb = dynamodb.NewFromConfig(cfg)
	repository = data.NewImageRepository(imageTable, ddb)

	lambda.Start(handler)
}

func handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (res, error) {
	images, err := repository.ListImages(ctx)
	if err != nil {
		return res{StatusCode: http.StatusInternalServerError, Body: "Could not retrieve images"}, err
	}

	for i := range images {
		images[i].Path = "https://" + imageDomain + "/" + images[i].Path
	}

	body, err := json.Marshal(&response{images})
	if err != nil {
		return res{StatusCode: http.StatusInternalServerError, Body: "Could not serialize response"}, err
	}
	return res{StatusCode: http.StatusOK, Body: string(body)}, nil
}
