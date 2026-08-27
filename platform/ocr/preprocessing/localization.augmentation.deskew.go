package preprocessing

import (
	"fmt"
	"image"
	"math"

	"gocv.io/x/gocv"
)

func (e *Storage) DeskewImage(img gocv.Mat) (gocv.Mat, float64, error) {
	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)

	thresh := gocv.NewMat()
	defer thresh.Close()
	gocv.Threshold(gray, &thresh, 0, 255, gocv.ThresholdBinaryInv|gocv.ThresholdOtsu)

	angle, err := houghBasedSkew(thresh)
	if err != nil {
		angle, err = contourBasedSkew(thresh)
		if err != nil {
			return img, 0, fmt.Errorf("deskew failed: %w", err)
		}
	}

	result := e.rotateImage(img, -angle)
	return result, angle, nil
}

func houghBasedSkew(binary gocv.Mat) (float64, error) {
	edges := gocv.NewMat()
	defer edges.Close()
	gocv.Canny(binary, &edges, 50, 150)

	lines := gocv.NewMat()
	defer lines.Close()
	gocv.HoughLines(edges, &lines, 1, float32(math.Pi/180), 100)

	if lines.Empty() {
		return 0, fmt.Errorf("no hough lines detected")
	}

	data, err := lines.DataPtrFloat32()
	if err != nil {
		return 0, fmt.Errorf("failed reading lines: %w", err)
	}

	var angleSum float64
	var count int
	for i := 0; i < len(data); i += 2 {
		theta := float64(data[i+1])
		angle := theta*180/math.Pi - 90
		if angle > -45 && angle < 45 {
			angleSum += angle
			count++
		}
	}

	if count == 0 {
		return 0, fmt.Errorf("no valid angles in hough lines")
	}
	return angleSum / float64(count), nil
}

func contourBasedSkew(binary gocv.Mat) (float64, error) {
	contours := gocv.FindContours(binary, gocv.RetrievalExternal, gocv.ChainApproxSimple)
	if contours.Size() == 0 {
		return 0, fmt.Errorf("no contours found")
	}

	for i := 0; i < contours.Size(); i++ {
		c := contours.At(i)
		rect := gocv.MinAreaRect(c)
		angle := rect.Angle
		if angle < -45 {
			angle = -(90 + angle)
		} else {
			angle = -angle
		}
		if math.Abs(angle) < 45 {
			return angle, nil
		}
	}
	return 0, fmt.Errorf("could not determine skew from contours")
}

func (e *Storage) rotateImage(img gocv.Mat, angle float64) gocv.Mat {
	rad := angle * math.Pi / 180
	cos := math.Abs(math.Cos(rad))
	sin := math.Abs(math.Sin(rad))
	newW := int(float64(img.Cols())*cos + float64(img.Rows())*sin)
	newH := int(float64(img.Cols())*sin + float64(img.Rows())*cos)

	center := image.Point{X: img.Cols() / 2, Y: img.Rows() / 2}

	rotationMatrix := gocv.GetRotationMatrix2D(center, angle, 1.0)
	defer rotationMatrix.Close()

	rotationMatrix.SetDoubleAt(0, 2, rotationMatrix.GetDoubleAt(0, 2)+float64(newW)/2-float64(img.Cols())/2)
	rotationMatrix.SetDoubleAt(1, 2, rotationMatrix.GetDoubleAt(1, 2)+float64(newH)/2-float64(img.Rows())/2)

	result := gocv.NewMat()
	gocv.WarpAffine(img, &result, rotationMatrix, image.Point{X: newW, Y: newH})
	return result
}
