package sync

import (
	"aidoc/internal/entity"
	"aidoc/platform/ocr/amazontextract"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/textract"
)

type Reader struct {
	textractClient *textract.Client
	dynamoClient   *dynamodb.Client
	config         entity.TextractConfig
}

func New(cfg entity.TextractConfig) (*Reader, error) {
	textractClient, dynamoClient, err := amazontextract.NewClients(cfg)
	if err != nil {
		return nil, err
	}
	return &Reader{
		textractClient: textractClient,
		dynamoClient:   dynamoClient,
		config:         cfg,
	}, nil
}
