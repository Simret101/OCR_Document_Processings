package amazontextract

import (
	"aidoc/internal/entity"
	"aidoc/internal/errors"
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/google/uuid"
)

func CreateDocument(input *entity.ReadInput) (*entity.OCRDocument, error) {
	usercode := "5000003826"

	if input.Bytes == nil && (input.S3Bucket == "" || input.S3Key == "") {
		return nil, errors.ErrInvalidDocument.New("Error provide either image bytes or S3 location")
	}

	return &entity.OCRDocument{
		ID:            uuid.NewString(),
		ExtractedData: make(map[string]interface{}),
		ProcessedAt:   time.Now(),
		Status:        "processing",
		UserID:        usercode,
		Metadata:      input.Metadata,
	}, nil
}
func SaveDocument(ctx context.Context, dynamoClient *dynamodb.Client, tableName string, doc *entity.OCRDocument) {
	if dynamoClient == nil || tableName == "" {
		return
	}

	item, err := attributevalue.MarshalMap(doc)
	if err != nil {
		log.Printf("failed to marshal document for DynamoDB: %v", err)
		return
	}

	_, err = dynamoClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	if err != nil {
		log.Printf("failed to save document to DynamoDB: %v", err)
	}
}

func BuildDocument(input *entity.ReadInput) *types.Document {
	if input.Bytes != nil {
		return &types.Document{
			Bytes: input.Bytes,
		}
	}
	return &types.Document{
		S3Object: &types.S3Object{
			Bucket: &input.S3Bucket,
			Name:   &input.S3Key,
		},
	}
}
