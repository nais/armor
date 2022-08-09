package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nais/armor/pkg/model"
	"github.com/sirupsen/logrus"
	"google.golang.org/genproto/googleapis/cloud/compute/v1"
	"io"
	"net/http"
)

const (
	EndpointUpdatePolicy = "/projects/{project}/policies/{policy}"
	EndpointUpdateRule   = "/projects/{project}/policies/{policy}/rules/{priority}"
)

func (h *Handler) UpdatePolicy(w http.ResponseWriter, r *http.Request) {
	h.log.WithFields(logrus.Fields{
		"method": "UpdatePolicy",
	})

	projectID := mux.Vars(r)["project"]
	policy := mux.Vars(r)["policy"]

	if ok, value := parse(projectID, policy); !ok {
		http.Error(w, fmt.Sprintf("unkown parameter: %s", value), http.StatusBadRequest)
		return
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil || len(reqBody) == 0 {
		http.Error(w, "error body", http.StatusBadRequest)
		return
	}

	request := model.ArmorRequestPolicy{}
	err = json.Unmarshal(reqBody, &request)
	if err != nil {
		h.log.Errorf("parse rule %v", err)
		http.Error(w, fmt.Sprintf("parse request body for project %s: policy %s", projectID, policy), http.StatusBadRequest)
		return
	}

	currentPolicy, err := h.getPolicy(projectID, policy)
	if err != nil {
		h.log.Errorf("failed to get policy %s: %v", policy, err)
		h.HttpError(err, w, projectID, securityTypePolicy)
		policyResponse(w, &compute.SecurityPolicy{})
		return
	}

	resource := compute.SecurityPolicy{}
	if err := request.MergePolicy(&resource, currentPolicy); err != nil {
		h.log.Warnf("failed to merge policy: %v", err)
		http.Error(w, fmt.Sprintf("merge policy %s for project %s", policy, projectID), http.StatusInternalServerError)
		policyResponse(w, &compute.SecurityPolicy{})
		return
	}

	if ok, err := h.securityClient.UpdatePolicy(h.ctx, &resource, projectID, policy); !ok {
		if err != nil {
			h.log.Errorf("failed to get policy %s: %v", policy, err)
			h.HttpError(err, w, projectID, securityTypePolicy)
			policyResponse(w, &compute.SecurityPolicy{})
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	return
}

func (h *Handler) UpdateRule(w http.ResponseWriter, r *http.Request) {
	h.log.WithFields(logrus.Fields{
		"method": "UpdateRule",
	})

	projectID := mux.Vars(r)["project"]
	policy := mux.Vars(r)["policy"]
	priority := mux.Vars(r)["priority"]

	if ok, value := parse(projectID, policy, priority); !ok {
		http.Error(w, fmt.Sprintf("unkown parameter: %s", value), http.StatusBadRequest)
		return
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil || len(reqBody) == 0 {
		http.Error(w, "error body", http.StatusBadRequest)
		return
	}

	request := model.ArmorRequestRule{}
	err = json.Unmarshal(reqBody, &request)
	if err != nil {
		h.log.Errorf("parse rule %v", err)
		http.Error(w, fmt.Sprintf("parse request body for project %s: policy %s", projectID, policy), http.StatusBadRequest)
		return
	}

	p, err := parseInt(priority)
	if err != nil {
		h.log.Errorf("failed to parse priority %s: %v", priority, err)
		http.Error(w, fmt.Sprintf("parse priority: %s", priority), http.StatusInternalServerError)
		return
	}

	currentRule, err := h.getRule(&p, projectID, policy)
	if err != nil {
		h.log.Errorf("failed to get rule %s: %v", priority, err)
		h.HttpError(err, w, projectID, securityTypeRule)
		ruleResponse(w, &compute.SecurityPolicyRule{})
		return
	}

	resource := compute.SecurityPolicyRule{}
	if err := request.MergeRule(&resource, currentRule); err != nil {
		h.log.Warnf("failed to merge rule: %v", err)
		http.Error(w, fmt.Sprintf("merge rule %s for project %s", priority, projectID), http.StatusInternalServerError)
		ruleResponse(w, &compute.SecurityPolicyRule{})
		return
	}

	if ok, err := h.securityClient.UpdateRule(h.ctx, &resource, projectID, policy); !ok {
		if err != nil {
			h.log.Errorf("failed to update rule %s: %v", priority, err)
			h.HttpError(err, w, projectID, securityTypeRule)
			ruleResponse(w, &compute.SecurityPolicyRule{})
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	return
}
