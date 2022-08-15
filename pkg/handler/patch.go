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
		HttpError(w, fmt.Sprintf("unkown parameter: %s", value), http.StatusBadRequest)
		return
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil || len(reqBody) == 0 {
		HttpError(w, "error body", http.StatusBadRequest)
		return
	}

	request := model.ArmorRequestPolicy{}
	err = json.Unmarshal(reqBody, &request)
	if err != nil {
		h.log.Errorf("parse rule %v", err)
		HttpError(w, fmt.Sprintf("parse request body for project %s: policy %s", projectID, policy), http.StatusBadRequest)
		return
	}

	currentPolicy, err := h.securityClient.GetPolicy(h.ctx, projectID, policy)
	if err != nil {
		h.log.Errorf("failed to get policy %s: %v", policy, err)
		h.HttpError(err, w, projectID, securityTypePolicy)
		response(w, interface{}(&compute.SecurityPolicy{}))
		return
	}

	resource := compute.SecurityPolicy{}
	if err := request.MergePolicy(&resource, currentPolicy); err != nil {
		h.log.Warnf("failed to merge policy: %v", err)
		HttpError(w, fmt.Sprintf("merge policy %s for project %s", policy, projectID), http.StatusInternalServerError)
		response(w, interface{}(&compute.SecurityPolicy{}))
		return
	}

	if ok, err := h.securityClient.UpdatePolicy(h.ctx, &resource, projectID, policy); !ok {
		if err != nil {
			h.log.Errorf("failed to get policy %s: %v", policy, err)
			h.HttpError(err, w, projectID, securityTypePolicy)
			response(w, interface{}(&compute.SecurityPolicy{}))
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
		HttpError(w, fmt.Sprintf("unkown parameter: %s", value), http.StatusBadRequest)
		return
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil || len(reqBody) == 0 {
		HttpError(w, "error body", http.StatusBadRequest)
		return
	}

	request := model.ArmorRequestRule{}
	err = json.Unmarshal(reqBody, &request)
	if err != nil {
		h.log.Errorf("parse rule %v", err)
		HttpError(w, fmt.Sprintf("parse request body for project %s: policy %s", projectID, policy), http.StatusBadRequest)
		return
	}

	p, err := parseInt(priority)
	if err != nil {
		h.log.Errorf("failed to parse priority %s: %v", priority, err)
		HttpError(w, fmt.Sprintf("parse priority: %s", priority), http.StatusInternalServerError)
		return
	}

	currentRule, err := h.securityClient.GetRule(h.ctx, &p, projectID, policy)
	if err != nil {
		h.log.Errorf("failed to get rule %s: %v", priority, err)
		h.HttpError(err, w, projectID, securityTypeRule)
		response(w, interface{}(&compute.SecurityPolicyRule{}))
		return
	}

	if h.cfg.IsProtectedRule(priority) {
		HttpError(w, fmt.Sprintf("forbidden to update protected rule %s", priority), http.StatusBadRequest)
		return
	}

	resource := compute.SecurityPolicyRule{}
	if err := request.MergeRule(&resource, currentRule); err != nil {
		h.log.Warnf("failed to merge rule: %v", err)
		HttpError(w, fmt.Sprintf("merge rule %s for project %s", priority, projectID), http.StatusInternalServerError)
		response(w, interface{}(&compute.SecurityPolicyRule{}))
		return
	}

	if ok, err := h.securityClient.UpdateRule(h.ctx, &resource, projectID, policy); !ok {
		if err != nil {
			h.log.Errorf("failed to update rule %s: %v", priority, err)
			h.HttpError(err, w, projectID, securityTypeRule)
			response(w, interface{}(&compute.SecurityPolicyRule{}))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	return
}
