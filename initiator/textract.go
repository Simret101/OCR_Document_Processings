package initiator

import (
	"context"

	"aidoc/internal/entity"
	"aidoc/platform/logger"
	"aidoc/platform/ocr/amazontextract/async"
	"aidoc/platform/ocr/amazontextract/sync"

	"github.com/spf13/viper"
)


func initSyncReader(log logger.Logger) *sync.Reader {

	cfg := entity.TextractConfig{
		Region:          viper.GetString("aws.region"),
		Endpoint:        viper.GetString("aws.endpoint"),
		AccessKeyID:     viper.GetString("aws.access_key_id"),
		SecretAccessKey: viper.GetString("aws.secret_access_key"),
		SessionToken:    viper.GetString("aws.session_token"),
		TableName:       viper.GetString("aws.table_name"),
		BucketName:      viper.GetString("aws.bucket_name"),
		KMSKeyID:        viper.GetString("aws.kms_key_id"),
		Timeout:         viper.GetDuration("aws.timeout"),
		MaxRetries:      viper.GetInt("aws.max_retries"),
	}

	reader, err := sync.New(cfg)

	if err != nil {
		log.Fatal(
			context.Background(),
			"failed initializing sync textract reader",
		)
	}

	return reader
}
func initAsyncReader(log logger.Logger) *async.Reader {

	cfg := entity.TextractConfig{
		Region:          viper.GetString("aws.region"),
		Endpoint:        viper.GetString("aws.endpoint"),
		AccessKeyID:     viper.GetString("aws.access_key_id"),
		SecretAccessKey: viper.GetString("aws.secret_access_key"),
		SessionToken:    viper.GetString("aws.session_token"),
		TableName:       viper.GetString("aws.table_name"),
		BucketName:      viper.GetString("aws.bucket_name"),
		KMSKeyID:        viper.GetString("aws.kms_key_id"),
		Timeout:         viper.GetDuration("aws.timeout"),
		MaxRetries:      viper.GetInt("aws.max_retries"),
	}

	reader, err := async.New(cfg)

	if err != nil {
		log.Fatal(
			context.Background(),
			"failed initializing async textract reader",
		)
	}

	return reader
}
