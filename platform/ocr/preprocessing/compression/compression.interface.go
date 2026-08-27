package compression

import "image"

type Compressor interface {
	Encode(img image.Image) ([]byte, error)
}
