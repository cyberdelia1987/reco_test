package app

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"

	"github.com/cyber/test-project/controllers"
	"github.com/cyber/test-project/services"
)

type AsanaService interface {
	services.AsanaUsersGetter
	services.AsanaProjectsGetter
}

type RouterConfig struct {
	AsanaService AsanaService
}

const pathPrefix = "/api/"

func NewRouter(cfg RouterConfig) (*mux.Router, error) {
	router := mux.NewRouter()

	chain := alice.New()

	baseRouter := router.PathPrefix(pathPrefix).Subrouter()

	baseRouter.
		Path("/users/get").
		Methods(http.MethodGet).
		Handler(chain.ThenFunc(controllers.AsanaGetUsers(cfg.AsanaService)))

	baseRouter.
		Path("/projects/get").
		Methods(http.MethodGet).
		Handler(chain.ThenFunc(controllers.AsanaGetProjects(cfg.AsanaService)))

	return router, nil
}
