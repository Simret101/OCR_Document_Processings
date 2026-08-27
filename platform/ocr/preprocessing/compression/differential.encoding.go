package compression

import (
	"image"
	"image/color"
)

func computeResiduals(pixels []int32, w, h int) []int32 {
	res := make([]int32, len(pixels))
	copy(res, pixels) 
	stride := w * 3   
	for i := stride; i < len(res); i++ {
		res[i] = pixels[i] - pixels[i-stride] 
	}
	return res
}

func reconstructFromResiduals(residuals []int32, w, h int) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	pixels := make([]int32, len(residuals))
	copy(pixels, residuals)
	stride := w * 3
	for i := stride; i < len(pixels); i++ {
		pixels[i] += pixels[i-stride] 
	}
	idx := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			co := pixels[idx]  
			cg := pixels[idx+1] 
			yc := pixels[idx+2] 
			r, g, b := yCoCgToRGB(co, cg, yc)
			img.SetNRGBA(x, y, color.NRGBA{R: r, G: g, B: b, A: 255})
			idx += 3
		}
	}
	return img
}
