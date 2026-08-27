package async

import (
	"aidoc/internal/entity"
	"context"
	"time"
)

type AsyncTextract interface {
	GetAnalysisResult(ctx context.Context, jobID string, features []entity.TextractFeature, timeout time.Duration) (*entity.OCRDocument, error)
	StartDocumentAnalysis(ctx context.Context, input *entity.ReadInput, features []entity.TextractFeature) (*string, error) 
}
