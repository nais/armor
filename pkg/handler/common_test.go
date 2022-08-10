package handler

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/googleapis/cloud/compute/v1"
	"os"
	"testing"
)

func Test_filterResult(t *testing.T) {
	for _, test := range []struct {
		name     string
		ruleType string
		version  string
		result   int
	}{
		{
			name:     "No rule type, just raw list of preconfigured rules",
			ruleType: "",
			result:   48,
		},
		{
			name:     "Filter by rule type and version",
			ruleType: "xss",
			version:  "v33",
			result:   2,
		},
		{
			name:     "Filter by rule type",
			ruleType: "xss",
			result:   4,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			wafRules := parsePreConfiguredFromFile(t)
			preConfig := filterResult(test.ruleType, test.version, wafRules)
			assert.Equal(t, test.result, len(preConfig))
		})
	}
}

func parsePreConfiguredFromFile(t *testing.T) []*compute.WafExpressionSet {
	data, err := os.ReadFile("testdata/preconfigured.json")
	assert.NoError(t, err)
	var preConfig []*compute.WafExpressionSet
	err = json.Unmarshal(data, &preConfig)
	assert.NoError(t, err)
	return preConfig
}
