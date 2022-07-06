package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	"github.com/nais/armor/pkg/model"
	"github.com/sirupsen/logrus"
	"google.golang.org/genproto/googleapis/cloud/compute/v1"
)

const (
	EndpointCreatePolicy = "/{project}/policy"
	EndpointCreateRule   = "/{project}/policy/{policy}/rule"
)

func (h *Handler) CreatePolicy(w http.ResponseWriter, r *http.Request) {
	h.log.WithFields(logrus.Fields{
		"method": "CreatePolicy",
	})

	projectID := mux.Vars(r)["project"]
	if ok, value := parse(projectID); !ok {
		http.Error(w, fmt.Sprintf("unkown parameter: %s", value), http.StatusBadRequest)
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil || len(reqBody) == 0 {
		http.Error(w, "error body", http.StatusBadRequest)
		return
	}

	request := model.ArmorRequestPolicy{}
	err = json.Unmarshal(reqBody, &request)
	if err != nil {
		h.log.Errorf("parse policy %v", err)
		http.Error(w, fmt.Sprintf("parse policy for project %s", projectID), http.StatusBadRequest)
		return
	}

	resource, err := h.createPolicy(&request, projectID)
	if err != nil {
		if ok := h.HttpError(err, w, projectID, securityTypePolicy); !ok {
			policyResponse(w, &compute.SecurityPolicy{})
			return
		}
		h.log.Errorf("error creating policy %v", err)
		http.Error(w, fmt.Sprintf("create policy %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	policiesResponse(w, []*compute.SecurityPolicy{resource})
	return
}

func (h *Handler) createPolicy(request *model.ArmorRequestPolicy, projectID string) (*compute.SecurityPolicy, error) {
	parsedPolicy, err := request.ParsePolicy()
	if err != nil {
		return nil, err
	}

	if parsedPolicy.Name == nil {
		return nil, fmt.Errorf("policy name is required")
	}

	if parsedPolicy.Rules != nil {
		if len(parsedPolicy.Rules) < 1 {
			securityPolicyRule := defaultRule(request.DefaultRuleAction)
			parsedPolicy.Rules = append(parsedPolicy.Rules, securityPolicyRule)
		}
	}

	if ok, err := h.client.CreatePolicy(h.ctx, parsedPolicy, projectID); !ok {
		if err != nil {
			return nil, err
		}

		h.log.Info("inserted policy ", parsedPolicy.Name)
	}
	return parsedPolicy, nil
}

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

func (h *Handler) CreateRule(w http.ResponseWriter, r *http.Request) {
	h.log.WithFields(logrus.Fields{
		"method": "CreateRule",
	})

	projectID := mux.Vars(r)["project"]
	policy := mux.Vars(r)["policy"]

	if ok, value := parse(projectID, policy); !ok {
		http.Error(w, fmt.Sprintf("unkown parameter: %s", value), http.StatusBadRequest)
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil || len(reqBody) == 0 {
		http.Error(w, "error body", http.StatusBadRequest)
		return
	}

	request := model.ArmorRequestRule{}
	err = json.Unmarshal(reqBody, &request)
	if err != nil {
		h.log.Errorf("parse rule %v", err)
		http.Error(w, fmt.Sprintf("parse request body for project %s: policy %s", projectID, policy), http.StatusBadRequest)
		return
	}

	resource, err := request.ParseRule()

	if ok, err := validRule(resource); !ok {
		h.log.Errorf("error validation of rule %v", err)
		http.Error(w, fmt.Sprintf("validation of rule: %v", err), http.StatusBadRequest)
		return
	}

	if err != nil {
		h.log.Errorf("parse rule %v", err)
		http.Error(w, fmt.Sprintf("parse rule for project %s: policy %s", projectID, policy), http.StatusInternalServerError)
		return
	}

	if ok, err := h.client.AddRule(h.ctx, resource, projectID, policy); !ok {
		if err != nil {
			if ok := h.HttpError(err, w, projectID, securityTypeRule); !ok {
				policyResponse(w, &compute.SecurityPolicy{})
				return
			}
			h.log.Errorf("error adding rule %v", err)
			http.Error(w, fmt.Sprintf("adding rule %s", projectID), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	ruleResponse(w, resource)
	return
}

func validRule(rule *compute.SecurityPolicyRule) (bool, error) {
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
