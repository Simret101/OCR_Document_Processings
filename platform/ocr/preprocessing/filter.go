package preprocessing

import (
	"image"

	"gocv.io/x/gocv"
)

func (p *Storage) ApplyFilter(src gocv.Mat) (gocv.Mat, error) {

	dst := gocv.NewMat()

	ksize := p.KernelSize
	if ksize < 3 {
		ksize = 3
	}
	if ksize%2 == 0 {
		ksize++
	}

	blurred := gocv.NewMat()
	defer blurred.Close()
	gocv.GaussianBlur(src, &blurred, image.Pt(ksize, ksize), 2.0, 2.0, gocv.BorderDefault)

	gocv.AddWeighted(src, 1.5, blurred, -0.5, 0, &dst)

	return dst, nil
}
