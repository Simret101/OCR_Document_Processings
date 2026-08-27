package initiator

import (
	"aidoc/platform/logger"
	"aidoc/platform/ocr/segmentation"

	"github.com/spf13/viper"
)

func initSegmentation(
	log logger.Logger,
) *segmentation.Segmentation {

	return segmentation.NewSegmentation(
		viper.GetFloat64("ocr.segmentation.max_value"),
		log,
	)
}
