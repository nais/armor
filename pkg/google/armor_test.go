package google

import (
	"context"
	"github.com/nais/armor/config"
	"github.com/nais/armor/pkg/fake"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/genproto/googleapis/cloud/compute/v1"
	"net/http/httptest"
	"testing"
)

func Test_ListPolicies(t *testing.T) {
	for _, test := range []struct {
		name       string
		project    string
		policyName string
		policies   int
		exists     bool
		runner     string
	}{
		{
			name:       "List all policies in a specific project",
			project:    "fake-project",
			policyName: "",
			policies:   3,
			exists:     true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, fakeClient, err := GoogleApiServer(test.policyName, test.exists)
			assert.NoError(t, err)

			it := fakeClient.ListPolicies(ctx, test.project)
			var policies []*compute.SecurityPolicy
			for {
				resp, err := it.Next()
				if err == iterator.Done {
					break
				}
				policies = append(policies, resp)
			}
			assert.Equal(t, test.policies, len(policies))
		})
	}
}

func Test_GetPolicy(t *testing.T) {
	existingPolicy := "test-2"
	ctx, fakeClient, err := GoogleApiServer(existingPolicy, true)
	assert.NoError(t, err)

	res, err := fakeClient.GetPolicy(ctx, "fake-project", "test-2")
	assert.NoError(t, err)

	assert.Equal(t, existingPolicy, *res.Name)
	assert.Equal(t, "test policy YOLO", *res.Description)
	assert.Equal(t, "compute#securityPolicy", *res.Kind)
	assert.Equal(t, "CLOUD_ARMOR", *res.Type)
	assert.Equal(t, "https://www.googleapis.com/compute/v1/projects/fake-project/global/securityPolicies/test-2", *res.SelfLink)
	assert.Equal(t, 3, len(res.Rules))
	assert.Equal(t, uint64(5663025914644165958), *res.Id)
	assert.Equal(t, false, *res.AdaptiveProtectionConfig.Layer7DdosDefenseConfig.Enable)
}

func Test_GetRule(t *testing.T) {
	existingRulePriority := int32(0)
	ctx, fakeClient, err := GoogleApiServer("test-2", true)
	assert.NoError(t, err)

	res, err := fakeClient.GetRule(ctx, &existingRulePriority, "fake-project", "test-2")
	assert.NoError(t, err)
	assert.Equal(t, existingRulePriority, *res.Priority)
}

func Test_UpdatePolicy(t *testing.T) {
	description := "test policy YOLO"
	securityPolicy := &compute.SecurityPolicy{
		Description: &description,
	}

	ctx, fakeClient, err := GoogleApiServer("test-2", true)
	assert.NoError(t, err)

	res, err := fakeClient.UpdatePolicy(ctx, securityPolicy, "fake-project", "test-2")
	assert.NoError(t, err)
	assert.Equal(t, true, res)
}

func Test_UpdateRule(t *testing.T) {
	description := "test policy YOLO"
	securityPolicyRule := &compute.SecurityPolicyRule{
		Description: &description,
	}

	ctx, fakeClient, err := GoogleApiServer("test-2", true)
	assert.NoError(t, err)

	res, err := fakeClient.UpdateRule(ctx, securityPolicyRule, "fake-project", "test-2")
	assert.NoError(t, err)
	assert.Equal(t, true, res)
}

func GoogleApiServer(existingPolicy string, exists bool) (context.Context, *SecurityClient, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, nil, err
	}

	ctx := context.Background()
	testServer := httptest.NewServer(fake.NewSecurityPoliciesRESTClient(existingPolicy, exists))
	opts := []option.ClientOption{
		option.WithEndpoint(testServer.URL),
		option.WithoutAuthentication(),
	}
	fakeClient := NewSecurityClient(cfg, ctx, log.WithField("component", "fake-client"), opts...)
	return ctx, fakeClient, nil
}
