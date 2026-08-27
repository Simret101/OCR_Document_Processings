package postprocessing

import (
	"aidoc/internal/entity"
	"context"
)

type PostProcessing interface {
	PostProcessPipeline(ctx context.Context, text string, doc *entity.OCRDocument) (string, error)
}
