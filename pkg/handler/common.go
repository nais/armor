package handler

import (
	"fmt"
	"google.golang.org/genproto/googleapis/cloud/compute/v1"
	"google.golang.org/protobuf/proto"
	"strings"
)

func defaultRule(defaultRuleAction string) *compute.SecurityPolicyRule {
	if len(defaultRuleAction) == 0 {
		defaultRuleAction = "deny(403)"
	}

	matcher := &compute.SecurityPolicyRuleMatcher{
		Config: &compute.SecurityPolicyRuleMatcherConfig{
			SrcIpRanges: []string{
				"*",
			},
		},
		VersionedExpr: proto.String(compute.SecurityPolicyRuleMatcher_SRC_IPS_V1.String()),
	}

	action := defaultRuleAction
	return &compute.SecurityPolicyRule{
		Action: &action,
		// lowest priority
		Priority:    proto.Int32(2147483647),
		Description: proto.String("Default rule, higher priority overrides it"),
		Match:       matcher,
	}
}

func filterResult(filter, version string, resource *compute.SecurityPoliciesListPreconfiguredExpressionSetsResponse) (filteredResponse []*compute.WafExpressionSet) {
	if filter == "" {
		filteredResponse = resource.GetPreconfiguredExpressionSets().GetWafRules().GetExpressionSets()
	} else {
		for _, expression := range resource.GetPreconfiguredExpressionSets().GetWafRules().GetExpressionSets() {
			// v33 is the latest version of preconfigured rules
			if strings.Contains(expression.GetId(), fmt.Sprintf("%s-%s", filter, version)) {
				filteredResponse = append(filteredResponse, expression)
			}
		}
	}
	return filteredResponse
}
