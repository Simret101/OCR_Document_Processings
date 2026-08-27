package postprocessing

func (p *PostProcessor) applyRegexRules(text string) string {
	result := text
	for _, rule := range p.rules {
		result = rule.pattern.ReplaceAllString(result, rule.replacement)
	}
	return result
}
