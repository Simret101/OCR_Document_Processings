package sync

import (
	"aidoc/internal/entity"
	x "errors"
	"fmt"

	"aidoc/platform/ocr/amazontextract"
	"context"

	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/smithy-go"
)

func (r *Reader) AnalyzeExpense(ctx context.Context, input *entity.ReadInput) (*entity.OCRDocument, error) {
	doc, err := amazontextract.CreateDocument(input)
	if err != nil {
		return nil, err
	}
	resp, err := r.analyzeExpense(ctx, input)
	if err != nil {
		fmt.Printf("Textract error: %T\n", err)
		fmt.Printf("Textract error: %+v\n", err)
		fmt.Printf("Textract error string: %s\n", err.Error())

		var opErr *smithy.OperationError
		if x.As(err, &opErr) {
			fmt.Printf("Service: %s\n", opErr.Service())
			fmt.Printf("Operation: %s\n", opErr.Operation())
			fmt.Printf("Underlying error: %T\n", opErr.Unwrap())
			fmt.Printf("Underlying error: %+v\n", opErr.Unwrap())
		}
		doc.Status = "failed"
		doc.ErrorMessage = err.Error()
		amazontextract.SaveDocument(ctx, r.dynamoClient, r.config.TableName, doc)
		return doc, err
	}

	amazontextract.PopulateExpenseDocument(doc, resp)
	doc.Status = "success"
	amazontextract.SaveDocument(ctx, r.dynamoClient, r.config.TableName, doc)
	return doc, nil
}

func (r *Reader) analyzeExpense(ctx context.Context, input *entity.ReadInput) (*textract.AnalyzeExpenseOutput, error) {
	params := &textract.AnalyzeExpenseInput{
		Document: amazontextract.BuildDocument(input),
	}
	return r.textractClient.AnalyzeExpense(ctx, params)
}
