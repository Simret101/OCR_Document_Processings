package amazontextract

import (
	"fmt"
	"strings"

	"aidoc/internal/entity"

	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

func PopulateDocument(doc *entity.OCRDocument, blocks []types.Block) {
	var sb strings.Builder
	var pages int

	for _, block := range blocks {
		if block.BlockType == types.BlockTypePage {
			pages++
		}
		if block.BlockType == types.BlockTypeLine && block.Text != nil {
			if sb.Len() > 0 {
				sb.WriteString("\n")
			}
			sb.WriteString(*block.Text)
		}
	}

	doc.RawText = sb.String()
	doc.PageCount = pages

	if doc.RawText == "" {
		doc.Confidence = 0
	} else {
		var total float64
		var count int
		for _, block := range blocks {
			if block.BlockType == types.BlockTypeLine && block.Confidence != nil {
				total += float64(*block.Confidence)
				count++
			}
		}
		if count > 0 {
			doc.Confidence = total / float64(count)
		}
	}
}

func PopulateExpenseDocument(doc *entity.OCRDocument, resp *textract.AnalyzeExpenseOutput) {
	expenseData := make([]map[string]interface{}, 0)

	for _, expenseDoc := range resp.ExpenseDocuments {
		docData := make(map[string]interface{})

		for _, summaryField := range expenseDoc.SummaryFields {
			if summaryField.Type == nil || summaryField.Type.Text == nil {
				continue
			}
			key := *summaryField.Type.Text
			val := ""
			if summaryField.ValueDetection != nil && summaryField.ValueDetection.Text != nil {
				val = *summaryField.ValueDetection.Text
			}
			if val != "" {
				docData[strings.ToLower(key)] = val
			}
		}

		var lineItems []map[string]interface{}
		for _, group := range expenseDoc.LineItemGroups {
			for _, lineItem := range group.LineItems {
				item := make(map[string]interface{})
				for _, prop := range lineItem.LineItemExpenseFields {
					if prop.Type == nil || prop.Type.Text == nil {
						continue
					}
					key := *prop.Type.Text
					val := ""
					if prop.ValueDetection != nil && prop.ValueDetection.Text != nil {
						val = *prop.ValueDetection.Text
					}
					if val != "" {
						item[strings.ToLower(key)] = val
					}
				}
				if len(item) > 0 {
					lineItems = append(lineItems, item)
				}
			}
		}

		if len(lineItems) > 0 {
			docData["line_items"] = lineItems
		}
		expenseData = append(expenseData, docData)

		if id, ok := docData["invoice_receipt_id"]; ok {
			doc.ExtractedData["receipt_id"] = id
		}
		if total, ok := docData["total"]; ok {
			doc.ExtractedData["total"] = total
		}
		if date, ok := docData["invoice_receipt_date"]; ok {
			doc.ExtractedData["receipt_date"] = date
		}
		if due, ok := docData["due_date"]; ok {
			doc.ExtractedData["due_date"] = due
		}
		if vendor, ok := docData["vendor_name"]; ok {
			doc.ExtractedData["vendor_name"] = vendor
		}
	}

	doc.ExtractedData["expense_details"] = expenseData

	if rawText, ok := doc.ExtractedData["receipt_id"]; ok {
		doc.RawText = fmt.Sprintf("%v", rawText)
	}
}

