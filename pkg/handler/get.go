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
	EndpointGetPolicy   = "/projects/{project}/policies/{policy}"
	EndpointGetPolicies = "/projects/{project}/policies"
	EndpointGetRule     = "/projects/{project}/policies/{policy}/rules/{priority}"
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

	resource, err := h.client.GetPolicy(h.ctx, projectID, policy)
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

	var policies []*compute.SecurityPolicy
	it := h.client.ListPolicies(h.ctx, projectID)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			if ok := h.HttpError(err, w, projectID, securityTypePolicy); !ok {
				policiesResponse(w, []*compute.SecurityPolicy{})
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

	resource, err := h.client.GetRule(h.ctx, &p, projectID, policy)
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
