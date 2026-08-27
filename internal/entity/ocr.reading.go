package entity

import "time"

type DocumentType string

const (
	DocTypeIDCard           DocumentType = "id_card"
	DocTypeBankStatement    DocumentType = "bank_statement"
	DocTypeEmploymentLetter DocumentType = "employment_letter"
	DocTypePaySlip          DocumentType = "payslip"
	DocTypeInvoice          DocumentType = "invoice"
	DocTypeReceipt          DocumentType = "receipt"
	DocTypeContract         DocumentType = "contract"
	DocTypePassport         DocumentType = "passport"
	DocTypeDriversLicense   DocumentType = "drivers_license"
	DocTypeOther            DocumentType = "other"
)

type TextractFeature string

const (
	TextractText      TextractFeature = "text"
	TextractForms     TextractFeature = "forms"
	TextractTables    TextractFeature = "tables"
	TextractExpense   TextractFeature = "expense"
	TextractSignature TextractFeature = "signature"
	TextractID        TextractFeature = "id"
	TextractLayout    TextractFeature = "layout"
	TextractQueries   TextractFeature = "queries"
)

type OCRDocument struct {
	ID                string                 `json:"id"`
	DocumentType      DocumentType           `json:"document_type"`
	RawText           string                 `json:"raw_text"`
	ExtractedData     map[string]interface{} `json:"extracted_data"`
	Confidence        float64                `json:"confidence"`
	ImagePath         string                 `json:"image_path,omitempty"`
	ProcessedAt       time.Time              `json:"processed_at"`
	Status            string                 `json:"status"`
	ErrorMessage      string                 `json:"error_message,omitempty"`
	LoanApplicationID string                 `json:"loan_application_id,omitempty"`
	UserID            string                 `json:"user_id,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
	PageCount         int                    `json:"page_count,omitempty"`
	Warnings          []string               `json:"warnings,omitempty"`
}

type TextractBlock struct {
	ID         string  `json:"id"`
	BlockType  string  `json:"block_type"`
	Text       string  `json:"text,omitempty"`
	Confidence float64 `json:"confidence"`
	Row        int     `json:"row,omitempty"`
	Column     int     `json:"column,omitempty"`
	Page       int     `json:"page,omitempty"`
}

type TextractResult struct {
	DocumentID string          `json:"document_id"`
	Blocks     []TextractBlock `json:"blocks"`
	RawJSON    string          `json:"raw_json,omitempty"`
	PageCount  int             `json:"page_count"`
}

type TextractConfig struct {
	Region          string
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	TableName       string
	BucketName      string
	KMSKeyID        string
	Timeout         time.Duration
	MaxRetries      int
}
type Query struct {
	Text  string
	Alias string
}


type ReadInput struct {
	Bytes             []byte
	Queries           []Query
	S3Bucket          string
	S3Key             string
	UserID            string
	Metadata          map[string]interface{}
}
