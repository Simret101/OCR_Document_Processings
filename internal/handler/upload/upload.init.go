package upload

import (
	"aidoc/internal/service"
	"aidoc/platform/logger"
)

type Upload struct {
	log                    logger.Logger
	service                service.Upload
	syncAllowedFileFormats []string
	syncMaxsizecap         int64
	asyncMaxsizecap        int64
}

func NewUpload(log logger.Logger, service service.Upload, syncAllowedFileFormats []string, syncMaxsizecap int64, asyncMaxsizecap int64) *Upload {
	return &Upload{
		log:                    log,
		service:                service,
		syncAllowedFileFormats: syncAllowedFileFormats,
		syncMaxsizecap: syncMaxsizecap,
		asyncMaxsizecap: asyncMaxsizecap,
	}
}
