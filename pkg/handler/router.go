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
	router := mux.NewRouter().StrictSlash(true)
	// Policy
	router.HandleFunc(EndpointGetPolicy, h.GetPolicy).Methods(http.MethodGet)
	router.HandleFunc(EndpointGetPolicies, h.GetPolicies).Methods(http.MethodGet)
	router.HandleFunc(EndpointCreatePolicy, h.CreatePolicy).Methods(http.MethodPost)
	router.HandleFunc(EndpointUpdatePolicy, h.UpdatePolicy).Methods(http.MethodPatch)
	router.HandleFunc(EndpointDeletePolicy, h.DeletePolicy).Methods(http.MethodDelete)
	// Rule
	router.HandleFunc(EndpointGetRule, h.GetRule).Methods(http.MethodGet)
	router.HandleFunc(EndpointCreateRule, h.CreateRule).Methods(http.MethodPost)
	router.HandleFunc(EndpointUpdateRule, h.UpdateRule).Methods(http.MethodPatch)
	router.HandleFunc(EndpointDeleteRule, h.DeleteRule).Methods(http.MethodDelete)
	return router
}
