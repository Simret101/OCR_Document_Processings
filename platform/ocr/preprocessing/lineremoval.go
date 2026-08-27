package preprocessing

import (
	"image"

	"gocv.io/x/gocv"
)

func (e *Storage) RemoveHorizontalLines(img gocv.Mat) gocv.Mat {

	gray := gocv.NewMat()
	defer gray.Close()

	gocv.CvtColor(img,&gray,gocv.ColorBGRToGray)


	thresh := gocv.NewMat()
	defer thresh.Close()

	gocv.Threshold(gray,&thresh,0,255,gocv.ThresholdBinaryInv|gocv.ThresholdOtsu)



	kernelWidth := img.Cols() / 10

	if kernelWidth < 25 {
		kernelWidth = 25
	}


	horizontalKernel := gocv.GetStructuringElement(gocv.MorphRect,image.Pt(kernelWidth, 1))

	defer horizontalKernel.Close()



	detectedLines := gocv.NewMat()
	defer detectedLines.Close()


	gocv.MorphologyEx(thresh,&detectedLines,gocv.MorphOpen,horizontalKernel)



	result := img.Clone()



	white := gocv.NewMatWithSize(result.Rows(),result.Cols(),result.Type())

	defer white.Close()



	white.SetTo(gocv.NewScalar(255,255,255,0))

	white.CopyToWithMask(&result,detectedLines)

	repairKernel := gocv.GetStructuringElement(gocv.MorphRect,image.Pt(1,6))

	defer repairKernel.Close()



	inverted := gocv.NewMat()
	defer inverted.Close()


	gocv.BitwiseNot(result,&inverted)



	repaired := gocv.NewMat()
	defer repaired.Close()



	gocv.MorphologyEx(inverted,&repaired,gocv.MorphClose,repairKernel)



	final := gocv.NewMat()


	gocv.BitwiseNot(repaired,&final)



	return final
}