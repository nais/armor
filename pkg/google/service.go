package google

import (
	compute "cloud.google.com/go/compute/apiv1"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
)

type ServiceClient struct {
	log    *logrus.Entry
	Client *compute.BackendServicesClient
}

func NewService(ctx context.Context, log *logrus.Entry) *ServiceClient {
	c, err := compute.NewBackendServicesRESTClient(ctx)
	if err != nil {
		fmt.Printf("NewBackendServicesRESTClient: %v", err)
	}

	return &ServiceClient{
		log:    log,
		Client: c,
	}
}

func (in *ServiceClient) SetSecurityPolicy(ctx context.Context, projectID, policy, backendService string) (bool, error) {
	req := &computepb.SetSecurityPolicyBackendServiceRequest{
		BackendService: backendService,
		Project:        projectID,
		SecurityPolicyReferenceResource: &computepb.SecurityPolicyReference{
			SecurityPolicy: &policy,
		},
	}
	op, err := in.Client.SetSecurityPolicy(ctx, req)
	if err != nil {
		if err != nil {
			return false, fmt.Errorf("insert policy to backend: %w", err)
		}
	}

	err = op.Wait(ctx)
	if err != nil {
		if err != nil {
			return false, fmt.Errorf("wait for backend: %w", err)
		}
	}

	return op.Done(), nil
}
