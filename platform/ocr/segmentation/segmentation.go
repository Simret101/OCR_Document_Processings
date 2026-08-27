package segmentation

import (
	"gocv.io/x/gocv"
)

func (s *Segmentation) OtsuSegmentThreshold(src gocv.Mat) (gocv.Mat, error) {
	dst := gocv.NewMat()
	maxVal := float32(s.maxVal)
	if maxVal == 0 {
		maxVal = 255
	}
	gocv.Threshold(src, &dst, 0, maxVal, gocv.ThresholdBinary|gocv.ThresholdOtsu)
	return dst, nil
}
