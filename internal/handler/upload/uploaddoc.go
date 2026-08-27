package upload

import (
	"aidoc/internal/constant"
	"aidoc/internal/errors"
	"aidoc/internal/response"
	"net/http"

	"go.uber.org/zap"
)

// UploadDocument accepts multipart form data: upload + optional 'name' field.
// @Summary Upload document
// @Description Uploads a document. Text is extracted (txt, md, csv, json, PDF, DOCX), summarized with AI, and indexed into the tenant knowledge base for semantic search.
// @Tags Documents
// @Accept multipart/form-data
// @Produce json
// @Param upload formData file true "Document file"
// @Success 200 {object} string "document uploaded"
// @Failure 400 {object} response.ErrorResponseFormat
// @Failure 403 {object} response.ErrorResponseFormat
// @Router /api/v1/upload [post]
func (h *Upload) UploadStatmentHandler(w http.ResponseWriter, r *http.Request) {
	var micCtx = constant.HANDER_UPLOAD
	file, header, err := r.FormFile("upload")
	if err != nil {
		h.log.Named(micCtx).Error(r.Context(), "missing file", zap.Error(err))

		response.SendErrorResponseFormated(w, errors.ErrInvalidUserInput.New("national id file is required"))
		return
	}
	contentType := header.Header.Get("Content-Type")
	var format bool
	if header.Size < h.syncMaxsizecap {
		for _, val := range h.syncAllowedFileFormats {
			if contentType != val {
				response.SendErrorResponseFormated(w, errors.ErrInvalidUserInput.New("file format not allowed."))
				return
			}
			format = true
			break

		}

	} else if header.Size > h.asyncMaxsizecap {
		response.SendErrorResponseFormated(w, errors.ErrInvalidUserInput.New("maximum allowed capacity passed."))
		return

	} else {
		format = false
	}

	defer file.Close()

	res, err := h.service.UploadStatement(r.Context(), file, header, format)
	if err != nil {
		h.log.Named(micCtx).Error(r.Context(), "upload failed", zap.Error(err))

		response.SendErrorResponseFormated(w, err)
		return
	}
	h.log.Named(micCtx).Info(r.Context(), "upload successfull")

	response.SendSuccessResponse(w, http.StatusOK, response.UPLOADDOCUMENT, res, "")
}
