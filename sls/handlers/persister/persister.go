package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/Jimeux/euregram/sls/data"
)

var (
	imageTable   = os.Getenv("IMAGE_TABLE")
	imageBucket  = os.Getenv("IMAGE_BUCKET")
	uploadBucket = os.Getenv("UPLOAD_BUCKET")
	imageHost    = os.Getenv("IMAGE_DOMAIN")

	ddbClient  *dynamodb.Client
	s3Client   *s3.Client
	repository = data.NewImageRepository(imageTable, ddbClient)
)

func main() {
	cfg, _ := config.LoadDefaultConfig(context.Background())
	ddbClient = dynamodb.NewFromConfig(cfg)
	s3Client = s3.NewFromConfig(cfg)
	repository = data.NewImageRepository(imageTable, ddbClient)

	lambda.Start(handler)
}

type res events.APIGatewayV2HTTPResponse

type urlPayload struct {
	URL     string `json:"url"`
	Caption string `json:"caption"`
}

func handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (res, error) {
	userID, ok := event.RequestContext.Authorizer.Lambda["user_id"].(string)
	if !ok || userID == "" {
		return res{StatusCode: http.StatusBadRequest, Body: "Invalid user_id"}, nil
	}
	username, ok := event.RequestContext.Authorizer.Lambda["username"].(string)
	if !ok || username == "" {
		return res{StatusCode: http.StatusBadRequest, Body: "Invalid username"}, nil
	}

	var u urlPayload
	if err := json.NewDecoder(strings.NewReader(event.Body)).Decode(&u); err != nil {
		return res{StatusCode: http.StatusInternalServerError, Body: "Could not parse request"}, nil
	}
	imgPath, err := transferImage(ctx, u)
	if err != nil {
		return res{StatusCode: http.StatusInternalServerError, Body: "failed to move image"}, err
	}

	img, err := repository.SaveImage(ctx, userID, username, u.Caption, imgPath)
	if err != nil {
		return res{StatusCode: http.StatusInternalServerError, Body: "failed to save image"}, err
	}

	fullURL, _ := url.Parse("https://" + imageHost)
	fullURL.Path = img.Path
	img.Path = fullURL.String()
	body, err := json.Marshal(img)
	if err != nil {
		return res{StatusCode: http.StatusInternalServerError, Body: "Failed to marshal image"}, err
	}
	return res{StatusCode: http.StatusOK, Body: string(body)}, nil
}

func transferImage(ctx context.Context, u urlPayload) (string, error) {
	parsed, err := url.Parse(u.URL)
	if err != nil {
		return "", err
	}

	// FIXME In a realistic scenario, we'd want some way of verifying
	//  the current user's privileges for the related S3 objects, e.g.
	//  including a user ID in the pre-signed path to be later validated.

	p := strings.TrimPrefix(parsed.Path, "/upload")
	fileName := path.Join("images", p)

	if _, err := s3Client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(imageBucket),
		CopySource: aws.String(url.PathEscape(path.Join(uploadBucket, parsed.Path))),
		Key:        aws.String(fileName),
	}); err != nil {
		return "", err
	}
	return fileName, nil
}
