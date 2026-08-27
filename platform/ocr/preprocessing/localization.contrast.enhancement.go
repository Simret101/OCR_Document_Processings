package preprocessing

import "gocv.io/x/gocv"



func (e *Storage)EnhanceContrast(img gocv.Mat) gocv.Mat {


	gray := gocv.NewMat()
	defer gray.Close()


	gocv.CvtColor(img,&gray,gocv.ColorBGRToGray,)


	clahe := gocv.NewCLAHEWithParams(e.clipLimit,e.tileGridSize)

	defer clahe.Close()


	result := gocv.NewMat()


	clahe.Apply(gray,&result)


	return result
}