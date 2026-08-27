package preprocessing

import (
	"aidoc/internal/entity"
	"errors"

	"gocv.io/x/gocv"
)

func (s *Storage) RunPipeline(src gocv.Mat, cfg entity.PipelineStep) (*PipelineResult, error) {
	res := &PipelineResult{}

	current := src.Clone()

	run := func(next *gocv.Mat, err error) error {
		if err != nil {
			current.Close()
			res.LastError = err
			return err
		}

		if next == nil {
			current.Close()
			err = errors.New("pipeline step returned nil")
			res.LastError = err
			return err
		}

		current.Close()
		current = *next
		res.StepsProcessed++

		return nil
	}

	if err := run(s.stepGray(&current)); err != nil {
		return res, err
	}

	if err := run(s.stepDeskew(&current, cfg)); err != nil {
		return res, err
	}

	if err := run(s.stepDewarp(&current, cfg)); err != nil {
		return res, err
	}

	if err := run(s.stepResize(&current, cfg)); err != nil {
		return res, err
	}

	if err := run(s.stepContrast(&current)); err != nil {
		return res, err
	}

	if err := run(s.stepRemoveLines(&current)); err != nil {
		return res, err
	}

	if err := run(s.stepDenoise(&current, cfg)); err != nil {
		return res, err
	}

	if err := run(s.stepFilter(&current)); err != nil {
		return res, err
	}

	if err := run(s.stepThreshold(&current)); err != nil {
		return res, err
	}

	if err := run(s.stepColorSpace(&current)); err != nil {
		return res, err
	}

	if err := run(s.stepColorTransfer(&current, cfg)); err != nil {
		return res, err
	}

	compressed, err := s.stepCompress(&current)
	if err != nil {
		current.Close()
		res.LastError = err
		return res, err
	}

	res.Processed = &current
	res.CompressedData = compressed

	return res, nil
}

type PipelineResult struct {
	Processed      *gocv.Mat
	CompressedData []byte
	StepsProcessed int
	LastError      error
}

func (s *Storage) stepGray(current *gocv.Mat) (*gocv.Mat, error) {
	if current == nil {
		return nil, nil
	}
	dst, err := s.Gray(*current)
	if err != nil {
		return nil, err
	}
	return &dst, nil
}

func (s *Storage) stepDeskew(current *gocv.Mat, step entity.PipelineStep) (*gocv.Mat, error) {
	if current == nil {
		return nil, nil
	}
	dst, _, err := s.DeskewImage(*current)
	if err != nil {
		return nil, err
	}
	return &dst, nil
}

func (s *Storage) stepDewarp(current *gocv.Mat, res entity.PipelineStep) (*gocv.Mat, error) {
	if current == nil || len(res.Points) == 0 {
		return current, nil
	}
	dst := s.DewarpDocument(*current, res.Points)
	return &dst, nil
}

func (s *Storage) stepResize(current *gocv.Mat,step entity.PipelineStep) (*gocv.Mat, error) {
	if current == nil {
		return nil, nil
	}
	buf, err := gocv.IMEncode(".png", *current)
	if err != nil {
		return nil, err
	}
	defer buf.Close()

	oldMax := s.maxDimension
	if step.MaxDim > 0 {
		s.maxDimension = step.MaxDim
	}
	resizedBytes, err := s.Resize(buf.GetBytes())
	if step.MaxDim > 0 {
		s.maxDimension = oldMax
	}
	if err != nil {
		return nil, err
	}

	src, err := gocv.IMDecode(resizedBytes, gocv.IMReadColor)
	if err != nil {
		return nil, err
	}
	return &src, nil
}

func (s *Storage) stepContrast(current *gocv.Mat) (*gocv.Mat, error) {
	if current == nil {
		return nil, nil
	}
	dst, err := s.Normalize(*current)
	if err != nil {
		return nil, err
	}
	return &dst, nil
}

func (s *Storage) stepRemoveLines(current *gocv.Mat) (*gocv.Mat, error) {
	if current == nil {
		return nil, nil
	}
	dst := s.RemoveHorizontalLines(*current)
	return &dst, nil
}

func (s *Storage) stepDenoise(current *gocv.Mat, res entity.PipelineStep) (*gocv.Mat, error) {
	if current == nil {
		return current, nil
	}
	dst, err := s.HybridDenoise(*current, res.Denoise.KernelSize)
	if err != nil {
		return nil, err
	}
	return &dst, nil
}

func (s *Storage) stepFilter(current *gocv.Mat) (*gocv.Mat, error) {
	if current == nil {
		return current, nil
	}
	dst, err := s.ApplyFilter(*current)
	if err != nil {
		return nil, err
	}
	return &dst, nil
}

func (s *Storage) stepThreshold(current *gocv.Mat) (*gocv.Mat, error) {
	if current == nil {
		return nil, nil
	}
	dst := s.AdaptiveThreshold(*current)
	return &dst, nil
}

func (s *Storage) stepColorSpace(current *gocv.Mat) (*gocv.Mat, error) {
	if current == nil {
		return nil, nil
	}
	dst, err := s.ToLAB(*current)
	if err != nil {
		return nil, err
	}
	return &dst, nil
}

func (s *Storage) stepColorTransfer(current *gocv.Mat, step entity.PipelineStep) (*gocv.Mat, error) {
	if current == nil || step.Target.Bytes == nil {
		return current, nil
	}
	target, err := gocv.IMDecode(step.Target.Bytes, gocv.IMReadColor)
	if err != nil {
		return nil, err
	}
	defer target.Close()

	dst, err := s.ColorTransfer(*current, target)
	if err != nil {
		return nil, err
	}
	return &dst, nil
}

func (s *Storage) stepCompress(current *gocv.Mat) ([]byte, error) {
	if current == nil {
		return nil, nil
	}
	img, err := current.ToImage()
	if err != nil {
		return nil, err
	}
	return s.compressor.Encode(img)
}
