package postprocessing

import (
	"strings"
)

func (p *PostProcessor) reconstructParagraphs(text string) string {
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		return text
	}

	var result []string

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			result = append(result, "")
			continue
		}

		if p.config.MergeHyphenatedWords && strings.HasSuffix(line, "-") {
			stem := strings.TrimSuffix(line, "-")
			if i+1 < len(lines) {
				nextLine := strings.TrimSpace(lines[i+1])
				result = append(result, stem+nextLine)
				i++
				continue
			}
		}

		if len(result) > 0 && len(line) < p.config.MaxLineLength && !isSentenceEnd(result[len(result)-1]) {
			result[len(result)-1] += " " + line
		} else {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

func isSentenceEnd(s string) bool {
	if len(s) == 0 {
		return true
	}
	last := s[len(s)-1]
	return last == '.' || last == '!' || last == '?' || last == ':'
}
