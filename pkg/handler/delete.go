package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	EndpointDeletePolicy = "/{project}/policy/{policy}"
	EndpointDeleteRule   = "/{project}/policy/{policy}/rule/{priority}"
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

	if ok, err := h.client.DeletePolicy(h.ctx, projectID, policy); !ok {
		if err != nil {
			if ok := h.HttpError(err, w, projectID, securityTypePolicy); !ok {
				return
			}
			h.log.Errorf("failed to get policy %s: %v", policy, err)
			http.Error(w, fmt.Sprintf("get policy %s for project %s", policy, projectID), http.StatusInternalServerError)
			return
		}
	}
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

	p, err := parseInt(priority)
	if err != nil {
		h.log.Errorf("failed to parse priority %s: %v", priority, err)
		http.Error(w, fmt.Sprintf("parse priority: %s", priority), http.StatusInternalServerError)
		return
	}

	if ok, err := h.client.RemoveRule(h.ctx, &p, projectID, policy); !ok {
		if err != nil {
			if ok := h.HttpError(err, w, projectID, securityTypeRule); !ok {
				return
			}
			h.log.Errorf("failed to get rule %s: %v", priority, err)
			http.Error(w, fmt.Sprintf("get rule %s for project %s", priority, projectID), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	return
}
