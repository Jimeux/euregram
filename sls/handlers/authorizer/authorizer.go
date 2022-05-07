package main

import (
	"context"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/Jimeux/euregram/sls/data"
	"github.com/Jimeux/euregram/sls/jwt"
)

var (
	userTable = os.Getenv("USER_TABLE")
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	userRepo  *data.UserRepository
)

var deny = events.APIGatewayV2CustomAuthorizerSimpleResponse{IsAuthorized: false}

func main() {
	cfg, _ := config.LoadDefaultConfig(context.Background())
	ddb := dynamodb.NewFromConfig(cfg)
	userRepo = data.NewUserRepository(userTable, ddb)

	lambda.Start(handleRequest)
}

func handleRequest(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {
	token := strings.TrimPrefix(event.Headers["authorization"], "Bearer ")
	if _, err := jwt.ParseAndValidate(token, jwtSecret); err != nil {
		return deny, nil
	}

	// TODO 2022/03/09 @Jimeux Could refresh token here

	user, err := userRepo.GetByToken(ctx, token)
	if err != nil {
		return deny, nil
	}
	if user == nil {
		return deny, nil
	}
	return events.APIGatewayV2CustomAuthorizerSimpleResponse{
		IsAuthorized: true,
		Context: map[string]any{
			"user_id":  user.PK,
			"username": user.GivenName,
		},
	}, nil
}
