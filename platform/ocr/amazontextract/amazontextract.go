package amazontextract

import (
	"context"
	"net/url"
	"strings"

	"aidoc/internal/entity"
	"aidoc/internal/errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/textract"
)

func NewClients(cfg entity.TextractConfig) (*textract.Client, *dynamodb.Client, error) {
	region := cfg.Region
	if cfg.Endpoint != "" {
		if endpointRegion := regionFromEndpoint(cfg.Endpoint); endpointRegion != "" {
			region = endpointRegion
		}
	}

	optFns := []func(*config.LoadOptions) error{
		config.WithRegion(region),
	}
	if cfg.AccessKeyID != "" && cfg.SecretAccessKey != "" {
		optFns = append(optFns, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, cfg.SessionToken),
		))
	}
	awsCfg, err := config.LoadDefaultConfig(context.Background(), optFns...)
	if err != nil {
		return nil, nil, errors.ErrUnableToGet.New("failed to load AWS config: %w", err)
	}

	textractClient := textract.NewFromConfig(awsCfg, func(o *textract.Options) {
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
		}
	})

	var dynamoClient *dynamodb.Client
	if cfg.TableName != "" {
		dynamoClient = dynamodb.NewFromConfig(awsCfg)
	}

	return textractClient, dynamoClient, nil
}

func regionFromEndpoint(endpoint string) string {
	u, err := url.Parse(endpoint)
	if err != nil {
		return ""
	}
	host := u.Host
	if i := strings.Index(host, "textract."); i >= 0 {
		rest := host[i+len("textract."):]
		if j := strings.Index(rest, "."); j > 0 {
			return rest[:j]
		}
	}
	return ""
}
