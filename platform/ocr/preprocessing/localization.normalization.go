package preprocessing

import (
	"gocv.io/x/gocv"
)

func (p *Storage) Normalize(src gocv.Mat) (gocv.Mat, error) {
	srcGray := src
	ownsGray := false
	if src.Channels() > 1 {
		srcGray = gocv.NewMat()
		ownsGray = true
		gocv.CvtColor(src, &srcGray, gocv.ColorBGRToGray)
	}

	srcFloat := gocv.NewMat()
	srcGray.ConvertTo(&srcFloat, gocv.MatTypeCV32F)

	if ownsGray {
		srcGray.Close()
	}

	mean := gocv.NewMat()
	stddev := gocv.NewMat()
	gocv.MeanStdDev(srcFloat, &mean, &stddev)

	meanVal := mean.GetFloatAt(0, 0)
	stddevVal := stddev.GetFloatAt(0, 0)

	mean.Close()
	stddev.Close()

	if stddevVal == 0 {
		stddevVal = 1
	}

	
	zscore := gocv.NewMat()
	gocv.AddWeighted(srcFloat, 1.0/float64(stddevVal), srcFloat, 0, -float64(meanVal)/float64(stddevVal), &zscore)

	srcFloat.Close()

	gocv.Normalize(zscore, &zscore, 0, 255, gocv.NormMinMax)

	dst := gocv.NewMat()
	zscore.ConvertTo(&dst, gocv.MatTypeCV8U)
	zscore.Close()

	return dst, nil
}

