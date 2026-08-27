package entity

import "bytes"

type BitReader struct {
	Data  []byte
	Pos   int
	Bits  byte
	Avail uint8
}

type BitWriter struct {
	Buf  bytes.Buffer
	Bits byte
	Pos  uint8
}
