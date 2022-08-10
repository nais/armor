package handler

import (
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupHttpRouter(h *Handler) *mux.Router {
	log.WithField("method", "SetupHttpRouter").Debug("setting up http router")
	r := mux.NewRouter().StrictSlash(true)
	r.Use(commonMiddleware)

	r.HandleFunc(EndpointIsAlive, h.isAlive).Methods(http.MethodGet)
	r.HandleFunc(EndpointIsReady, h.isReady).Methods(http.MethodGet)
	// Policy
	r.HandleFunc(EndpointGetPolicy, h.GetPolicy).Methods(http.MethodGet)
	r.HandleFunc(EndpointGetPolicies, h.GetPolicies).Methods(http.MethodGet)
	r.HandleFunc(EndpointCreatePolicy, h.CreatePolicy).Methods(http.MethodPost)
	r.HandleFunc(EndpointUpdatePolicy, h.UpdatePolicy).Methods(http.MethodPatch)
	r.HandleFunc(EndpointDeletePolicy, h.DeletePolicy).Methods(http.MethodDelete)
	// Rule
	r.HandleFunc(EndpointGetRule, h.GetRule).Methods(http.MethodGet)
	r.HandleFunc(EndpointCreateRule, h.CreateRule).Methods(http.MethodPost)
	r.HandleFunc(EndpointUpdateRule, h.UpdateRule).Methods(http.MethodPatch)
	r.HandleFunc(EndpointDeleteRule, h.DeleteRule).Methods(http.MethodDelete)
	// Preconfigured rules
	r.HandleFunc(EndpointGetPreConfiguredRules, h.GetPreConfiguredRules).Methods(http.MethodGet)
	// Add policy to backend
	r.HandleFunc(EndpointSetPolicyBackend, h.SetPolicyBackend).Methods(http.MethodPost)
	r.HandleFunc(EndpointGetBackendServices, h.GetBackendServices).Methods(http.MethodGet)
	return r
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
