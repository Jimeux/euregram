package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/Jimeux/euregram/sls/data"
	"github.com/Jimeux/euregram/sls/jwt"
)

var (
	stateTable = os.Getenv("STATE_TABLE")
	userTable  = os.Getenv("USER_TABLE")
	jwtSecret  = []byte(os.Getenv("JWT_SECRET"))
	oauthConf  = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes:       []string{"email", "profile"},
	}

	stateRepo *data.StateRepository
	userRepo  *data.UserRepository
)

type ConfirmPayload struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

type res events.APIGatewayV2HTTPResponse

func main() {
	cfg, _ := config.LoadDefaultConfig(context.Background())
	ddb := dynamodb.NewFromConfig(cfg)
	stateRepo = data.NewAuthRepository(stateTable, ddb)
	userRepo = data.NewUserRepository(userTable, ddb)

	lambda.Start(handleRequest)
}

func handleRequest(ctx context.Context, event events.APIGatewayV2HTTPRequest) (res, error) {
	switch event.RouteKey {
	case "GET /api/auth/init":
		result, err := handleInit(ctx)
		if err != nil {
			log.Println(err)
		}
		return result, err
	case "POST /api/auth/confirm":
		result, err := handleConfirm(ctx, []byte(event.Body))
		if err != nil {
			log.Println(err)
		}
		return result, err
	default:
		return res{StatusCode: http.StatusNotFound}, nil
	}
}

func handleInit(ctx context.Context) (res, error) {
	state, err := stateRepo.GenerateState(ctx)
	if err != nil {
		return res{StatusCode: http.StatusInternalServerError}, err
	}
	redirectURL := oauthConf.AuthCodeURL(state.PK, oauth2.AccessTypeOffline, oauth2.ApprovalForce)

	return res{
		StatusCode: http.StatusOK,
		Body:       `{"redirect_url": "` + redirectURL + `"}`,
	}, nil
}

func handleConfirm(ctx context.Context, body []byte) (res, error) {
	var payload ConfirmPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return res{
			StatusCode: http.StatusBadRequest,
			Body:       `{"message": "invalid request body"}`,
		}, err
	}

	state, err := stateRepo.GetState(ctx, payload.State)
	if err != nil {
		return res{StatusCode: http.StatusInternalServerError}, nil
	}
	if state == nil {
		return res{StatusCode: http.StatusForbidden}, err
	}
	user, err := confirmCode(ctx, payload.Code)
	if err != nil {
		return res{
			StatusCode: http.StatusBadRequest,
			Body:       `{"message": "invalid request params: ` + err.Error() + `"}`,
		}, err
	}

	return res{
		StatusCode: http.StatusOK,
		Body: `{
				 "user_id": "` + user.PK + `",
				 "username": "` + user.GivenName + `",                     
				 "access_token": "` + user.Token + `",
				 "locale": "` + user.Locale + `"
			   }`,
	}, nil
}

func confirmCode(ctx context.Context, code string) (*data.User, error) {
	googleToken, err := oauthConf.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodGet,
		"https://www.googleapis.com/oauth2/v2/userinfo?access_token="+googleToken.AccessToken, nil)
	req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ui GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&ui); err != nil {
		return nil, err
	}

	jwtToken, err := jwt.Generate(ui.ID, ui.GivenName, jwtSecret)
	if err != nil {
		return nil, err
	}

	user := data.NewUser(
		ui.ID, ui.Email, ui.GivenName, ui.FamilyName, ui.VerifiedEmail, ui.Picture, jwtToken, ui.Locale)
	if err := userRepo.Save(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
