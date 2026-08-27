package compression

import (
	"aidoc/internal/entity"
	"aidoc/platform/logger"
)

type Compression struct {
	log logger.Logger
	bitreader *entity.BitReader
	bitWriter *entity.BitWriter
}
func NewCompression(log logger.Logger,bitreader *entity.BitReader,bitWriter *entity.BitWriter)*Compression{
	return &Compression{
		log:log,
		bitreader: bitreader,
		bitWriter: bitWriter,
	}
}
