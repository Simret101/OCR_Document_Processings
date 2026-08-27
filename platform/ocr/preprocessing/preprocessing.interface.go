package preprocessing

import (
	"aidoc/internal/entity"

	"gocv.io/x/gocv"
)

type Preprocessing interface {
	RunPipeline(src gocv.Mat, cfg entity.PipelineStep) (*PipelineResult, error)
}
