package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/sony/gobreaker"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/cyber/test-project/appcontext"
	"github.com/cyber/test-project/logging"
	"github.com/cyber/test-project/models"
)

var successCodes = map[int]bool{
	http.StatusOK:           true,
	http.StatusAccepted:     true,
	http.StatusNoContent:    true,
	http.StatusResetContent: true,
}

type httpRequest struct {
	method  string
	path    string
	body    any
	query   url.Values
	headers map[string]string
}

func (r httpRequest) toHttpRequest(ctx context.Context, baseUrl string) (*http.Request, error) {
	logger := logging.FromContext(ctx)

	var bodyReader io.Reader
	if r.body != nil {
		bodyBytes, err := json.Marshal(r.body)
		if err != nil {
			logger.Error("Failed to serialize request body", zap.Error(err))
			return nil, err
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(r.method, baseUrl, bodyReader)
	if err != nil {
		logger.Error("Failed to create request", zap.Error(err))
		return nil, err
	}

	req.URL = req.URL.ResolveReference(&url.URL{Path: r.path})

	if len(r.headers) > 0 {
		for key, value := range r.headers {
			req.Header.Add(key, value)
		}
	}

	if len(r.query) > 0 {
		q := req.URL.Query()
		for key, value := range r.query {
			for _, param := range value {
				q.Add(key, param)
			}
		}
	}

	return req.WithContext(ctx), nil
}

type ClientOptions struct {
	ServiceName string
	BaseClient  *http.Client
	BaseURL     string
}

type httpClient struct {
	serviceName    string
	baseClient     *http.Client
	baseUrl        string
	circuitBreaker CircuitBreaker
}

func newHttpClient(options ClientOptions) *httpClient {
	return &httpClient{
		serviceName:    options.ServiceName,
		baseClient:     options.BaseClient,
		baseUrl:        options.BaseURL,
		circuitBreaker: noCircuitBreaker{},
	}
}

func (c httpClient) doRequest(ctx context.Context, req httpRequest) ([]byte, error) {
	logger := logging.FromContext(ctx).
		With(zap.String("method", req.method)).
		With(zap.String("path", req.path)).
		With(zap.String("service", c.serviceName))
	ctx = appcontext.WithLogger(ctx, logger)

	httpReq, err := req.toHttpRequest(ctx, c.baseUrl)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(httpReq)
	if err != nil {
		return nil, models.ErrServiceFailure{ServiceName: c.serviceName}
	}
	defer closeBody(resp, logger)

	logger = logger.With(zap.Int("status_code", resp.StatusCode))

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read response body", zap.Error(err))
		return nil, err
	}

	if successCodes[resp.StatusCode] {
		return respBodyBytes, nil
	}

	return nil, c.handleErrorResponse(ctx, logger, resp.StatusCode, respBodyBytes)
}

func (c httpClient) do(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	logger := logging.FromContext(ctx)

	resp, err := c.circuitBreaker.Execute(func() (any, error) {
		resp, httpErr := c.baseClient.Do(req)
		if httpErr != nil {
			logger.Error("could not perform HTTP request",
				logging.DebugField(func() zapcore.Field {
					reqDump, _ := httputil.DumpRequest(req, true)
					return zap.ByteString("request_dump", reqDump)
				}),
				zap.Error(httpErr),
			)
			return nil, models.ErrServiceFailure{ServiceName: c.serviceName}
		}

		if resp.StatusCode >= http.StatusInternalServerError {
			logger.Debug("got a 5xx HTTP status code from "+c.serviceName, zap.Int("status_code", resp.StatusCode))
			return nil, models.ErrServiceFailure{ServiceName: c.serviceName}
		}

		return resp, nil
	})

	response, _ := resp.(*http.Response)

	if err != nil {
		if errors.Is(err, gobreaker.ErrOpenState) || errors.Is(err, gobreaker.ErrTooManyRequests) {
			return nil, models.ErrRateLimitExceeded{ServiceName: c.serviceName}
		}

		return nil, models.ErrServiceFailure{ServiceName: c.serviceName}
	}

	return response, nil
}

func (c httpClient) handleErrorResponse(ctx context.Context, logger *zap.Logger, statusCode int, respBodyBytes []byte) error {
	switch statusCode {
	case http.StatusTooManyRequests:
		return models.ErrRateLimitExceeded{ServiceName: c.serviceName}
	default:
		return models.ErrServiceFailure{ServiceName: c.serviceName}
	}
}

func closeBody(resp *http.Response, logger *zap.Logger) {
	_, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		logger.Debug("Failed to discard response body", zap.Error(err))
	}

	err = resp.Body.Close()
	if err != nil {
		logger.Debug("Failed to close response body", zap.Error(err))
	}
}
