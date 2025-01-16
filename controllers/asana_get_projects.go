package controllers

import (
	"net/http"
	"strconv"

	"github.com/cyber/test-project/clients"
	"github.com/cyber/test-project/services"
	"github.com/cyber/test-project/transport"
)

func AsanaGetProjects(service services.AsanaProjectsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			limit = 50
		}

		var archived *bool
		if r.URL.Query().Has("archived") {
			archivedVal, err := strconv.ParseBool(r.URL.Query().Get("archived"))
			if err == nil {
				archived = &archivedVal
			}
		}

		req := clients.GetProjectsRequest{
			Workspace: r.URL.Query().Get("workspace"),
			Team:      r.URL.Query().Get("team"),
			Limit:     limit,
			Offset:    r.URL.Query().Get("offset"),
			Archived:  archived,
		}

		projects, err := service.GetProjects(ctx, req)
		if err != nil {
			transport.SendJson(ctx, w, http.StatusInternalServerError, err)
			return
		}

		transport.SendJson(ctx, w, http.StatusOK, projects)
	}
}
