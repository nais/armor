package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
	"google.golang.org/genproto/googleapis/cloud/compute/v1"
	"net/http"
	"strings"
)

const (
	EndpointGetPolicy             = "/projects/{project}/policies/{policy}"
	EndpointGetPolicies           = "/projects/{project}/policies"
	EndpointGetRule               = "/projects/{project}/policies/{policy}/rules/{priority}"
	EndpointGetPreConfiguredRules = "/projects/{project}/preConfiguredRules"
	EndpointGetBackendServices    = "/projects/{project}/backendservices"
)

func (h *Handler) GetPolicy(w http.ResponseWriter, r *http.Request) {
	h.log.WithFields(logrus.Fields{
		"method": "GetPolicy",
	})

	projectID := mux.Vars(r)["project"]
	policy := mux.Vars(r)["policy"]

	if ok, value := parse(projectID, policy); !ok {
		http.Error(w, fmt.Sprintf("unkown parameter: %s", value), http.StatusBadRequest)
		return
	}

	resource, err := h.securityClient.GetPolicy(h.ctx, projectID, policy)
	if err != nil {
		if ok := h.HttpError(err, w, projectID, securityTypePolicy); !ok {
			policyResponse(w, &compute.SecurityPolicy{})
			return
		}
		h.log.Errorf("failed to get policy %s: %v", policy, err)
		http.Error(w, fmt.Sprintf("get policy %s for project %s", policy, projectID), http.StatusInternalServerError)
		return
	}

	policyResponse(w, resource)
	return
}

func (h *Handler) GetPolicies(w http.ResponseWriter, r *http.Request) {
	h.log.WithFields(logrus.Fields{
		"method": "GetPolicies",
	})

	projectID := mux.Vars(r)["project"]

	if ok, value := parse(projectID); !ok {
		http.Error(w, fmt.Sprintf("unkown parameter: %s", value), http.StatusBadRequest)
		return
	}

	policies := []*compute.SecurityPolicy{}
	it := h.securityClient.ListPolicies(h.ctx, projectID)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			if ok := h.HttpError(err, w, projectID, securityTypePolicy); !ok {
				policiesResponse(w, policies)
				return
			}
			h.log.Errorf("failed to list policies %s: %v", projectID, err)
			http.Error(w, fmt.Sprintf("cant get polices for project: %s", projectID), http.StatusInternalServerError)
			return
		}
		policies = append(policies, resp)
	}

	policiesResponse(w, policies)
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
		http.Error(w, fmt.Sprintf("unkown parameter: %s", value), http.StatusBadRequest)
		return
	}

	p, err := parseInt(priority)
	if err != nil {
		h.log.Errorf("failed to parse priority %s: %v", priority, err)
		http.Error(w, fmt.Sprintf("parse priority: %s", priority), http.StatusInternalServerError)
		return
	}

	resource, err := h.securityClient.GetRule(h.ctx, &p, projectID, policy)
	if err != nil {
		if ok := h.HttpError(err, w, projectID, securityTypeRule); !ok {
			ruleResponse(w, &compute.SecurityPolicyRule{})
			return
		}
		h.log.Errorf("failed to get rule %s: %v", policy, err)
		http.Error(w, fmt.Sprintf("trying to get rule for policy %s", policy), http.StatusInternalServerError)
		return
	}

	ruleResponse(w, resource)
	return
}

func (h *Handler) GetPreConfiguredRules(w http.ResponseWriter, r *http.Request) {
	h.log.WithFields(logrus.Fields{
		"method": "GetPreConfigured",
	})

	projectID := mux.Vars(r)["project"]
	filter := r.URL.Query().Get("filter")

	if ok, value := parse(projectID); !ok {
		http.Error(w, fmt.Sprintf("unkown parameter: %s", value), http.StatusBadRequest)
		return
	}

	resource, err := h.securityClient.ListPreConfiguredRules(h.ctx, projectID)
	var filteredResponse []*compute.WafExpressionSet
	if err != nil {
		if ok := h.HttpError(err, w, projectID, securityTypeRule); !ok {
			preConfiguredResponse(w, filteredResponse)
			return
		}
		h.log.Errorf("failed to pre configured rules for %s: %v", projectID, err)
		http.Error(w, fmt.Sprintf("trying to get preconfigured rules for project %s", projectID), http.StatusInternalServerError)
		return
	}

	if filter == "" {
		filteredResponse = resource.GetPreconfiguredExpressionSets().WafRules.GetExpressionSets()
	} else {
		for _, expression := range resource.GetPreconfiguredExpressionSets().WafRules.GetExpressionSets() {
			if strings.Contains(expression.GetId(), filter) {
				filteredResponse = append(filteredResponse, expression)
			}
		}
	}

	preConfiguredResponse(w, filteredResponse)
	return
}

func (h *Handler) GetBackendServices(w http.ResponseWriter, r *http.Request) {
	h.log.WithFields(logrus.Fields{
		"method": "GetPolicies",
	})

	projectID := mux.Vars(r)["project"]

	if ok, value := parse(projectID); !ok {
		http.Error(w, fmt.Sprintf("unkown parameter: %s", value), http.StatusBadRequest)
		return
	}

	backends := []*compute.BackendService{}
	it := h.serviceClient.ListBackendServices(h.ctx, projectID)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			if ok := h.HttpError(err, w, projectID, securityTypePolicy); !ok {
				backendResponse(w, backends)
				return
			}
			h.log.Errorf("failed to list backend services %s: %v", projectID, err)
			http.Error(w, fmt.Sprintf("cant get backend services for project: %s", projectID), http.StatusInternalServerError)
			return
		}
		backends = append(backends, resp)
	}

	backendResponse(w, backends)
	return
}
