package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	EndpointDeletePolicy = "/projects/{project}/policies/{policy}"
	EndpointDeleteRule   = "/projects/{project}/policies/{policy}/rules/{priority}"
)

func (h *Handler) DeletePolicy(w http.ResponseWriter, r *http.Request) {
	h.log.WithFields(logrus.Fields{
		"method": "DeletePolicy",
	})

	projectID := mux.Vars(r)["project"]
	policy := mux.Vars(r)["policy"]

	if ok, value := parse(projectID, policy); !ok {
		http.Error(w, fmt.Sprintf("unkown parameter: %s", value), http.StatusBadRequest)
		return
	}

	if ok, err := h.securityClient.DeletePolicy(h.ctx, projectID, policy); !ok {
		if err != nil {
			h.log.Errorf("failed to delete policy %s: %v", policy, err)
			h.HttpError(err, w, projectID, securityTypePolicy)
			return
		}
	}

	h.log.Debug("deleted policy: ", policy)
	w.WriteHeader(http.StatusOK)
	return
}

func (h *Handler) DeleteRule(w http.ResponseWriter, r *http.Request) {
	h.log.WithFields(logrus.Fields{
		"method": "DeleteRule",
	})

	projectID := mux.Vars(r)["project"]
	policy := mux.Vars(r)["policy"]
	priority := mux.Vars(r)["priority"]

	if ok, value := parse(projectID, policy, priority); !ok {
		http.Error(w, fmt.Sprintf("unkown parameter: %s", value), http.StatusBadRequest)
		return
	}

	var err error
	p, err := parseInt(priority)
	if err != nil {
		h.log.Errorf("failed to parse priority %s: %v", priority, err)
		http.Error(w, fmt.Sprintf("parse priority: %s", priority), http.StatusInternalServerError)
		return
	}

	if ok, err := h.securityClient.RemoveRule(h.ctx, &p, projectID, policy); !ok {
		if err != nil {
			h.log.Errorf("failed to get rule %s: %v", priority, err)
			h.HttpError(err, w, projectID, securityTypeRule)
			return
		}
	}

	h.log.Debug("deleted rule: ", policy)
	w.WriteHeader(http.StatusOK)
	return
}
