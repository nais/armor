package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
	"google.golang.org/genproto/googleapis/cloud/compute/v1"
	"net/http"
)

const (
	EndpointGetPolicy             = "/projects/{project}/policies/{policy}"
	EndpointGetPolicies           = "/projects/{project}/policies"
	EndpointGetRule               = "/projects/{project}/policies/{policy}/rules/{priority}"
	EndpointGetPreConfiguredRules = "/projects/{project}/preConfiguredRules"
	EndpointGetBackendServices    = "/projects/{project}/backendServices"
)

func (h *Handler) GetPolicy(w http.ResponseWriter, r *http.Request) {
	h.log.WithFields(logrus.Fields{
		"method": "GetPolicy",
	})

	projectID := mux.Vars(r)["project"]
	policy := mux.Vars(r)["policy"]

	if ok, value := parse(projectID, policy); !ok {
		HttpError(w, fmt.Sprintf("unkown parameter: %s", value), http.StatusBadRequest)
		return
	}

	resource, err := h.securityClient.GetPolicy(h.ctx, projectID, policy)
	if err != nil {
		h.log.Errorf("failed to get policy %s: %v", policy, err)
		h.HttpError(err, w, projectID, securityTypePolicy)
		response(w, interface{}(&compute.SecurityPolicy{}))
		return
	}

	h.log.Debug("got policy: ", resource)
	response(w, interface{}(resource))
	return
}

func (h *Handler) GetPolicies(w http.ResponseWriter, r *http.Request) {
	h.log.WithFields(logrus.Fields{
		"method": "GetPolicies",
	})

	projectID := mux.Vars(r)["project"]

	if ok, value := parse(projectID); !ok {
		HttpError(w, fmt.Sprintf("unkown parameter: %s", value), http.StatusBadRequest)
		return
	}

	policies := []*compute.SecurityPolicy{}
	it := h.securityClient.ListPolicies(h.ctx, projectID)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err == nil {
			policies = append(policies, resp)
			continue
		}

		h.log.Errorf("failed to list policies %s: %v", projectID, err)
		h.HttpError(err, w, projectID, securityTypePolicy)
		response(w, interface{}(policies))
		return
	}

	h.log.Debug("got policies: ", policies)
	response(w, interface{}(policies))
	return
}

func (h *Handler) GetRule(w http.ResponseWriter, r *http.Request) {
	h.log.WithFields(logrus.Fields{
		"method": "GetRule",
	})

	projectID := mux.Vars(r)["project"]
	policy := mux.Vars(r)["policy"]
	priority := mux.Vars(r)["priority"]

	if ok, value := parse(projectID, policy, priority); !ok {
		HttpError(w, fmt.Sprintf("unkown parameter: %s", value), http.StatusBadRequest)
		return
	}

	p, err := parseInt(priority)
	if err != nil {
		h.log.Errorf("failed to parse priority %s: %v", priority, err)
		HttpError(w, fmt.Sprintf("parse priority: %s", priority), http.StatusInternalServerError)
		return
	}

	resource, err := h.securityClient.GetRule(h.ctx, &p, projectID, policy)
	if err != nil {
		h.log.Errorf("failed to get rule %s: %v", policy, err)
		h.HttpError(err, w, projectID, securityTypeRule)
		response(w, interface{}(&compute.SecurityPolicyRule{}))
		return
	}

	h.log.Debug("got rule: ", resource)
	response(w, interface{}(resource))
	return
}

func (h *Handler) GetPreConfiguredRules(w http.ResponseWriter, r *http.Request) {
	h.log.WithFields(logrus.Fields{
		"method": "GetPreConfiguredRules",
	})

	projectID := mux.Vars(r)["project"]
	ruleType := r.URL.Query().Get("rule-type")
	version := r.URL.Query().Get("version")

	if ok, value := parse(projectID, ruleType, version); !ok {
		HttpError(w, fmt.Sprintf("unkown parameter: %s", value), http.StatusBadRequest)
		return
	}

	resource, err := h.securityClient.ListPreConfiguredRules(h.ctx, projectID)
	var filteredResponse []*compute.WafExpressionSet
	if err != nil {
		h.log.Errorf("failed to pre configured rules for %s: %v", projectID, err)
		h.HttpError(err, w, projectID, securityTypeRule)
		response(w, interface{}(filteredResponse))
		return
	}

	h.log.Debug("got pre configured rules: ", resource)

	filteredResponse = filterResult(ruleType, version, resource.GetPreconfiguredExpressionSets().GetWafRules().GetExpressionSets())
	response(w, interface{}(filteredResponse))
	return
}

func (h *Handler) GetBackendServices(w http.ResponseWriter, r *http.Request) {
	h.log.WithFields(logrus.Fields{
		"method": "GetBackendServices",
	})

	projectID := mux.Vars(r)["project"]

	if ok, value := parse(projectID); !ok {
		HttpError(w, fmt.Sprintf("unkown parameter: %s", value), http.StatusBadRequest)
		return
	}

	backends := []*compute.BackendService{}
	it := h.serviceClient.ListBackendServices(h.ctx, projectID)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err == nil {
			backends = append(backends, resp)
			continue
		}

		h.log.Errorf("failed to list backend services %s: %v", projectID, err)
		h.HttpError(err, w, projectID, securityTypePolicy)
		response(w, interface{}(backends))
		return
	}

	h.log.Debug("got backend services: ", backends)
	response(w, interface{}(backends))
}
