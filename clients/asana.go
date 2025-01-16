package clients

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"go.uber.org/zap"

	"github.com/cyber/test-project/logging"
	"github.com/cyber/test-project/models"
)

type AsanaClient struct {
	baseClient *httpClient
}

const (
	getUsersEndpoint    = "/api/1.0/users"
	getProjectsEndpoint = "/api/1.0/projects"
)

func NewAsanaClient(options ClientOptions) *AsanaClient {
	return &AsanaClient{
		baseClient: newHttpClient(options),
	}
}

type GetUsersRequest struct {
	Workspace string
	Team      string
	Limit     int
	Offset    string
	Token     string
}

func (a AsanaClient) GetUsers(ctx context.Context, request GetUsersRequest) (models.AsanaGetUsersResponse, error) {
	logger := logging.FromContext(ctx)
	logger = logger.With(zap.String("operation_name", "asana_get_users"))
	ctx = logging.WithLogger(ctx, logger)

	query := url.Values{}
	if request.Workspace != "" {
		query.Set("workspace", request.Workspace)
	}

	if request.Team != "" {
		query.Set("team", request.Team)
	}
	if request.Limit != 0 {
		query.Set("limit", strconv.Itoa(request.Limit))
	}
	if request.Offset != "" {
		query.Set("offset", request.Offset)
	}

	req := httpRequest{
		method: http.MethodGet,
		path:   getUsersEndpoint,
		query:  query,
		headers: map[string]string{
			"Accept":        "application/json",
			"Authorization": "Bearer " + request.Token,
		},
	}

	resp, err := a.baseClient.doRequest(ctx, req)
	// @TODO: process 429
	if err != nil {
		return models.AsanaGetUsersResponse{}, models.ErrServiceFailure{ServiceName: "asana"}
	}

	var response models.AsanaGetUsersResponse
	err = json.Unmarshal(resp, &response)
	if err != nil {
		logger.Error("Failed to unmarshal get_users response", zap.Error(err))
		return models.AsanaGetUsersResponse{}, models.ErrServiceFailure{ServiceName: "asana"}
	}

	return response, nil
}

type GetProjectsRequest struct {
	Workspace string
	Team      string
	Limit     int
	Offset    string
	Token     string
	Archived  *bool
}

func (a AsanaClient) GetProjects(ctx context.Context, request GetProjectsRequest) (models.AsanaGetProjectsResponse, error) {
	logger := logging.FromContext(ctx)
	logger = logger.With(zap.String("operation_name", "asana_get_projects"))
	ctx = logging.WithLogger(ctx, logger)

	query := url.Values{}
	if request.Workspace != "" {
		query.Set("workspace", request.Workspace)
	}

	if request.Team != "" {
		query.Set("team", request.Team)
	}
	if request.Limit != 0 {
		query.Set("limit", strconv.Itoa(request.Limit))
	}
	if request.Offset != "" {
		query.Set("offset", request.Offset)
	}
	if request.Archived != nil {
		query.Set("archived", strconv.FormatBool(*request.Archived))
	}

	req := httpRequest{
		method: http.MethodGet,
		path:   getProjectsEndpoint,
		query:  query,
		headers: map[string]string{
			"Accept":        "application/json",
			"Authorization": "Bearer " + request.Token,
		},
	}

	resp, err := a.baseClient.doRequest(ctx, req)
	// @TODO: process 429
	if err != nil {
		return models.AsanaGetProjectsResponse{}, models.ErrServiceFailure{ServiceName: "asana"}
	}

	var response models.AsanaGetProjectsResponse
	err = json.Unmarshal(resp, &response)
	if err != nil {
		logger.Error("Failed to unmarshal get_users response", zap.Error(err))
		return models.AsanaGetProjectsResponse{}, models.ErrServiceFailure{ServiceName: "asana"}
	}

	return response, nil
}
