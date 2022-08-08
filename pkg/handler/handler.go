package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/nais/armor/pkg/google"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/googleapi"
	"google.golang.org/genproto/googleapis/cloud/compute/v1"
)

type Handler struct {
	log            *logrus.Entry
	securityClient *google.SecurityClient
	serviceClient  *google.ServiceClient
	ctx            context.Context
}

const (
	securityTypePolicy = "policy"
	securityTypeRule   = "rule"
)

func NewHandler(ctx context.Context, securityClient *google.SecurityClient, serviceClient *google.ServiceClient, log *logrus.Entry) *Handler {
	return &Handler{
		log:            log.WithField("subsystem", "handler"),
		securityClient: securityClient,
		serviceClient:  serviceClient,
		ctx:            ctx,
	}
}

func (h *Handler) getPolicy(projectID, policy string) (*compute.SecurityPolicy, error) {
	resource, err := h.securityClient.GetPolicy(h.ctx, projectID, policy)
	if err != nil {
		return nil, err
	}
	return resource, nil
}

func (h *Handler) getRule(priority *int32, projectID, policy string) (*compute.SecurityPolicyRule, error) {
	resource, err := h.securityClient.GetRule(h.ctx, priority, projectID, policy)
	if err != nil {
		return nil, err
	}
	return resource, nil
}

func (h *Handler) HttpError(err error, w http.ResponseWriter, projectID, resource string) bool {
	if ErrorType(err, http.StatusNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return false
	}
	if ErrorType(err, http.StatusBadRequest) {
		h.log.Warnf("%s: %v", resource, err)
		http.Error(w, fmt.Sprintf("%s resource %s: %s", resource, projectID, err.Error()), http.StatusBadRequest)
		return false
	}
	if ErrorType(err, http.StatusConflict) {
		h.log.Warnf("failed %s: %v", resource, err)
		http.Error(w, fmt.Sprintf("%s exists in %s", resource, projectID), http.StatusConflict)
		return false
	}
	return true
}

func ErrorType(err error, code int) bool {
	var e *googleapi.Error
	if ok := errors.As(err, &e); ok {
		if e.Code == code {
			return true
		}
	}
	return false
}

func policiesResponse(w http.ResponseWriter, policies []*compute.SecurityPolicy) {
	err := json.NewEncoder(w).Encode(policies)
	if err != nil {
		http.Error(w, fmt.Sprintf("encode %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
}

func policyResponse(w http.ResponseWriter, policy *compute.SecurityPolicy) {
	err := json.NewEncoder(w).Encode(policy)
	if err != nil {
		http.Error(w, fmt.Sprintf("encode %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
}

func ruleResponse(w http.ResponseWriter, rules *compute.SecurityPolicyRule) {
	err := json.NewEncoder(w).Encode(rules)
	if err != nil {
		http.Error(w, fmt.Sprintf("encode %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
}

func preConfiguredResponse(w http.ResponseWriter, rules []*compute.WafExpressionSet) {
	var err error
	if rules != nil {
		err = json.NewEncoder(w).Encode(rules)
	} else {
		err = json.NewEncoder(w).Encode([]*compute.WafExpressionSet{})
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("encode %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
}

func parse(input ...string) (bool, string) {
	// This will only match sequences of one or more sequences
	// of alphanumeric characters separated by a single -
	regex := "^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$"
	for _, v := range input {
		if !regexp.MustCompile(regex).MatchString(v) {
			return false, v
		}
	}
	return true, ""
}

func parseInt(i string) (int32, error) {
	p, err := strconv.ParseInt(i, 10, 32)
	if err != nil {
		return int32(0), err
	}
	return int32(p), nil
}
