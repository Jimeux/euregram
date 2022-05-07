package data

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Image struct {
	PK        string    `dynamodbav:"pk" json:"-"` // user_id
	SK        string    `dynamodbav:"sk" json:"id"`
	UserID    string    `dynamodbav:"user_id" json:"-"`
	Username  string    `dynamodbav:"username" json:"username"`
	Caption   string    `dynamodbav:"caption" json:"caption"`
	Path      string    `dynamodbav:"url" json:"url"`
	CreatedAt time.Time `dynamodbav:"created_at" json:"created_at"`
}

type ImageRepository struct {
	table *string
	db    *dynamodb.Client
}

func NewImageRepository(table string, db *dynamodb.Client) *ImageRepository {
	return &ImageRepository{table: &table, db: db}
}

func (r *ImageRepository) ListImages(ctx context.Context) ([]Image, error) {
	res, err := r.db.Scan(ctx, &dynamodb.ScanInput{
		Limit:     aws.Int32(200),
		TableName: r.table,
	})
	if err != nil {
		return nil, fmt.Errorf("error scanning images: %v", err)
	}

	var images []Image
	if err := attributevalue.UnmarshalListOfMaps(res.Items, &images); err != nil {
		return nil, fmt.Errorf("error unmarshalling images: %v", err)
	}
	return images, nil
}

func (r *ImageRepository) SaveImage(ctx context.Context, userID, username, caption, imgPath string) (*Image, error) {
	id, createdAt := NewULID()
	rec := &Image{
		PK:        userID,
		SK:        id,
		UserID:    userID,
		Username:  username,
		Caption:   caption,
		Path:      imgPath,
		CreatedAt: createdAt,
	}

	av, err := attributevalue.MarshalMap(rec)
	if err != nil {
		return nil, fmt.Errorf("error marshalling new image item: %v", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: r.table,
	}
	if _, err := r.db.PutItem(ctx, input); err != nil {
		return nil, fmt.Errorf("error calling PutItem for image item: %v", err)
	}
	return rec, nil
}
