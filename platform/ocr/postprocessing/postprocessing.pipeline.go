package postprocessing

import (
	"context"

	"aidoc/internal/entity"
)

func (p *PostProcessor) PostProcessPipeline(
	ctx context.Context,
	text string,
	doc *entity.OCRDocument,
) (string, error) {

	if err := ctx.Err(); err != nil {
		return "", err
	}

	result := text

	if doc != nil {
		result = p.correctCharacters(result, doc)
	}

	result = p.normalizeEncoding(result)

	if p.config.NormalizeWhitespace {
		result = p.normalizeWhitespace(result)
	}

	if p.config.ReconstructParagraphs {
		result = p.reconstructParagraphs(result)
	}

	if p.config.RemoveHeadersFooters {
		result = p.removeHeadersFooters(result)
	}

	result = p.applyRegexRules(result)

	if doc != nil && p.config.MinLineConfidence > 0 {
		result = p.filterByConfidence(result, doc)
	}

	if err := ctx.Err(); err != nil {
		return "", err
	}

	return result, nil
}
