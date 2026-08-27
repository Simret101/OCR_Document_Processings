package upload

import (
	"aidoc/internal/dto/upload"
	"aidoc/internal/entity"
	"aidoc/internal/errors"
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gocv.io/x/gocv"
)

// UploadAvatarHandler godoc
//
//	@Summary		Upload document
//	@Description	Uploads a document image (jpeg,jpg, png, webp), processes it, and returns the scanned ocr outputt.
//	@Tags			UploadDoc
//	@Accept			multipart/form-data
//	@Produce		json
//
//	@Param			upload			formData	file												true	"Upload image file (jpeg,jpg, png, webp)"
//
//	@Success		200				{object}	response.SuccessResponse{data=map[string]string}	"document uploaded successfully"
//	@Failure		400				{object}	response.ErrorResponse								"Bad Request"
//	@Failure		401				{object}	response.ErrorResponse								"Unauthorized"
//	@Failure		403				{object}	response.ErrorResponse								"Forbidden"
//	@Failure		500				{object}	response.ErrorResponse								"Internal Server Error"
//
//	@Router			/api/v1/upload [post]
func (s *Upload) UploadStatement(ctx context.Context, file multipart.File, header *multipart.FileHeader, format bool) (*upload.UploadStatementResponse, error) {

	data, contentType, err := s.readAndValidate(file, header)
	if err != nil {
		return nil, err
	}

	if format {
		data, err = s.preprocessIfImage(data, contentType)
		if err != nil {
			return nil, err
		}
	}

	doc, err := s.processOCR(ctx, data)
	if err != nil {
		return nil, err
	}

	documentID, err := s.uploadDocument(ctx, data, contentType)
	if err != nil {
		return nil, err
	}

	if doc != nil && doc.RawText != "" {
		txtKey := fmt.Sprintf(
			"statements/%s/%s.txt",
			time.Now().Format("2006/01"),
			documentID,
		)
		if _, err := s.storage.Upload(ctx, txtKey, strings.NewReader(doc.RawText), int64(len(doc.RawText)), "text/plain"); err != nil {
			s.log.Warn(ctx, "failed to store ocr text", zap.Error(err))
		}
	}

	return &upload.UploadStatementResponse{
		DocumentID: documentID,
		Status:     "uploaded",
		OCRText:    doc.RawText,
		Confidence: doc.Confidence,
	}, nil
}

func (s *Upload) readAndValidate(file multipart.File, header *multipart.FileHeader) ([]byte, string, error) {

	maxSize := s.maxSizeMB * 1024 * 1024
	if header.Size > maxSize {
		return nil, "", errors.ErrInvalidUserInput.New(
			"file size exceeds maximum allowed limit (%d MB)",
			s.maxSizeMB,
		)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, "", errors.ErrInvalidUserInput.Wrap(err, "failed to read file")
	}

	contentType := http.DetectContentType(data)

	allowed := map[string]bool{
		"application/pdf": true,
		"image/png":       true,
		"image/jpeg":      true,
		"image/tiff":      true,
	}

	if !allowed[contentType] {
		return nil, "", errors.ErrInvalidUserInput.New("unsupported file type: %s", contentType)
	}

	return data, contentType, nil
}

func (s *Upload) preprocessIfImage(data []byte, contentType string) ([]byte, error) {

	switch contentType {

	case "image/png",
		"image/jpeg",
		"image/tiff":

		img, err := gocv.IMDecode(data, gocv.IMReadColor)
		if err != nil {
			return nil, err
		}
		defer img.Close()

		result, err := s.preProcessing.RunPipeline(img, entity.PipelineStep{})
		if err != nil {
			return nil, err
		}
		defer result.Processed.Close()

		return result.CompressedData, nil

	default:
		return data, nil
	}
}

func (s *Upload) processOCR(ctx context.Context, data []byte) (*entity.OCRDocument, error) {

	doc, err := s.synctextract.AnalyzeExpense(ctx, &entity.ReadInput{
		Bytes: data,
	})
	if err != nil {
		return nil, err
	}

	cleaned, err := s.postProcessing.PostProcessPipeline(ctx, doc.RawText, doc)
	if err != nil {
		return nil, err
	}

	doc.RawText = cleaned

	return doc, nil
}

func (s *Upload) uploadDocument(ctx context.Context, data []byte, contentType string) (string, error) {

	id := uuid.NewString()

	ext := ".pdf"

	switch contentType {
	case "image/png":
		ext = ".png"

	case "image/jpeg":
		ext = ".jpg"

	case "image/tiff":
		ext = ".tiff"
	}

	objectKey := fmt.Sprintf(
		"statements/%s/%s%s",
		time.Now().Format("2006/01"),
		id,
		ext,
	)

	_, err := s.storage.Upload(ctx, objectKey, bytes.NewReader(data), int64(len(data)), contentType)

	if err != nil {
		return "", err
	}

	return id, nil
}
