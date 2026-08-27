package postprocessing

import (
	"strings"

	"aidoc/internal/entity"
)

func (p *PostProcessor) correctCharacters(text string, doc *entity.OCRDocument) string {
	result := text

	for _, cc := range p.config.CharCorrections {
		if doc != nil && p.hasNumericContext(result, cc.Wrong) {
			continue
		}
		result = strings.ReplaceAll(result, cc.Wrong, cc.Correct)
	}

	result = p.applyContextualCorrections(result)
	return result
}
