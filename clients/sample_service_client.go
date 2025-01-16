package clients

import (
	"context"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/cyber/test-project/logging"
	"github.com/cyber/test-project/models"
)

type SampleServiceClient struct {
	client *httpClient
}

func NewSampleServiceClient(options ClientOptions) *SampleServiceClient {
	return &SampleServiceClient{
		client: newHttpClient(options),
	}
}

type SomeActionRequest struct {
	Field1 string `json:"field1"`
	Field2 string `json:"field2"`
}

func (s SampleServiceClient) SomeAction(ctx context.Context, request SomeActionRequest) (*models.SampleResponse, error) {
	logger := logging.FromContext(ctx).With(zap.String("operation", "some_action"))
	ctx = logging.WithLogger(ctx, logger)

	req := httpRequest{
		method: http.MethodGet,
		path:   "/some/action",
		body:   request,
		headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer ",
		},
	}
	resp, err := s.client.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	var someResponse models.SampleResponse
	err = json.Unmarshal(resp, &someResponse)
	if err != nil {
		logger.Error("Failed to unmarshal response", zap.Error(err))
		return nil, err
	}

	return &someResponse, nil
}
