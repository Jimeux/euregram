package data

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type User struct {
	PK            string `dynamodbav:"pk"`
	Token         string `dynamodbav:"token"` // GSI
	Email         string `dynamodbav:"email"`
	GivenName     string `dynamodbav:"given_name"`
	FamilyName    string `dynamodbav:"family_name"`
	VerifiedEmail bool   `dynamodbav:"verified_email"`
	Picture       string `dynamodbav:"picture"`
	Locale        string `dynamodbav:"locale"`
}

func NewUser(id, email, gName, fName string, verifiedEmail bool, picture, token, locale string) *User {
	return &User{
		PK:            id,
		GivenName:     gName,
		FamilyName:    fName,
		Email:         email,
		VerifiedEmail: verifiedEmail,
		Picture:       picture,
		Token:         token,
		Locale:        locale,
	}
}

type UserRepository struct {
	table    *string
	gsiToken *string
	db       *dynamodb.Client
	prefix   string
}

func NewUserRepository(table string, db *dynamodb.Client) *UserRepository {
	return &UserRepository{
		table:    &table,
		gsiToken: aws.String("token"),
		db:       db,
	}
}

func (r *UserRepository) Save(ctx context.Context, user *User) error {
	av, err := attributevalue.MarshalMap(user)
	if err != nil {
		return fmt.Errorf("error marshalling new user item: %v", err)
	}
	if _, err := r.db.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      av,
		TableName: r.table,
	}); err != nil {
		return fmt.Errorf("error calling PutItem for user item: %v", err)
	}
	return nil
}

func (r *UserRepository) GetByToken(ctx context.Context, token string) (*User, error) {
	expr, err := expression.NewBuilder().WithKeyCondition(
		expression.KeyEqual(expression.Key("token"), expression.Value(token)),
	).Build()
	if err != nil {
		return nil, err
	}

	output, err := r.db.Query(ctx, &dynamodb.QueryInput{
		TableName:                 r.table,
		IndexName:                 r.gsiToken,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	})
	if err != nil {
		return nil, err
	}
	if len(output.Items) != 1 {
		return nil, nil
	}

	var user User
	if err := attributevalue.UnmarshalMap(output.Items[0], &user); err != nil {
		return nil, err
	}
	return &user, nil
}
