package postprocessing

import (
	"fmt"
	"regexp"
	"strings"
)

func (p *PostProcessor) normalizeWhitespace(text string) string {
	result := regexp.MustCompile(`[ \t]+`).ReplaceAllString(text, " ")

	pattern := fmt.Sprintf(`\n{%d,}`, p.config.MaxConsecutiveBlanks+1)
	result = regexp.MustCompile(pattern).ReplaceAllString(result, strings.Repeat("\n", p.config.MaxConsecutiveBlanks))

	result = regexp.MustCompile(`[ \t]+\n`).ReplaceAllString(result, "\n")

	result = regexp.MustCompile(`\n[ \t]+`).ReplaceAllString(result, "\n")

	result = strings.TrimSpace(result)
	return result
}
