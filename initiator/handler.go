package initiator

import (
	"aidoc/internal/handler"
	uploadhandler "aidoc/internal/handler/upload"
	"aidoc/platform/logger"

	"github.com/spf13/viper"
)

type Handler struct {
	Upload handler.Upload
}

func initHandler(
	service *Service,
	log logger.Logger,
) *Handler {

	uploadHandler := uploadhandler.NewUpload(
		log,
		service.Upload,
		[]string{
			"application/pdf",
			"application/jpeg",
			"application/png",
			"application/tiff",
		},
		viper.GetInt64("upload.sync_max_size"),
		viper.GetInt64("upload.async_max_size"),
	)

	return &Handler{
		Upload: uploadHandler,
	}
}
