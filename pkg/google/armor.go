package google

import (
	"context"
	"fmt"

	compute "cloud.google.com/go/compute/apiv1"
	"github.com/nais/armor/config"
	"github.com/sirupsen/logrus"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
)

type Client struct {
	log    *logrus.Entry
	Client *compute.SecurityPoliciesClient
	Config *config.Config
}

func NewClient(cfg *config.Config, ctx context.Context, log *logrus.Entry) *Client {
	c, err := compute.NewSecurityPoliciesRESTClient(ctx)
	if err != nil {
		fmt.Printf("NewInstancesRESTClient: %v", err)
	}
	log.Info("created armored client")

	return &Client{
		log:    log,
		Client: c,
		Config: cfg,
	}
}

func (in *Client) ListPolicies(ctx context.Context, projectID string) *compute.SecurityPolicyIterator {
	req := &computepb.ListSecurityPoliciesRequest{
		Project: projectID,
	}

	return in.Client.List(ctx, req)
}

func (in *Client) GetPolicy(ctx context.Context, projectID, policyName string) (*computepb.SecurityPolicy, error) {
	req := &computepb.GetSecurityPolicyRequest{
		Project:        projectID,
		SecurityPolicy: policyName,
	}

	result, err := in.Client.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("get policy: %w", err)
	}

	return result, nil
}

func (in *Client) CreatePolicy(ctx context.Context, policy *computepb.SecurityPolicy, projectID string) (bool, error) {
	req := &computepb.InsertSecurityPolicyRequest{
		Project:                projectID,
		SecurityPolicyResource: policy,
	}

	op, err := in.Client.Insert(ctx, req)
	if err != nil {
		return false, fmt.Errorf("insert policy: %w", err)
	}

	err = op.Wait(ctx)
	if err != nil {
		return false, fmt.Errorf("wait policy: %w", err)
	}

	return op.Done(), nil
}

func (in *Client) DeletePolicy(ctx context.Context, projectID, policyName string) (bool, error) {
	req := &computepb.DeleteSecurityPolicyRequest{
		SecurityPolicy: policyName,
		Project:        projectID,
	}

	op, err := in.Client.Delete(ctx, req)
	if err != nil {
		return false, fmt.Errorf("delete policy: %w", err)
	}

	err = op.Wait(ctx)
	if err != nil {
		return false, fmt.Errorf("wait policy: %w", err)
	}

	return op.Done(), nil
}

func (in *Client) UpdatePolicy(ctx context.Context, policy *computepb.SecurityPolicy, projectID, policyName string) (bool, error) {
	req := &computepb.PatchSecurityPolicyRequest{
		SecurityPolicy:         policyName,
		Project:                projectID,
		SecurityPolicyResource: policy,
	}

	op, err := in.Client.Patch(ctx, req)
	if err != nil {
		return false, fmt.Errorf("update policy: %w", err)
	}

	err = op.Wait(ctx)
	if err != nil {
		return false, fmt.Errorf("wait policy: %w", err)
	}

	return op.Done(), nil
}

func (in *Client) GetRule(ctx context.Context, priority *int32, projectID, policyName string) (*computepb.SecurityPolicyRule, error) {
	req := &computepb.GetRuleSecurityPolicyRequest{
		SecurityPolicy: policyName,
		Project:        projectID,
		Priority:       priority,
	}

	rule, err := in.Client.GetRule(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("add rule: %w", err)
	}

	return rule, nil
}

func (in *Client) AddRule(ctx context.Context, resource *computepb.SecurityPolicyRule, projectID, policyName string) (bool, error) {
	req := &computepb.AddRuleSecurityPolicyRequest{
		SecurityPolicy:             policyName,
		Project:                    projectID,
		SecurityPolicyRuleResource: resource,
	}

	op, err := in.Client.AddRule(ctx, req)
	if err != nil {
		return false, fmt.Errorf("add rule: %w", err)
	}

	err = op.Wait(ctx)
	if err != nil {
		return false, fmt.Errorf("wait rule: %w", err)
	}

	return op.Done(), nil
}

func (in *Client) UpdateRule(ctx context.Context, resource *computepb.SecurityPolicyRule, projectID, policyName string) (bool, error) {
	req := &computepb.PatchRuleSecurityPolicyRequest{
		SecurityPolicy:             policyName,
		Project:                    projectID,
		SecurityPolicyRuleResource: resource,
	}

	op, err := in.Client.PatchRule(ctx, req)
	if err != nil {
		return false, fmt.Errorf("patch rule: %w", err)
	}

	err = op.Wait(ctx)
	if err != nil {
		return false, fmt.Errorf("wait rule: %w", err)
	}

	return op.Done(), nil
}

func (in *Client) RemoveRule(ctx context.Context, priority *int32, projectID, policyName string) (bool, error) {
	req := &computepb.RemoveRuleSecurityPolicyRequest{
		SecurityPolicy: policyName,
		Project:        projectID,
		Priority:       priority,
	}

	op, err := in.Client.RemoveRule(ctx, req)
	if err != nil {
		return false, fmt.Errorf("remove rule: %w", err)
	}

	err = op.Wait(ctx)
	if err != nil {
		return false, fmt.Errorf("wait rule: %w", err)
	}

	return op.Done(), nil
}
