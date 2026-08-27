package preprocessing

import (
	"image"

	"gocv.io/x/gocv"
)


func (p *Storage) Resize(data []byte) ([]byte, error) {

	src, err := gocv.IMDecode(data, gocv.IMReadColor)
	if err != nil {
		return nil, err
	}
	defer src.Close()

	srcWidth := src.Cols()
	srcHeight := src.Rows()

	if srcWidth <= p.maxDimension && srcHeight <= p.maxDimension {
		buf, err := gocv.IMEncode(".png", src)
		if err != nil {
			return nil, err
		}
		return buf.GetBytes(), nil
	}

	newWidth, newHeight := calculateSize(srcWidth,srcHeight,p.maxDimension)

	dst := gocv.NewMat()
	defer dst.Close()

	gocv.Resize(src,&dst,image.Pt(newWidth, newHeight),0,0,chooseInterpolation(srcWidth,srcHeight,newWidth,newHeight))

	buf, err := gocv.IMEncode(".png", dst)
	if err != nil {
		return nil, err
	}

	return buf.GetBytes(), nil
}

func calculateSize(width,height,maxDimension int) (int, int) {

	if width >= height {
		scale := float64(maxDimension) / float64(width)
		return maxDimension, int(float64(height) * scale)
	}

	scale := float64(maxDimension) / float64(height)
	return int(float64(width) * scale), maxDimension
}

func chooseInterpolation(srcWidth,srcHeight,dstWidth,dstHeight int) gocv.InterpolationFlags {

	scaleX := float64(dstWidth) / float64(srcWidth)
	scaleY := float64(dstHeight) / float64(srcHeight)

	scale := scaleX
	if scaleY > scale {
		scale = scaleY
	}

	if scale < 1.0 {
		return gocv.InterpolationArea
	}

	return gocv.InterpolationCubic
}