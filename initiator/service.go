package initiator

import (
	"context"

	"aidoc/internal/entity"
	"aidoc/internal/service"
	uploadservice "aidoc/internal/service/upload"
	"aidoc/platform/logger"
	"aidoc/platform/ocr/amazontextract/async"
	synctextract "aidoc/platform/ocr/amazontextract/sync"
	"aidoc/platform/ocr/tesseract"

	"github.com/minio/minio-go/v7"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Service struct {
	Upload service.Upload
}

func initService(
	persistence *Persistence,
	minioClient *minio.Client,

	log logger.Logger,
) *Service {

	textractConfig := entity.TextractConfig{
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

	asyncReader, err := async.New(textractConfig)

	if err != nil {
		log.Fatal(
			context.Background(),
			"failed initializing async textract",
		)
	}

	var syncReader synctextract.SyncTextract
	switch viper.GetString("ocr.provider") {
	case "tesseract":
		tessCfg := tesseract.Config{
			BinaryPath:  viper.GetString("tesseract.binary_path"),
			Languages:   viper.GetString("tesseract.languages"),
			PageSegMode: viper.GetInt("tesseract.page_seg_mode"),
			Timeout:     viper.GetDuration("tesseract.timeout"),
		}
		syncReader, err = tesseract.New(tessCfg)
		if err != nil {
			log.Warn(context.Background(), "failed initializing tesseract", zap.Error(err))
		}
	default:
		syncReader, err = synctextract.New(textractConfig)
		if err != nil {
			log.Fatal(
				context.Background(),
				"failed initializing sync textract",
			)
		}
	}
	preprocessing := initPreProcessing(minioClient, log)
	postprocessing:=initPostProcessing()

	uploadSvc := uploadservice.NewUpload(
		log,
		asyncReader,
		syncReader,
		preprocessing,
		postprocessing,
		viper.GetInt64("upload.max_size_mb"),
		persistence.ImageStore,
		viper.GetInt64("upload.sync_max_size"),
		viper.GetInt64("upload.async_max_size"),
		[]string{
			"PDF",
			"JPEG",
			"PNG",
			"TIFF",
		},
	)

	return &Service{
		Upload: uploadSvc,
	}
}
