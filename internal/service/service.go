package service

import (
	"aidoc/internal/dto/upload"
	"context"
	"mime/multipart"
)

type Upload interface {
	UploadStatement(ctx context.Context,file multipart.File,header *multipart.FileHeader,format bool) (*upload.UploadStatementResponse, error)
}
