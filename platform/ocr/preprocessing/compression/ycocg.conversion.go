package compression

func rgbToYCoCg(r, g, b uint8) (int32, int32, int32) {
	ri, gi, bi := int32(r), int32(g), int32(b)
	co := ri - bi         
	tmp := bi + (co >> 1) 
	cg := gi - tmp        
	y := tmp + (cg >> 1)  
	return co, cg, y
}

func yCoCgToRGB(co, cg, y int32) (uint8, uint8, uint8) {
	tmp := y - (cg >> 1) 
	g := cg + tmp        
	b := tmp - (co >> 1) 
	r := co + b        
	return clamp(r), clamp(g), clamp(b)
}

func clamp(v int32) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v)
}
