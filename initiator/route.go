package initiator

import (
	"aidoc/internal/routing/upload"
	"aidoc/platform/logger"

	"github.com/go-chi/chi/v5"
)

func initRoute(grp chi.Router, log logger.Logger, handler *Handler) {
	upload.Upload(grp, handler.Upload, log)
	

}
