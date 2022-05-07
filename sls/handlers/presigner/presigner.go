package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

var (
	presignClient *s3.PresignClient
	uploadBucket  = os.Getenv("UPLOAD_BUCKET")
	imageDomain   = os.Getenv("IMAGE_DOMAIN")
)

const (
	maxContentLength = 5e+6 // 5MB in bytes
	urlExpiry        = 2 * time.Minute
	uploadPath       = "upload"
)

func main() {
	cfg, _ := config.LoadDefaultConfig(context.Background())
	presignClient = s3.NewPresignClient(s3.NewFromConfig(cfg))

	lambda.Start(handler)
}

type request struct {
	ContentLength int64  `json:"contentLength"`
	ContentType   string `json:"contentType"`
}

type response struct {
	URL string `json:"url"`
}

type res events.APIGatewayV2HTTPResponse

func handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (res, error) {
	var req request
	if err := json.NewDecoder(strings.NewReader(event.Body)).Decode(&req); err != nil {
		return res{StatusCode: http.StatusBadRequest, Body: "Could not parse request body"}, err
	}
	if !validateContentType(req.ContentType) {
		return res{StatusCode: http.StatusBadRequest, Body: "Invalid Content-Type"}, nil
	}
	if req.ContentLength > maxContentLength {
		return res{StatusCode: http.StatusBadRequest, Body: "Maximum Content-Length is 5MB"}, nil
	}

	uri, err := presignURL(ctx, req.ContentLength, req.ContentType)
	if err != nil {
		return res{StatusCode: http.StatusInternalServerError, Body: "Could not pre-sign"}, err
	}
	body, err := json.Marshal(&response{uri})
	if err != nil {
		return res{StatusCode: http.StatusInternalServerError, Body: "Could not prepare response"}, err
	}
	return res{StatusCode: http.StatusOK, Body: string(body)}, nil
}

func presignURL(ctx context.Context, contentLength int64, contentType string) (string, error) {
	fileName := uuid.New().String()
	input := &s3.PutObjectInput{
		ContentType:   aws.String(contentType),
		ContentLength: contentLength,
		Bucket:        aws.String(uploadBucket),
		Key:           aws.String(path.Join(uploadPath, fileName)),
	}
	req, err := presignClient.PresignPutObject(ctx, input, s3.WithPresignExpires(urlExpiry))
	if err != nil {
		return "", err
	}
	return convertURL(req.URL), nil
}

func convertURL(s3URL string) string {
	parsed, _ := url.Parse(s3URL)
	parsed.Host = imageDomain
	return parsed.String()
}

func validateContentType(contentType string) bool {
	switch contentType {
	case "image/gif",
		"image/heic",
		"image/jpeg",
		"image/png",
		"image/webp":
		return true
	}
	return false
}
