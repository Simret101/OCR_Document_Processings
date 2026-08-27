package async

import (
	"aidoc/internal/entity"
	"aidoc/internal/errors"

	"aidoc/platform/ocr/amazontextract"
	"context"

	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

func (r *Reader) StartDocumentAnalysis(ctx context.Context, input *entity.ReadInput, features []entity.TextractFeature) (*string, error) {
	if input.S3Bucket == "" || input.S3Key == "" {
		return nil, errors.ErrInvalidDocument.New("Error: async analysis requires S3 input")
	}

	params := &textract.StartDocumentAnalysisInput{
		DocumentLocation: &types.DocumentLocation{
			S3Object: &types.S3Object{
				Bucket: &input.S3Bucket,
				Name:   &input.S3Key,
			},
		},
		FeatureTypes: amazontextract.BuildFeatureTypes(features),
	}

	if len(input.Queries) > 0 {
		params.QueriesConfig = amazontextract.BuildQueriesConfig(input.Queries)
	}

	resp, err := r.textractClient.StartDocumentAnalysis(ctx, params)
	if err != nil {
		return nil, err
	}
	return resp.JobId, nil
}
