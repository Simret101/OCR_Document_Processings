package amazontextract

import "github.com/aws/aws-sdk-go-v2/service/textract"

type Service struct {
	client *textract.Client
}

func NewTextract(
	client *textract.Client,
) *Service {

	return &Service{
		client: client,
	}
}
