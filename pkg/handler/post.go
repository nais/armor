package handler

import (
	"encoding/json"
	"fmt"
	"github.com/nais/armor/pkg/validation"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nais/armor/pkg/model"
	"github.com/sirupsen/logrus"
	"google.golang.org/genproto/googleapis/cloud/compute/v1"
)

const (
	EndpointCreatePolicy     = "/projects/{project}/policies"
	EndpointCreateRule       = "/projects/{project}/policies/{policy}/rules"
	EndpointSetPolicyBackend = "/projects/{project}/policies/{policy}/backendServices/{backend}"
)

func (h *Handler) CreatePolicy(w http.ResponseWriter, r *http.Request) {
	h.log.WithFields(logrus.Fields{
		"method": "CreatePolicy",
	})

	projectID := mux.Vars(r)["project"]
	if ok, value := parse(projectID); !ok {
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
		h.log.Errorf("parse policy %v", err)
		HttpError(w, fmt.Sprintf("parse policy for project %s", projectID), http.StatusBadRequest)
		return
	}

	resource, err := h.createPolicy(&request, projectID)
	if err != nil {
		h.log.Errorf("error creating policy %v", err)
		h.HttpError(err, w, projectID, securityTypePolicy)
		return
	}

	h.log.Debug("inserted policy ", resource.Name)
	w.WriteHeader(http.StatusCreated)
	return
}

func (h *Handler) CreateRule(w http.ResponseWriter, r *http.Request) {
	h.log.WithFields(logrus.Fields{
		"method": "CreateRule",
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

	request := model.ArmorRequestRule{}
	err = json.Unmarshal(reqBody, &request)
	if err != nil {
		h.log.Errorf("parse rule %v", err)
		HttpError(w, fmt.Sprintf("parse request body for project %s: policy %s", projectID, policy), http.StatusBadRequest)
		return
	}

	resource, err := request.ParseRule()
	if ok, err := validation.Rule(resource); !ok {
		h.log.Errorf("error validation of rule %v", err)
		HttpError(w, fmt.Sprintf("validation of rule: %v", err), http.StatusBadRequest)
		return
	}

	if h.cfg.IsProtectedRule(strconv.Itoa(int(*resource.Priority))) {
		HttpError(w, fmt.Sprintf("forbidden to create protected priority %d", *resource.Priority), http.StatusBadRequest)
		return
	}

	if err != nil {
		h.log.Errorf("parse rule %v", err)
		HttpError(w, fmt.Sprintf("parse rule for project %s: policy %s", projectID, policy), http.StatusInternalServerError)
		return
	}

	if ok, err := h.securityClient.AddRule(h.ctx, resource, projectID, policy); !ok {
		if err != nil {
			h.log.Errorf("error adding rule %v", err)
			h.HttpError(err, w, projectID, securityTypeRule)
			return
		}
	}

	h.log.Debug("inserted rule ", resource.Priority)
	w.WriteHeader(http.StatusCreated)
	return
}

func (h *Handler) SetPolicyBackend(w http.ResponseWriter, r *http.Request) {
	h.log.WithFields(logrus.Fields{
		"method": "SetPolicyBackend",
	})

	projectID := mux.Vars(r)["project"]
	policy := mux.Vars(r)["policy"]
	backend := mux.Vars(r)["backend"]

	if ok, value := parse(projectID, policy, backend); !ok {
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

	if ok, err := h.serviceClient.SetSecurityPolicy(h.ctx, projectID, resource.SelfLink, backend); !ok {
		if err != nil {
			h.log.Errorf("error setting policy backend %v", err)
			h.HttpError(err, w, projectID, securityTypeRule)
			response(w, interface{}(&compute.SecurityPolicy{}))
			return
		}
	}

	h.log.Debug("policy backend is set to ", backend)
	w.WriteHeader(http.StatusCreated)
}
