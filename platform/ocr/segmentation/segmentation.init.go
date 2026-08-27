package segmentation

import "aidoc/platform/logger"

type Segmentation struct {
	maxVal float64
	log    logger.Logger
}

func NewSegmentation(maxVal float64, log logger.Logger) *Segmentation {
	return &Segmentation{
		maxVal: maxVal,
		log:    log,
	}
}
