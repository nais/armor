package handler

import (
	"github.com/golang/protobuf/proto"
	"google.golang.org/genproto/googleapis/cloud/compute/v1"
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
