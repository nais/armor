package handler

import (
	"encoding/json"
	"fmt"
	"google.golang.org/genproto/googleapis/cloud/compute/v1"
	"google.golang.org/protobuf/proto"
	"net/http"
	"regexp"
	"strconv"
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

func response(w http.ResponseWriter, response interface{}) {
	err := json.NewEncoder(w).Encode(response)
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
		if v == "" {
			continue
		}
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
