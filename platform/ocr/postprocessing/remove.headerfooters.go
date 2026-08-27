package postprocessing

import (
	"strings"
)

func (p *PostProcessor) removeHeadersFooters(text string) string {
	lines := strings.Split(text, "\n")
	if len(lines) <= 3 {
		return text
	}

	var result []string

	for i, line := range lines {
		isHeader := false
		for _, re := range p.headerRe {
			if re.MatchString(strings.TrimSpace(line)) {
				isHeader = true
				break
			}
		}
		if isHeader && i <= 3 {
			continue
		}

		isFooter := false
		for _, re := range p.footerRe {
			if re.MatchString(strings.TrimSpace(line)) {
				isFooter = true
				break
			}
		}
		if isFooter && i >= len(lines)-3 {
			continue
		}

		result = append(result, line)
	}

	return strings.Join(result, "\n")
}
