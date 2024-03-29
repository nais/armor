package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/nais/armor/config"
	"github.com/nais/armor/pkg/google"
	"github.com/nais/armor/pkg/model"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/googleapi"
	"google.golang.org/genproto/googleapis/cloud/compute/v1"
	"net/http"
)

type Handler struct {
	ctx            context.Context
	log            *logrus.Entry
	cfg            *config.Config
	securityClient *google.SecurityClient
	serviceClient  *google.ServiceClient
}

const (
	securityTypePolicy = "policy"
	securityTypeRule   = "rule"
)

func NewHandler(ctx context.Context, cfg *config.Config, securityClient *google.SecurityClient, serviceClient *google.ServiceClient, log *logrus.Entry) *Handler {
	return &Handler{
		log:            log.WithField("subsystem", "handler"),
		securityClient: securityClient,
		serviceClient:  serviceClient,
		ctx:            ctx,
		cfg:            cfg,
	}
}

func (h *Handler) createPolicy(request *model.ArmorRequestPolicy, projectID string) (*compute.SecurityPolicy, error) {
	parsedPolicy, err := request.ParsePolicy()
	if err != nil {
		return nil, err
	}

	if parsedPolicy.Name == nil {
		return nil, fmt.Errorf("policy name is required")
	}

	if ok, err := h.securityClient.CreatePolicy(h.ctx, parsedPolicy, projectID); !ok {
		if err != nil {
			return nil, err
		}
	}
	return parsedPolicy, nil
}

func (h *Handler) HttpError(err error, w http.ResponseWriter, projectID, resource string) {
	if ErrorType(err, http.StatusNotFound) {
		w.WriteHeader(http.StatusNotFound)
	}
	if ErrorType(err, http.StatusBadRequest) {
		h.log.Warnf("%s: %v", resource, err)
		HttpError(w, fmt.Sprintf("%s resource %s: %s", resource, projectID, err.Error()), http.StatusBadRequest)
	}
	if ErrorType(err, http.StatusConflict) {
		h.log.Warnf("failed %s: %v", resource, err)
		HttpError(w, fmt.Sprintf("%s exists in %s", resource, projectID), http.StatusConflict)
	}
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
