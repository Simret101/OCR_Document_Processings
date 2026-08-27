package preprocessing

import (
	"errors"

	"gocv.io/x/gocv"
)

var errDivisionByZero = errors.New("color transfer: source channel has zero standard deviation")

func (c *Storage) ToHSV(src gocv.Mat) (gocv.Mat, error) {

	dst := gocv.NewMat()

	gocv.CvtColor(src, &dst, gocv.ColorBGRToHSV)

	return dst, nil
}

func (c *Storage) ToLAB(src gocv.Mat) (gocv.Mat, error) {

	dst := gocv.NewMat()

	gocv.CvtColor(src, &dst, gocv.ColorBGRToLab)

	return dst, nil
}

func (c *Storage) GrayToBGR(src gocv.Mat) (gocv.Mat, error) {

	dst := gocv.NewMat()

	gocv.CvtColor(src, &dst, gocv.ColorGrayToBGR)

	return dst, nil
}

func (c *Storage) ColorTransfer(source, target gocv.Mat) (gocv.Mat, error) {

	srcLAB := gocv.NewMat()
	gocv.CvtColor(source, &srcLAB, gocv.ColorBGRToLab)
	srcLAB.ConvertTo(&srcLAB, gocv.MatTypeCV32F)
	defer srcLAB.Close()

	tarLAB := gocv.NewMat()
	gocv.CvtColor(target, &tarLAB, gocv.ColorBGRToLab)
	tarLAB.ConvertTo(&tarLAB, gocv.MatTypeCV32F)
	defer tarLAB.Close()

	srcCh := gocv.Split(srcLAB)
	tarCh := gocv.Split(tarLAB)

	channels := make([]gocv.Mat, 3)
	for i := 0; i < 3; i++ {
		lSrc, stdSrc := imageStats(srcCh[i])
		lTar, stdTar := imageStats(tarCh[i])

		if stdSrc.Val1 == 0 {
			for j := 0; j < 3; j++ {
				srcCh[j].Close()
				tarCh[j].Close()
				if j < i {
					channels[j].Close()
				}
			}
			return gocv.NewMat(), errDivisionByZero
		}

		ratio := stdTar.Val1 / stdSrc.Val1

		gocv.AddWeighted(tarCh[i], ratio, tarCh[i], 0, lSrc.Val1-lTar.Val1*ratio, &channels[i])

		srcCh[i].Close()
		tarCh[i].Close()
	}

	dst := gocv.NewMat()
	gocv.Merge(channels, &dst)
	for _, m := range channels {
		m.Close()
	}

	dst.ConvertTo(&dst, gocv.MatTypeCV8U)

	result := gocv.NewMat()
	gocv.CvtColor(dst, &result, gocv.ColorLabToBGR)
	dst.Close()

	return result, nil
}

func imageStats(mat gocv.Mat) (gocv.Scalar, gocv.Scalar) {
	meanMat := gocv.NewMat()
	stdMat := gocv.NewMat()
	defer meanMat.Close()
	defer stdMat.Close()

	gocv.MeanStdDev(mat, &meanMat, &stdMat)

	mean := meanMat.GetDoubleAt(0, 0)
	std := stdMat.GetDoubleAt(0, 0)

	return gocv.NewScalar(mean, mean, mean, 0),
		gocv.NewScalar(std, std, std, 0)
}
