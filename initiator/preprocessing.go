package initiator

import (
	"aidoc/internal/entity"
	"aidoc/platform/logger"
	"aidoc/platform/ocr/preprocessing"
	"aidoc/platform/ocr/preprocessing/compression"
	"image"

	"github.com/minio/minio-go/v7"
	"github.com/spf13/viper"
)

func initPreProcessing(minioClient *minio.Client,
	log logger.Logger,
) preprocessing.Preprocessing {

	compressor := initCompression(log)

	return preprocessing.NewStorage(
		minioClient,
		viper.GetString("minio.bucket"),
		viper.GetInt("ocr.width"),
		viper.GetInt("ocr.height"),
		viper.GetInt("ocr.max_dimension"),
		viper.GetFloat64("ocr.clip_limit"),
		image.Point{
			X: viper.GetInt("ocr.tile_width"),
			Y: viper.GetInt("ocr.tile_height"),
		},
		viper.GetInt("ocr.kernel_size"),
		compressor,
	)
}

func initCompression(log logger.Logger) *compression.Compression {

	bitReader := &entity.BitReader{}

	bitWriter := &entity.BitWriter{}

	return compression.NewCompression(
		log,
		bitReader,
		bitWriter,
	)
}
