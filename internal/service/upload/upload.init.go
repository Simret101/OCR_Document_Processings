package upload

import (
	"aidoc/internal/repository"
	"aidoc/platform/logger"
	"aidoc/platform/ocr/amazontextract/async"
	"aidoc/platform/ocr/amazontextract/sync"
	"aidoc/platform/ocr/postprocessing"
	"aidoc/platform/ocr/preprocessing"
)

type Upload struct {
	log                    logger.Logger
	asynctextract          async.AsyncTextract
	synctextract           sync.SyncTextract
	preProcessing          preprocessing.Preprocessing
	postProcessing         postprocessing.PostProcessing
	maxSizeMB              int64
	storage                repository.Image
	syncMaxsizecap         int64
	asyncMaxsizecap        int64
	syncAllowedFileFormats []string
}

func NewUpload(log logger.Logger, asynctextract async.AsyncTextract, synctextract sync.SyncTextract, preProcessing preprocessing.Preprocessing, postProcessing postprocessing.PostProcessing, maxSizeMB int64, storage repository.Image, syncMaxsizecap int64, asyncMaxsizecap int64, syncAllowedFileFormats []string) *Upload {
	syncAllowedFileFormats = []string{"PDF", "JPEG", "PNG", "TIFF"}
	return &Upload{
		log:                    log,
		asynctextract: asynctextract,
		synctextract: synctextract,
		preProcessing:          preProcessing,
		postProcessing:         postProcessing,
		maxSizeMB:              maxSizeMB,
		storage:                storage,
		syncMaxsizecap:         syncMaxsizecap,
		asyncMaxsizecap:        asyncMaxsizecap,
		syncAllowedFileFormats: syncAllowedFileFormats,
	}
}
