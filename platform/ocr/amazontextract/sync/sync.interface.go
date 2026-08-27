package sync

import (
	"aidoc/internal/entity"
	"context"
)

type SyncTextract interface {
	AnalyzeExpense(ctx context.Context, input *entity.ReadInput) (*entity.OCRDocument, error)
}
