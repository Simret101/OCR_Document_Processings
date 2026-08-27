package postprocessing

import "regexp"

var contextualCorrections = []struct {
	re   *regexp.Regexp
	repl string
}{
	{regexp.MustCompile(`\b1nvoice\s+#`), `Invoice #`},
	{regexp.MustCompile(`\b1nvoice\b`), `Invoice`},
	{regexp.MustCompile(`\b(\d+)l(\d+)\b`), `${1}1${2}`},
	{regexp.MustCompile(`\bO(\d+)\b`), `0${1}`},
	{regexp.MustCompile(`(\d+)O(\d+)\b`), `${1}0${2}`},
	{regexp.MustCompile(`\brn([a-z])`), `m${1}`},
}

func (p *PostProcessor) applyContextualCorrections(text string) string {
	result := text
	for _, c := range contextualCorrections {
		result = c.re.ReplaceAllString(result, c.repl)
	}
	return result
}
