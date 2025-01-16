package models

import (
	"time"
)

type ResourceList []TypedResource

type AsanaGetUsersResponse struct {
	Data     []AsanaUser   `json:"data"`
	NextPage AsanaNextPage `json:"next_page,omitempty"`
}

type TypedResource interface {
	GetGid() string
	GetResourceType() string
}

type BaseResource struct {
	Gid          string `json:"gid"`
	ResourceType string `json:"resource_type"`
}

type AsanaUser struct {
	BaseResource
	Email      string           `json:"email"`
	Name       string           `json:"name"`
	Photo      AsanaPhoto       `json:"photo"`
	Workspaces []AsanaWorkspace `json:"workspaces"`
}

func (t AsanaUser) GetGid() string {
	return t.BaseResource.Gid
}

func (t AsanaUser) GetResourceType() string {
	return t.BaseResource.ResourceType
}

type AsanaPhoto struct {
	Image1024x1024 string `json:"image_1024x1024"`
	Image128x128   string `json:"image_128x128"`
	Image21x21     string `json:"image_21x21"`
	Image27x27     string `json:"image_27x27"`
	Image36x36     string `json:"image_36x36"`
	Image60x60     string `json:"image_60x60"`
}

type AsanaWorkspace struct {
	BaseResource
	Name string `json:"name"`
}

type AsanaNextPage struct {
	Offset string `json:"offset"`
	Path   string `json:"path"`
	Uri    string `json:"uri"`
}

type AsanaGetProjectsResponse struct {
	Data     []AsanaProjectResource `json:"data"`
	NextPage AsanaNextPage          `json:"next_page,omitempty"`
}

type AsanaProjectResource struct {
	BaseResource
	Archived      bool               `json:"archived"`
	Color         string             `json:"color"`
	CreatedAt     *time.Time         `json:"created_at"`
	CurrentStatus AsanaProjectStatus `json:"current_status"`
}

func (t AsanaProjectResource) GetGid() string {
	return t.BaseResource.Gid
}

func (t AsanaProjectResource) GetResourceType() string {
	return t.BaseResource.ResourceType
}

type AsanaAuthor struct {
	BaseResource
	Name string `json:"name"`
}

type AsanaProjectStatus struct {
	ResourceType        string `json:"resource_type"`
	*AsanaAuthor        `json:"author,omitempty"`
	Color               string                   `json:"color"`
	CreatedAt           *time.Time               `json:"created_at"`
	CreatedBy           *AsanaAuthor             `json:"created_by,omitempty"`
	HtmlText            string                   `json:"html_text"`
	ModifiedAt          *time.Time               `json:"modified_at"`
	Text                string                   `json:"text"`
	Title               string                   `json:"title"`
	CurrentStatusUpdate AsanaProjectStatusUpdate `json:"current_status_update"`
}

type AsanaProjectStatusUpdate struct {
	BaseResource
	ResourceSubtype string `json:"resource_subtype"`
	Title           string `json:"title"`
}
