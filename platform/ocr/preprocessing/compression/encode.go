package compression

import (
	"bytes"
	"aidoc/internal/entity"
	"encoding/binary"
	"image"
)

func (c *Compression) Encode(img image.Image) ([]byte, error) {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	pixels := make([]int32, 0, w*h*3)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			co, cg, yc := rgbToYCoCg(uint8(r>>8), uint8(g>>8), uint8(b>>8))
			pixels = append(pixels, co, cg, yc)
		}
	}

	residuals := computeResiduals(pixels, w, h)

	freq := make(map[int32]int)
	for _, v := range residuals {
		freq[v]++
	}
	tree := buildHuffmanTree(freq)

	var buf bytes.Buffer
	write32 := func(v int32) { binary.Write(&buf, binary.BigEndian, v) }
	write32(int32(w))
	write32(int32(h))

	if tree == nil {
		write32(0)
		for _, v := range residuals {
			write32(v)
		}
		return buf.Bytes(), nil
	}

	codes := make(map[int32]string)
	generateCodes(tree, "", codes)
	write32(int32(len(codes)))
	for sym, code := range codes {
		write32(sym)
		write32(int32(len(code)))
		codeBits := 0
		for _, c := range code {
			codeBits = (codeBits << 1)
			if c == '1' {
				codeBits |= 1
			}
		}
		write32(int32(codeBits))
	}

	var bw entity.BitWriter
	for _, v := range residuals {
		c.writeCode(codes[v])
	}
	c.flush()
	encoded := bw.Buf.Bytes()
	write32(int32(len(encoded)))
	buf.Write(encoded)
	return buf.Bytes(), nil
}

