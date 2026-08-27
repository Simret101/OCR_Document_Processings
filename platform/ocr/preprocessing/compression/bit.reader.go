package compression

func (r *Compression) readBit() (int, bool) {
	if r.bitreader.Avail == 0 {
		if r.bitreader.Pos >= len(r.bitreader.Data) {
			return 0, false
		}
		r.bitreader.Bits = r.bitreader.Data[r.bitreader.Pos]
		r.bitreader.Pos++
		r.bitreader.Avail = 8
	}
	r.bitreader.Avail--
	return int((r.bitreader.Bits >> r.bitreader.Avail) & 1), true
}
