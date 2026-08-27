package async

import (
	"aidoc/internal/entity"
	"aidoc/internal/errors"
	"aidoc/platform/ocr/amazontextract"
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/google/uuid"
)

func (r *Reader) GetAnalysisResult(ctx context.Context, jobID string, features []entity.TextractFeature, timeout time.Duration) (*entity.OCRDocument, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var blocks []types.Block
	nextToken := ""

	for {
		params := &textract.GetDocumentAnalysisInput{
			JobId: aws.String(jobID),
		}
		if nextToken != "" {
			params.NextToken = aws.String(nextToken)
		}

		resp, err := r.textractClient.GetDocumentAnalysis(ctx, params)
		if err != nil {
			return nil, err
		}

		if resp.JobStatus == types.JobStatusFailed {
			return nil, errors.ErrTextractFailed.New("Error: job %s failed", jobID)
		}

		if resp.JobStatus == types.JobStatusInProgress {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(2 * time.Second):
				continue
			}
		}

		blocks = append(blocks, resp.Blocks...)

		if resp.NextToken == nil || *resp.NextToken == "" {
			break
		}
		nextToken = *resp.NextToken
	}

	doc := &entity.OCRDocument{
		ID:            uuid.NewString(),
		ExtractedData: make(map[string]interface{}),
		ProcessedAt:   time.Now(),
	}
	amazontextract.PopulateDocument(doc, blocks)
	amazontextract.ExtractStructuredData(doc, blocks, features)
	return doc, nil
}
