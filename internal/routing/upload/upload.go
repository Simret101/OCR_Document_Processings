package upload

import (
	"aidoc/internal/handler"
	"aidoc/internal/routing/mainroute"
	"net/http"

	"aidoc/platform/logger"

	"github.com/go-chi/chi/v5"
)

func Upload(grp chi.Router, handle handler.Upload, log logger.Logger) {
	upload := []mainroute.Route{
		{
			Method:  http.MethodPost,
			Path:    "/api/v1/upload",
			Handler: handle.UploadStatmentHandler,
		},
	}
	mainroute.RegisterRoute(grp, upload, log)
}
