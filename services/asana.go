package services

import (
	"context"

	"github.com/cyber/test-project/clients"
	"github.com/cyber/test-project/models"
)

type AsanaUsersGetter interface {
	GetUsers(context.Context, clients.GetUsersRequest) (models.AsanaGetUsersResponse, error)
}

type AsanaProjectsGetter interface {
	GetProjects(context.Context, clients.GetProjectsRequest) (models.AsanaGetProjectsResponse, error)
}

type AsanaService struct {
	client      *clients.AsanaClient
	accessToken string
	dataDumper  Dumper
}

func NewAsanaService(client *clients.AsanaClient, accessToken string, dumper Dumper) *AsanaService {
	return &AsanaService{
		client:      client,
		accessToken: accessToken,
		dataDumper:  dumper,
	}
}

func (a AsanaService) GetUsers(ctx context.Context, request clients.GetUsersRequest) (models.AsanaGetUsersResponse, error) {
	request.Token = a.accessToken
	response, err := a.client.GetUsers(ctx, request)
	if err != nil {
		return models.AsanaGetUsersResponse{}, err
	}

	//a.dataDumper.DumpList(ctx, response.Data)

	return response, nil
}

func (a AsanaService) GetProjects(ctx context.Context, request clients.GetProjectsRequest) (models.AsanaGetProjectsResponse, error) {
	request.Token = a.accessToken
	response, err := a.client.GetProjects(ctx, request)
	if err != nil {
		return models.AsanaGetProjectsResponse{}, err
	}

	//a.dataDumper.DumpList(ctx, response.Data)

	return response, nil
}
