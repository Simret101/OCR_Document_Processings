package entity

type CorrectionRule struct {
	Pattern     string
	Replacement string
	Description string
}

type CharCorrection struct {
	Wrong       string
	Correct     string
	Description string
}

type PostProcessConfig struct {
	NormalizeWhitespace   bool
	ReconstructParagraphs bool
	RemoveHeadersFooters  bool
	HeaderPatterns        []string
	FooterPatterns        []string
	CharCorrections       []CharCorrection
	CorrectionRules       []CorrectionRule
	MinLineConfidence     float64
	MaxLineLength         int
	MergeHyphenatedWords  bool
	FixEllipsis           bool
	NormalizeQuotes       bool
	StripNonASCII         bool
	MaxConsecutiveBlanks  int
	EnableErrorDetection  bool
	EnableCandidateSelection bool
	Dictionary            []string
	MaxEditDistance       int
}

type PostProcessStep string

const (
	PostStepCharCorrection       PostProcessStep = "char_correction"
	PostStepWhitespaceNormalize  PostProcessStep = "whitespace_normalize"
	PostStepParagraphReconstruct PostProcessStep = "paragraph_reconstruct"
	PostStepHeaderFooterRemove   PostProcessStep = "header_footer_remove"
	PostStepEncodingNormalize    PostProcessStep = "encoding_normalize"
	PostStepRegexCorrection      PostProcessStep = "regex_correction"
	PostStepConfidenceFilter     PostProcessStep = "confidence_filter"
	PostStepErrorDetection       PostProcessStep = "error_detection"
	PostStepCandidateSelection   PostProcessStep = "candidate_selection"
)

func DefaultCharCorrections() []CharCorrection {
	return []CharCorrection{
		{Wrong: "0", Correct: "O", Description: "zero vs capital O"},
		{Wrong: "O", Correct: "0", Description: "capital O vs zero in numeric context"},
		{Wrong: "1", Correct: "l", Description: "one vs lowercase l"},
		{Wrong: "l", Correct: "1", Description: "lowercase l vs one in numeric context"},
		{Wrong: "|", Correct: "I", Description: "pipe vs capital I"},
		{Wrong: "rn", Correct: "m", Description: "rn ligature"},
		{Wrong: "cl", Correct: "d", Description: "cl ligature"},
		{Wrong: "5", Correct: "S", Description: "five vs capital S"},
		{Wrong: "§", Correct: "S", Description: "section symbol vs S"},
	}
}

func DefaultPostProcessConfig() PostProcessConfig {
	return PostProcessConfig{
		NormalizeWhitespace:     true,
		ReconstructParagraphs:   true,
		RemoveHeadersFooters:    true,
		MergeHyphenatedWords:    true,
		FixEllipsis:             true,
		NormalizeQuotes:         true,
		StripNonASCII:           false,
		MaxConsecutiveBlanks:    2,
		MaxLineLength:           200,
		MinLineConfidence:       0.3,
		EnableErrorDetection:    true,
		EnableCandidateSelection: true,
		MaxEditDistance:         2,
		CharCorrections:         DefaultCharCorrections(),
		HeaderPatterns: []string{
			`^Page\s+\d+\s*(?:of\s+\d+)?$`,
			`^\d+\s+of\s+\d+\s*$`,
			`^Document\s+(?:ID|Number|#)[:\s]*.+$`,
			`^Confidential`,
			`^DRAFT`,
		},
		FooterPatterns: []string{
			`Page\s+\d+\s*(?:of\s+\d+)?\s*$`,
			`\d+\s+of\s+\d+\s*$`,
			`Copyright\s+©?\s*\d{4}`,
			`All\s+rights\s+reserved`,
		},
		CorrectionRules: []CorrectionRule{
			{Pattern: `[♪♫♬]`, Replacement: "", Description: "remove music symbols"},
			{Pattern: `[●■◆▲▼○□▷◁]`, Replacement: "", Description: "remove geometric shapes"},
			{Pattern: `[•·]`, Replacement: " ", Description: "bullet to space"},
			{Pattern: `[           ​]`, Replacement: " ", Description: "unicode whitespace to space"},
			{Pattern: `[‐‑–—]`, Replacement: "-", Description: "hyphens/dashes to standard hyphen"},
			{Pattern: `[‘’ʻʼ]`, Replacement: "'", Description: "curly quotes to straight apostrophe"},
		},
		Dictionary: DefaultDictionary(),
	}
}

func DefaultDictionary() []string {
	return []string{
		"the", "be", "to", "of", "and", "a", "in", "that", "have", "I",
		"it", "for", "not", "on", "with", "he", "as", "you", "do", "at",
		"this", "but", "his", "by", "from", "they", "we", "say", "her", "she",
		"or", "an", "will", "my", "one", "all", "would", "there", "their", "what",
		"so", "up", "out", "if", "about", "who", "get", "which", "go", "me",
		"when", "make", "can", "like", "time", "no", "just", "him", "know", "take",
		"people", "into", "year", "your", "good", "some", "could", "them", "see", "other",
		"than", "then", "now", "look", "only", "come", "its", "over", "think", "also",
		"back", "after", "use", "two", "how", "our", "work", "first", "well", "way",
		"even", "new", "want", "because", "any", "these", "give", "day", "most", "us",
		"invoice", "amount", "total", "date", "number", "name", "address", "phone", "email",
		"payment", "balance", "due", "tax", "price", "quantity", "description", "unit",
		"customer", "supplier", "vendor", "purchase", "order", "receipt", "account",
		"service", "product", "code", "rate", "charge", "fee", "discount", "subtotal",
	}
}
