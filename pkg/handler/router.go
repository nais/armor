package handler

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nais/armor/pkg/google"
	"github.com/sirupsen/logrus"
)

type Application struct {
	Router  *mux.Router
	Client  *google.Client
	Log     *logrus.Entry
	Context context.Context
}

func NewApp(ctx context.Context, client *google.Client, log *logrus.Entry) *Application {
	return &Application{
		Router:  mux.NewRouter().StrictSlash(true),
		Client:  client,
		Log:     log,
		Context: ctx,
	}
}

func (a *Application) SetupHttpRouter(h *Handler) *mux.Router {
	a.Router.HandleFunc(EndpointIsAlive, h.isAlive).Methods(http.MethodGet)
	a.Router.HandleFunc(EndpointIsReady, h.isReady).Methods(http.MethodGet)
	// Policy
	a.Router.HandleFunc(EndpointGetPolicy, h.GetPolicy).Methods(http.MethodGet)
	a.Router.HandleFunc(EndpointGetPolicies, h.GetPolicies).Methods(http.MethodGet)
	a.Router.HandleFunc(EndpointCreatePolicy, h.CreatePolicy).Methods(http.MethodPost)
	a.Router.HandleFunc(EndpointUpdatePolicy, h.UpdatePolicy).Methods(http.MethodPatch)
	a.Router.HandleFunc(EndpointDeletePolicy, h.DeletePolicy).Methods(http.MethodDelete)
	// Rule
	a.Router.HandleFunc(EndpointGetRule, h.GetRule).Methods(http.MethodGet)
	a.Router.HandleFunc(EndpointCreateRule, h.CreateRule).Methods(http.MethodPost)
	a.Router.HandleFunc(EndpointUpdateRule, h.UpdateRule).Methods(http.MethodPatch)
	a.Router.HandleFunc(EndpointDeleteRule, h.DeleteRule).Methods(http.MethodDelete)
	return a.Router
}
