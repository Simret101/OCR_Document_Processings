package compression

import (
	"bytes"
	"aidoc/internal/entity"
	"encoding/binary"
	"image"
)

func (c *Compression) Decode(data []byte) (image.Image, error) {
	r := bytes.NewReader(data)

	var w, h, numSymbols int32
	binary.Read(r, binary.BigEndian, &w)
	binary.Read(r, binary.BigEndian, &h)
	binary.Read(r, binary.BigEndian, &numSymbols)

	if numSymbols == 0 {
		residuals := make([]int32, w*h*3)
		for i := range residuals {
			binary.Read(r, binary.BigEndian, &residuals[i])
		}
		return reconstructFromResiduals(residuals, int(w), int(h)), nil
	}

	root := &entity.TrieNode{}
	for i := int32(0); i < numSymbols; i++ {
		var sym, clen, codeBits int32
		binary.Read(r, binary.BigEndian, &sym)
		binary.Read(r, binary.BigEndian, &clen)
		binary.Read(r, binary.BigEndian, &codeBits)
		node := root
		for j := clen - 1; j >= 0; j-- {
			if (codeBits>>uint(j))&1 == 0 {
				if node.Left == nil {
					node.Left = &entity.TrieNode{}
				}
				node = node.Left
			} else {
				if node.Right == nil {
					node.Right = &entity.TrieNode{}
				}
				node = node.Right
			}
		}
		node.Value = &sym
	}

	var encodedLen int32
	binary.Read(r, binary.BigEndian, &encodedLen)
	encoded := make([]byte, encodedLen)
	r.Read(encoded)

	c.bitreader.Data = encoded
	c.bitreader.Pos = 0
	c.bitreader.Bits = 0
	c.bitreader.Avail = 0
	residuals := make([]int32, 0, w*h*3)
	target := int(w) * int(h) * 3
	node := root
	for len(residuals) < target {
		bit, ok := c.readBit()
		if !ok {
			break
		}
		if bit == 0 {
			node = node.Left
		} else {
			node = node.Right
		}
		if node == nil {
			break
		}
		if node.Value != nil {
			residuals = append(residuals, *node.Value)
			node = root 
		}
	}
	return reconstructFromResiduals(residuals, int(w), int(h)), nil
}
