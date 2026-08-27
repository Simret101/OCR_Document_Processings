package postprocessing

import (
	"math"
	"strings"
)

func (p *PostProcessor) hasNumericContext(text, substr string) bool {
	pos := 0
	for {
		idx := strings.Index(text[pos:], substr)
		if idx < 0 {
			break
		}
		absPos := pos + idx
		start := int(math.Max(0, float64(absPos-3)))
		end := int(math.Min(float64(len(text)), float64(absPos+len(substr)+3)))
		if p.numericRe.MatchString(text[start:end]) {
			return true
		}
		pos = absPos + 1
	}
	return false
}
