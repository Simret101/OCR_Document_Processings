package compression

func (w *Compression) writeBit(b int) {
	if b != 0 {
		w.bitWriter.Bits |= 1 << (7 - w.bitWriter.Pos)
	}
	w.bitWriter.Pos++
	if w.bitWriter.Pos == 8 {
		w.bitWriter.Buf.WriteByte(w.bitWriter.Bits)
		w.bitWriter.Bits = 0
		w.bitWriter.Pos = 0
	}
}

func (w *Compression) writeCode(code string) {
	for _, c := range code {
		if c == '1' {
			w.writeBit(1)
		} else {
			w.writeBit(0)
		}
	}
}

func (w *Compression) flush() {
	if w.bitWriter.Pos > 0 {
		w.bitWriter.Buf.WriteByte(w.bitWriter.Bits)
	}
}
