package postprocessing

import (
	"aidoc/internal/entity"
	"strings"
)

func (p *PostProcessor) filterByConfidence(text string, doc *entity.OCRDocument) string {
	if doc == nil || doc.Metadata == nil {
		return text
	}

	confMap, ok := doc.Metadata["line_confidences"].(map[string]float64)
	if !ok {
		return text
	}

	lines := strings.Split(text, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			result = append(result, line)
			continue
		}

		conf, exists := confMap[trimmed]
		if !exists || conf >= p.config.MinLineConfidence {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}
