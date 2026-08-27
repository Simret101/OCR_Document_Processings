package segmentation

import "gocv.io/x/gocv"

type Segmentations interface {
	OtsuSegmentThreshold(src gocv.Mat) (gocv.Mat, error)
}
