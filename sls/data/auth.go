package data

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	stateTTL = time.Hour
)

type State struct {
	PK       string    `dynamodbav:"pk"`
	CreateAt time.Time `dynamodbav:"created_at"`
	TTL      int64     `dynamodbav:"ttl"`
}

type StateRepository struct {
	table *string
	db    *dynamodb.Client
}

func NewAuthRepository(table string, db *dynamodb.Client) *StateRepository {
	return &StateRepository{table: &table, db: db}
}

func (r *StateRepository) GenerateState(ctx context.Context) (*State, error) {
	id, createdAt := NewULID()
	rec := &State{
		PK:       id,
		TTL:      TTL(stateTTL),
		CreateAt: createdAt,
	}

	av, err := attributevalue.MarshalMap(rec)
	if err != nil {
		return nil, fmt.Errorf("error marshalling new auth item: %v", err)
	}
	if _, err := r.db.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      av,
		TableName: r.table,
	}); err != nil {
		return nil, fmt.Errorf("error calling PutItem for auth item: %v", err)
	}
	return rec, nil
}

func (r *StateRepository) GetState(ctx context.Context, st string) (*State, error) {
	item, err := r.db.GetItem(ctx, &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: st},
		},
		TableName:      r.table,
		ConsistentRead: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	if item.Item == nil {
		return nil, nil
	}

	// FIXME 2022/03/11 @Jimeux Delete or invalidate directly after first access

	var state State
	if err := attributevalue.UnmarshalMap(item.Item, &state); err != nil {
		return nil, err
	}
	return &state, nil
}
