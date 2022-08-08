package validation

import (
	"fmt"
	"google.golang.org/genproto/googleapis/cloud/compute/v1"
)

func Rule(rule *compute.SecurityPolicyRule) (bool, error) {
	if rule.Action == nil {
		return false, fmt.Errorf("action is required")
	}

	if rule.Priority == nil {
		return false, fmt.Errorf("priority is required")
	}

	if rule.Preview == nil {
		return false, fmt.Errorf("preview is required")
	}

	if rule.Match == nil {
		return false, fmt.Errorf("match is required")
	}

	if rule.Action != nil && *rule.Action == "rate_based_ban" || *rule.Action == "throttle" {
		if rule.RateLimitOptions == nil {
			return false, fmt.Errorf("rate limit options is required when rate_based_ban or throttle is used")
		}
	}

	if rule.Action != nil && *rule.Action == "redirect" {
		if rule.RedirectOptions == nil {
			return false, fmt.Errorf("redirect options is required when redirect is used")
		}
	}

	if rule.Match.VersionedExpr != nil {
		if rule.Match.Config == nil {
			return false, fmt.Errorf("match config is required when match versioned expr is used")
		}
	}

	return true, nil
}
