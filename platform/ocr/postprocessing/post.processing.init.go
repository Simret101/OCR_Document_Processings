package postprocessing

import (
	"aidoc/internal/entity"
	"regexp"
	"strings"
)

type PostProcessor struct {
	config     entity.PostProcessConfig
	rules      []compiledRule
	headerRe   []*regexp.Regexp
	footerRe   []*regexp.Regexp
	numericRe  *regexp.Regexp
	dictLookup map[string]bool
}

type compiledRule struct {
	pattern     *regexp.Regexp
	replacement string
	description string
}

func NewPostProcessor(cfg entity.PostProcessConfig) *PostProcessor {
	rules := make([]compiledRule, 0, len(cfg.CorrectionRules))
	for _, r := range cfg.CorrectionRules {
		re, err := regexp.Compile(r.Pattern)
		if err != nil {
			continue
		}
		rules = append(rules, compiledRule{
			pattern: re, 
			replacement: r.Replacement,
			description: r.Description,
		})
	}

	headerRe := make([]*regexp.Regexp, len(cfg.HeaderPatterns))
	for i, hp := range cfg.HeaderPatterns {
		headerRe[i] = regexp.MustCompile(hp)
	}
	footerRe := make([]*regexp.Regexp, len(cfg.FooterPatterns))
	for i, fp := range cfg.FooterPatterns {
		footerRe[i] = regexp.MustCompile(fp)
	}

	dictLookup := make(map[string]bool, len(cfg.Dictionary))
	for _, w := range cfg.Dictionary {
		dictLookup[strings.ToLower(w)] = true
	}

	return &PostProcessor{
		config:     cfg,
		rules:      rules,
		headerRe:   headerRe,
		footerRe:   footerRe,
		numericRe:  regexp.MustCompile(`\d`),
		dictLookup: dictLookup,
	}
}
