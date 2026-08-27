package preprocessing

import "gocv.io/x/gocv"


func (e *Storage)AdaptiveThreshold(img gocv.Mat) gocv.Mat {


	gray := gocv.NewMat()
	defer gray.Close()


	gocv.CvtColor(img,&gray,gocv.ColorBGRToGray)


	result := gocv.NewMat()


	gocv.AdaptiveThreshold(gray,&result,255,gocv.AdaptiveThresholdGaussian,gocv.ThresholdBinary,31,15)


	return result
}