package upload

type UploadStatementResponse struct {
	DocumentID string  `json:"document_id"`
	Status     string  `json:"status"`
	OCRText    string  `json:"ocr_text,omitempty"`
	Confidence float64 `json:"confidence,omitempty"`
}

