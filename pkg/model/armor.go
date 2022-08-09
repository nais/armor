package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/imdario/mergo"
	"google.golang.org/genproto/googleapis/cloud/compute/v1"
)

type ArmorRequestPolicy struct {
	DefaultRuleAction string                  `json:"default_rule_action,omitempty"`
	SecurityPolicy    *compute.SecurityPolicy `json:"policy"`
}

type ArmorRequestRule struct {
	SecurityPolicyRule *compute.SecurityPolicyRule `json:"rule"`
}

type ArmorPreConfiguredRulesRequest struct {
	Version string `json:"version"`
}

func (in *ArmorRequestPolicy) ParsePolicy() (*compute.SecurityPolicy, error) {
	var instance *compute.SecurityPolicy

	requestPolicy, err := ToBytes(in.SecurityPolicy)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(requestPolicy, &instance)
	if err != nil {
		return nil, fmt.Errorf("parse policy: %w", err)
	}

	return instance, nil
}

func (in *ArmorRequestRule) ParseRule() (*compute.SecurityPolicyRule, error) {
	var instance *compute.SecurityPolicyRule

	requestRule, err := ToBytes(in.SecurityPolicyRule)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(requestRule, &instance)
	if err != nil {
		return nil, fmt.Errorf("parse rule: %w", err)
	}
	return instance, nil
}

func (in *ArmorRequestPolicy) MergePolicy(dest, requestPolicy *compute.SecurityPolicy) error {
	policyUpdate, err := ToBytes(in.SecurityPolicy)
	if err != nil {
		return err
	}

	err = json.Unmarshal(policyUpdate, &dest)
	if err != nil {
		return fmt.Errorf("parse policy: %w", err)
	}

	// Rules cannot be updated with patch, please use addRule, removeRule, or patchRule instead
	requestPolicy.Rules = nil

	err = mergo.Merge(dest, requestPolicy)
	if err != nil {
		return fmt.Errorf("unable to merge default security policy values: %s", err)
	}
	return nil
}

func (in *ArmorRequestRule) MergeRule(dest, requestRule *compute.SecurityPolicyRule) error {
	ruleUpdate, err := ToBytes(in.SecurityPolicyRule)
	if err != nil {
		return err
	}

	err = json.Unmarshal(ruleUpdate, &dest)
	if err != nil {
		return fmt.Errorf("parse rule: %w", err)
	}

	err = mergo.Merge(dest, requestRule)
	if err != nil {
		return fmt.Errorf("unable to merge default security rule values: %s", err)
	}
	return nil
}

func ToBytes(resource interface{}) ([]byte, error) {
	reqBodyBytes := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBytes).Encode(resource)
	if err != nil {
		return nil, err
	}
	return reqBodyBytes.Bytes(), err
}
