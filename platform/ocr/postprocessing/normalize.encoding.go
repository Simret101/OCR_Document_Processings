package postprocessing

import (
	"regexp"
	"strings"
	"unicode"
)

func (p *PostProcessor) normalizeEncoding(text string) string {
	replacements := map[rune]rune{
		'\u2018': '\'', '\u2019': '\'',
		'\u201C': '"', '\u201D': '"',
		'\u2013': '-', '\u2014': '-',
		'\u00A0': ' ',
		'\u2026': '.',
	}
	result := strings.Map(func(r rune) rune {
		if repl, ok := replacements[r]; ok {
			return repl
		}
		if p.config.StripNonASCII && r > 127 && !unicode.IsLetter(r) {
			return -1
		}
		return r
	}, text)

	if p.config.FixEllipsis {
		result = regexp.MustCompile(`\.{3,}`).ReplaceAllString(result, "...")
	}
	if p.config.NormalizeQuotes {
		result = regexp.MustCompile(`「`).ReplaceAllString(result, `"`)
		result = regexp.MustCompile(`」`).ReplaceAllString(result, `"`)
		result = regexp.MustCompile(`『`).ReplaceAllString(result, `"`)
		result = regexp.MustCompile(`』`).ReplaceAllString(result, `"`)
	}

	return result
}
